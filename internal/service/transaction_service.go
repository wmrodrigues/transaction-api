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
	GetBalanceByUserId(ctx context.Context, userId string) (*domain.Transaction, error)
	Create(ctx context.Context, transaction *domain.Transaction) error
	GetByUserId(ctx context.Context, userId string, pagination domain.Pagination) (*domain.Page[domain.Transaction], error)
}

type transactionService struct {
	transactionRepository repository.TransactionRepository
	userRepository        repository.UserRepository
	transactionManager    database.Manager
}

func NewTransactionService(transactionRepository repository.TransactionRepository, userRepository repository.UserRepository, transactionManager database.Manager) TransactionService {
	return &transactionService{transactionRepository: transactionRepository, userRepository: userRepository, transactionManager: transactionManager}
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

func (t *transactionService) GetBalanceByUserId(ctx context.Context, userId string) (*domain.Transaction, error) {
	if userId == "" {
		return nil, errors.New("user id cannot be empty and must be a valid UUID")
	}
	transaction, err := t.transactionRepository.GetBalanceByUserId(ctx, userId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("balance not found for specified user")
	}
	return transaction, err
}

func (t *transactionService) Create(ctx context.Context, transaction *domain.Transaction) error {
	return t.transactionManager.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := transaction.Validate(); err != nil {
			return fmt.Errorf("error validating transaction data: %w", err)
		}
		err := t.validateUserExists(ctx, transaction.UserID)
		if err != nil {
			return fmt.Errorf("user id %s doesn't exists: %w", transaction.UserID, err)
		}
		fromBalance, err := t.transactionRepository.GetBalanceByUserId(ctx, transaction.UserID)
		if err != nil {
			return fmt.Errorf("error getting fromBalance by user id: %w", err)
		}
		if fromBalance == nil {
			return errors.New("fromBalance not found for specified user")
		}

		// sending money transaction
		if transaction.ToUserID != nil && transaction.FromUserID != transaction.ToUserID {
			err := t.validateUserExists(ctx, *transaction.ToUserID)
			if err != nil {
				return fmt.Errorf("user id %s doesn't exists: %w", *transaction.ToUserID, err)
			}
			if fromBalance.Balance+transaction.Amount*-1 < 0 {
				return errors.New("insufficient balance to make this transaction")
			}

			toBalance, err := t.transactionRepository.GetBalanceByUserId(ctx, *transaction.ToUserID)
			if err != nil {
				return fmt.Errorf("error getting balance by user id: %w", err)
			}
			if toBalance == nil {
				return errors.New("balance not found for specified user")
			}

			//first we register the withdrawal
			transactionWithdrawal := transaction.Clone()
			transactionWithdrawal.Amount *= -1
			transactionWithdrawal.Balance = fromBalance.Balance + transactionWithdrawal.Amount
			err = t.transactionRepository.Create(ctx, &transactionWithdrawal)

			// then we register the deposit
			transaction.Balance = toBalance.Balance + transaction.Amount
			transaction.UserID = *transaction.ToUserID
			err = t.transactionRepository.Create(ctx, transaction)
			return err
		}

		// when sending money to yourself, we just update the balance
		transaction.Balance = fromBalance.Balance + transaction.Amount
		err = t.transactionRepository.Create(ctx, transaction)
		if err != nil {
			return fmt.Errorf("error creating transaction for user %s: %w", transaction.UserID, err)
		}
		return nil
	})
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
