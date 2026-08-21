package service

import (
	"context"
	"errors"
	"testing"
	"transaction-api/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type mockTransactionRepository struct {
	transaction     *domain.Transaction
	balanceByUserID map[string]*domain.Transaction
	err             error
	getBalanceErr   error
	createErr       error
	created         []*domain.Transaction
}

func (m *mockTransactionRepository) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	return m.transaction, m.err
}

func (m *mockTransactionRepository) GetBalanceByUserId(ctx context.Context, id string) (*domain.Transaction, error) {
	if m.getBalanceErr != nil {
		return nil, m.getBalanceErr
	}
	if m.balanceByUserID != nil {
		if tx, ok := m.balanceByUserID[id]; ok {
			return tx, nil
		}
	}
	return m.transaction, m.err
}

func (m *mockTransactionRepository) Create(ctx context.Context, transaction *domain.Transaction) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.created = append(m.created, transaction)
	return m.err
}

func (m *mockTransactionRepository) GetByUserId(ctx context.Context, userId string, pagination domain.Pagination) (*domain.Page[domain.Transaction], error) {
	return nil, m.err
}

func TestGetBalanceByUserId(t *testing.T) {
	expected := &domain.Transaction{
		ID:      "tx123",
		UserID:  "123",
		Amount:  500,
		Balance: 1500,
	}
	repository := &mockTransactionRepository{transaction: expected}
	transactionService := NewTransactionService(repository, &mockUserRepository{}, &mockTransactionManager{})
	result, err := transactionService.GetBalanceByUserId(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Balance != expected.Balance {
		t.Fatalf("expected balance %d, got %d", expected.Balance, result.Balance)
	}
}

func TestGetBalanceByUserId_EmptyUserID(t *testing.T) {
	transactionService := NewTransactionService(&mockTransactionRepository{}, &mockUserRepository{}, &mockTransactionManager{})
	_, err := transactionService.GetBalanceByUserId(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty user ID, got nil")
	}
	if err.Error() != "user id cannot be empty and must be a valid UUID" {
		t.Fatalf("expected empty ID error message, got %v", err)
	}
}

func TestCreateTransaction_SelfDeposit(t *testing.T) {
	userID := uuid.New().String()
	transactionRepository := &mockTransactionRepository{
		transaction: &domain.Transaction{
			UserID:  userID,
			Balance: 1000,
		},
	}
	userRepository := &mockUserRepository{
		user: &domain.User{ID: userID, Name: "Wash", Email: "wash@example.com"},
	}
	transactionService := NewTransactionService(transactionRepository, userRepository, &mockTransactionManager{})
	transaction := &domain.Transaction{
		UserID:   userID,
		Amount:   500,
		Currency: "SGD",
	}
	err := transactionService.Create(context.Background(), transaction)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if transaction.Balance != 1500 {
		t.Fatalf("expected balance 1500, got %d", transaction.Balance)
	}
}

func TestCreateTransaction_ValidationFailure(t *testing.T) {
	transactionService := NewTransactionService(&mockTransactionRepository{}, &mockUserRepository{}, &mockTransactionManager{})
	tests := []struct {
		name        string
		transaction domain.Transaction
	}{
		{"missing user_id", domain.Transaction{Currency: "SGD", Amount: 100}},
		{"missing currency", domain.Transaction{UserID: uuid.New().String(), Amount: 100}},
		{"unsupported currency", domain.Transaction{UserID: uuid.New().String(), Currency: "XYZ", Amount: 100}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tx := tc.transaction
			err := transactionService.Create(context.Background(), &tx)
			if err == nil {
				t.Fatalf("expected validation error for %s, got nil", tc.name)
			}
		})
	}
}

func TestCreateTransaction_UserDoesNotExist(t *testing.T) {
	userRepository := &mockUserRepository{err: gorm.ErrRecordNotFound}
	transactionService := NewTransactionService(&mockTransactionRepository{}, userRepository, &mockTransactionManager{})
	transaction := &domain.Transaction{
		UserID:   "123",
		Amount:   500,
		Currency: "SGD",
	}
	err := transactionService.Create(context.Background(), transaction)
	if err == nil {
		t.Fatal("expected error for non-existent user, got nil")
	}
}

