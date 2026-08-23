package service

import (
	"context"
	"errors"
	"fmt"
	"transaction-api/internal/database"
	"transaction-api/internal/domain"
	"transaction-api/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionService interface {
	GetByID(ctx context.Context, id string) (*domain.Transaction, error)
	GetBalanceByUserId(ctx context.Context, userID string) ([]domain.Account, error)
	Create(ctx context.Context, transaction *domain.Transaction) error
	GetByUserId(ctx context.Context, userId string, pagination domain.Pagination) (*domain.Page[domain.Transaction], error)
}

type transactionService struct {
	transactionRepository repository.TransactionRepository
	userRepository        repository.UserRepository
	accountRepository     repository.AccountRepository
	transactionManager    database.Manager
}

func NewTransactionService(transactionRepository repository.TransactionRepository, userRepository repository.UserRepository, accountRepository repository.AccountRepository, transactionManager database.Manager) TransactionService {
	return &transactionService{
		transactionRepository: transactionRepository,
		userRepository:        userRepository,
		accountRepository:     accountRepository,
		transactionManager:    transactionManager,
	}
}

func (t *transactionService) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	if id == "" {
		return nil, errors.New("transaction id cannot be empty and must be a valid UUID")
	}
	transaction, err := t.transactionRepository.GetByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("transaction not found")
	}
	return transaction, err
}

func (t *transactionService) GetBalanceByUserId(ctx context.Context, userID string) ([]domain.Account, error) {
	if userID == "" {
		return nil, errors.New("user id cannot be empty and must be a valid UUID")
	}
	return t.accountRepository.GetByUserID(ctx, userID)
}

func (t *transactionService) Create(ctx context.Context, transaction *domain.Transaction) error {
	return t.transactionManager.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := transaction.Validate(); err != nil {
			return fmt.Errorf("error validating transaction data: %w", err)
		}
		if err := t.validateUserExists(ctx, transaction.UserID); err != nil {
			return fmt.Errorf("user id %s doesn't exists: %w", transaction.UserID, err)
		}
		if transaction.ToUserID != nil && *transaction.ToUserID != transaction.UserID {
			return t.transfer(ctx, transaction)
		}
		return t.deposit(ctx, transaction)
	})
}

func (t *transactionService) deposit(ctx context.Context, transaction *domain.Transaction) error {
	account, err := t.lockAccount(ctx, transaction.UserID, transaction.Currency)
	if err != nil {
		return err
	}
	if err := t.accountRepository.UpdateBalance(ctx, account.ID, account.Balance+transaction.Amount); err != nil {
		return err
	}
	if err := t.transactionRepository.Create(ctx, transaction); err != nil {
		return fmt.Errorf("error creating transaction for user %s: %w", transaction.UserID, err)
	}
	return nil
}

func (t *transactionService) transfer(ctx context.Context, transaction *domain.Transaction) error {
	receiverID := *transaction.ToUserID
	if err := t.validateUserExists(ctx, receiverID); err != nil {
		return fmt.Errorf("user id %s doesn't exists: %w", receiverID, err)
	}
	// sorting the users (A->B and B->A) so no deadlock occurs in case one tries to transfer to another at the same time
	firstUser, secondUser := transaction.UserID, receiverID
	if firstUser > secondUser {
		firstUser, secondUser = secondUser, firstUser
	}
	firstAccount, err := t.lockAccount(ctx, firstUser, transaction.Currency)
	if err != nil {
		return err
	}
	secondAccount, err := t.lockAccount(ctx, secondUser, transaction.Currency)
	if err != nil {
		return err
	}
	from, to := firstAccount, secondAccount
	if firstUser != transaction.UserID {
		from, to = secondAccount, firstAccount
	}
	if from.Balance < transaction.Amount {
		return errors.New("insufficient balance to make this transaction")
	}
	if err := t.accountRepository.UpdateBalance(ctx, from.ID, from.Balance-transaction.Amount); err != nil {
		return fmt.Errorf("error debiting sender account: %w", err)
	}
	if err := t.accountRepository.UpdateBalance(ctx, to.ID, to.Balance+transaction.Amount); err != nil {
		return fmt.Errorf("error crediting receiver account: %w", err)
	}
	// record the ledger: a withdrawal for the sender and a transactionDeposit for the receiver.
	transactionWithdrawal := transaction.Clone()
	transactionWithdrawal.Amount = -transaction.Amount
	if err := t.transactionRepository.Create(ctx, &transactionWithdrawal); err != nil {
		return fmt.Errorf("error recording withdrawal: %w", err)
	}
	transactionDeposit := transaction.Clone()
	transactionDeposit.UserID = receiverID
	if err := t.transactionRepository.Create(ctx, &transactionDeposit); err != nil {
		return fmt.Errorf("error recording deposit: %w", err)
	}
	return nil
}

func (t *transactionService) lockAccount(ctx context.Context, userID, currency string) (*domain.Account, error) {
	account, err := t.accountRepository.GetForUpdate(ctx, userID, currency)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("no %s account found for user %s", currency, userID)
	}
	if err != nil {
		return nil, fmt.Errorf("error locking account: %w", err)
	}
	return account, nil
}

func (t *transactionService) validateUserExists(ctx context.Context, id string) error {
	_, err := t.userRepository.GetByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	}
	if err != nil {
		return fmt.Errorf("error getting user by id: %w", err)
	}
	return nil
}

func (t *transactionService) GetByUserId(ctx context.Context, userId string, pagination domain.Pagination) (*domain.Page[domain.Transaction], error) {
	if err := uuid.Validate(userId); err != nil {
		return nil, errors.New("user_id must be a valid UUID")
	}
	return t.transactionRepository.GetByUserId(ctx, userId, pagination)
}