func TestCreateTransaction_NoExistingBalance(t *testing.T) {
	userID := "123"
	transactionRepository := &mockTransactionRepository{getBalanceErr: errors.New("no balance record")}
	userRepository := &mockUserRepository{user: &domain.User{ID: userID, Name: "Wash", Email: "wash@example.com"}}
	transactionService := NewTransactionService(transactionRepository, userRepository, &mockTransactionManager{})
	tx := &domain.Transaction{
		UserID:   userID,
		Amount:   500,
		Currency: "SGD",
	}
	err := transactionService.Create(context.Background(), tx)
	if err == nil {
		t.Fatal("expected error when balance not found, got nil")
	}
}

func TestCreateTransaction_Transfer(t *testing.T) {
	senderID := uuid.New().String()
	receiverID := uuid.New().String()

	transactionRepository := &mockTransactionRepository{
		balanceByUserID: map[string]*domain.Transaction{
			senderID:   {UserID: senderID, Balance: 1000},
			receiverID: {UserID: receiverID, Balance: 200},
		},
	}
	userRepository := &mockUserRepository{user: &domain.User{ID: senderID, Name: "Wash", Email: "wash@example.com"}}
	transactionService := NewTransactionService(transactionRepository, userRepository, &mockTransactionManager{})
	transaction := &domain.Transaction{
		UserID:   senderID,
		ToUserID: &receiverID,
		Amount:   300,
		Currency: "SGD",
	}
	err := transactionService.Create(context.Background(), transaction)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(transactionRepository.created) != 2 {
		t.Fatalf("expected 2 transaction records (withdrawal + deposit), got %d", len(transactionRepository.created))
	}
	withdrawal := transactionRepository.created[0]
	if withdrawal.Amount != -300 {
		t.Fatalf("expected withdrawal amount -300, got %d", withdrawal.Amount)
	}
	if withdrawal.Balance != 700 {
		t.Fatalf("expected sender balance 700 after withdrawal, got %d", withdrawal.Balance)
	}
	deposit := transactionRepository.created[1]
	if deposit.Balance != 500 {
		t.Fatalf("expected receiver balance 500 after deposit, got %d", deposit.Balance)
	}
}

func TestCreateTransaction_Transfer_InsufficientBalance(t *testing.T) {
	senderID := uuid.New().String()
	receiverID := uuid.New().String()
	transactionRepository := &mockTransactionRepository{balanceByUserID: map[string]*domain.Transaction{senderID: {UserID: senderID, Balance: 100}}}
	userRepository := &mockUserRepository{user: &domain.User{ID: senderID, Name: "Wash", Email: "wash@example.com"}}
	transactionService := NewTransactionService(transactionRepository, userRepository, &mockTransactionManager{})
	transaction := &domain.Transaction{
		UserID:   senderID,
		ToUserID: &receiverID,
		Amount:   500,
		Currency: "SGD",
	}
	err := transactionService.Create(context.Background(), transaction)
	if err == nil {
		t.Fatal("expected insufficient balance error, got nil")
	}
	if err.Error() != "insufficient balance to make this transaction" {
		t.Fatalf("expected insufficient balance message, got %v", err)
	}
}

func TestCreateTransaction_Transfer_RecipientDoesNotExist(t *testing.T) {
	senderID := uuid.New().String()
	receiverID := uuid.New().String()
	transactionRepository := &mockTransactionRepository{balanceByUserID: map[string]*domain.Transaction{
		senderID: {UserID: senderID, Balance: 1000}},
	}
	// using a mock that fails for the recipient lookup
	userRepository := &mockUserRepository{err: gorm.ErrRecordNotFound}
	transactionService := NewTransactionService(transactionRepository, userRepository, &mockTransactionManager{})
	transaction := &domain.Transaction{
		UserID:   senderID,
		ToUserID: &receiverID,
		Amount:   300,
		Currency: "SGD",
	}
	err := transactionService.Create(context.Background(), transaction)
	if err == nil {
		t.Fatal("expected error for non-existent recipient, got nil")
	}
}
