package service

import (
	"context"
	"testing"
	"transaction-api/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type mockTransactionRepository struct {
	transaction *domain.Transaction
	err         error
	createErr   error
	created     []*domain.Transaction
}

func (m *mockTransactionRepository) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
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

type mockAccountRepository struct {
	accounts        map[string]*domain.Account
	byUser          map[string][]domain.Account
	getForUpdateErr error
	updateErr       error
	createErr       error
	updated         map[string]int64
	created         []*domain.Account
}

func (m *mockAccountRepository) Create(ctx context.Context, account *domain.Account) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.created = append(m.created, account)
	return nil
}

func acctKey(userID, currency string) string {
	return userID + "|" + currency
}

func (m *mockAccountRepository) GetForUpdate(ctx context.Context, userID, currency string) (*domain.Account, error) {
	if m.getForUpdateErr != nil {
		return nil, m.getForUpdateErr
	}
	if m.accounts != nil {
		if account, ok := m.accounts[acctKey(userID, currency)]; ok {
			return account, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAccountRepository) GetByUserID(ctx context.Context, userID string) ([]domain.Account, error) {
	if m.byUser != nil {
		return m.byUser[userID], nil
	}
	return nil, nil
}

func (m *mockAccountRepository) UpdateBalance(ctx context.Context, accountID string, newBalance int64) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if m.updated == nil {
		m.updated = map[string]int64{}
	}
	m.updated[accountID] = newBalance
	return nil
}

func TestGetBalances(t *testing.T) {
	accounts := []domain.Account{{ID: "account1", UserID: "123", Currency: "SGD", Balance: 1500}}
	accountRepository := &mockAccountRepository{byUser: map[string][]domain.Account{"123": accounts}}
	transactionService := NewTransactionService(&mockTransactionRepository{}, &mockUserRepository{}, accountRepository, &mockTransactionManager{})
	result, err := transactionService.GetBalanceByUserId(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 account, got %d", len(result))
	}
	if result[0].Balance != 1500 {
		t.Fatalf("expected balance 1500, got %d", result[0].Balance)
	}
}

func TestGetBalances_EmptyUserID(t *testing.T) {
	transactionService := NewTransactionService(&mockTransactionRepository{}, &mockUserRepository{}, &mockAccountRepository{}, &mockTransactionManager{})
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
	account := &domain.Account{ID: "account1", UserID: userID, Currency: "SGD", Balance: 1000}
	accountRepository := &mockAccountRepository{accounts: map[string]*domain.Account{acctKey(userID, "SGD"): account}}
	transactionRepository := &mockTransactionRepository{}
	userRepository := &mockUserRepository{user: &domain.User{ID: userID, Name: "Wash", Email: "wash@example.com"}}
	transactionService := NewTransactionService(transactionRepository, userRepository, accountRepository, &mockTransactionManager{})
	transaction := &domain.Transaction{
		UserID:     userID,
		FromUserID: &userID,
		Amount:     500,
		Currency:   "SGD",
	}
	err := transactionService.Create(context.Background(), transaction)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if accountRepository.updated["account1"] != 1500 {
		t.Fatalf("expected account balance 1500, got %d", accountRepository.updated["account1"])
	}
	if len(transactionRepository.created) != 1 {
		t.Fatalf("expected 1 ledger record, got %d", len(transactionRepository.created))
	}
}

func TestCreateTransaction_ValidationFailure(t *testing.T) {
	transactionService := NewTransactionService(&mockTransactionRepository{}, &mockUserRepository{}, &mockAccountRepository{}, &mockTransactionManager{})
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
	userID := uuid.New().String()
	userRepository := &mockUserRepository{missingIDs: map[string]bool{userID: true}}
	transactionService := NewTransactionService(&mockTransactionRepository{}, userRepository, &mockAccountRepository{}, &mockTransactionManager{})
	transaction := &domain.Transaction{
		UserID:   userID,
		Amount:   500,
		Currency: "SGD",
	}
	err := transactionService.Create(context.Background(), transaction)
	if err == nil {
		t.Fatal("expected error for non-existent user, got nil")
	}
}

func TestCreateTransaction_NoExistingAccount(t *testing.T) {
	userID := uuid.New().String()
	userRepository := &mockUserRepository{user: &domain.User{ID: userID, Name: "Wash", Email: "wash@example.com"}}
	// empty account repository, GetForUpdate returns ErrRecordNotFound
	transactionService := NewTransactionService(&mockTransactionRepository{}, userRepository, &mockAccountRepository{}, &mockTransactionManager{})
	transaction := &domain.Transaction{
		UserID:     userID,
		FromUserID: &userID,
		Amount:     500,
		Currency:   "SGD",
	}
	err := transactionService.Create(context.Background(), transaction)
	if err == nil {
		t.Fatal("expected error when account not found, got nil")
	}
}

func TestCreateTransaction_Transfer(t *testing.T) {
	senderID := uuid.New().String()
	receiverID := uuid.New().String()

	senderAccount := &domain.Account{ID: "accountSender", UserID: senderID, Currency: "SGD", Balance: 1000}
	receiverAccount := &domain.Account{ID: "accountReceiver", UserID: receiverID, Currency: "SGD", Balance: 200}
	accountRepository := &mockAccountRepository{accounts: map[string]*domain.Account{
		acctKey(senderID, "SGD"):   senderAccount,
		acctKey(receiverID, "SGD"): receiverAccount,
	}}
	transactionRepository := &mockTransactionRepository{}
	userRepository := &mockUserRepository{user: &domain.User{ID: senderID, Name: "Wash", Email: "wash@example.com"}}
	transactionService := NewTransactionService(transactionRepository, userRepository, accountRepository, &mockTransactionManager{})
	transaction := &domain.Transaction{
		UserID:     senderID,
		FromUserID: &senderID,
		ToUserID:   &receiverID,
		Amount:     300,
		Currency:   "SGD",
	}
	err := transactionService.Create(context.Background(), transaction)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if accountRepository.updated["accountSender"] != 700 {
		t.Fatalf("expected sender balance 700, got %d", accountRepository.updated["accountSender"])
	}
	if accountRepository.updated["accountReceiver"] != 500 {
		t.Fatalf("expected receiver balance 500, got %d", accountRepository.updated["accountReceiver"])
	}
	if len(transactionRepository.created) != 2 {
		t.Fatalf("expected 2 ledger records (withdrawal + deposit), got %d", len(transactionRepository.created))
	}
	withdrawal := transactionRepository.created[0]
	if withdrawal.Amount != -300 {
		t.Fatalf("expected withdrawal amount -300, got %d", withdrawal.Amount)
	}
	deposit := transactionRepository.created[1]
	if deposit.UserID != receiverID {
		t.Fatalf("expected deposit user id %s, got %s", receiverID, deposit.UserID)
	}
}

func TestCreateTransaction_Transfer_InsufficientBalance(t *testing.T) {
	senderID := uuid.New().String()
	receiverID := uuid.New().String()
	senderAccount := &domain.Account{ID: "accountSender", UserID: senderID, Currency: "SGD", Balance: 100}
	receiverAccount := &domain.Account{ID: "accountReceiver", UserID: receiverID, Currency: "SGD", Balance: 0}
	accountRepository := &mockAccountRepository{accounts: map[string]*domain.Account{
		acctKey(senderID, "SGD"):   senderAccount,
		acctKey(receiverID, "SGD"): receiverAccount,
	}}
	userRepository := &mockUserRepository{user: &domain.User{ID: senderID, Name: "Wash", Email: "wash@example.com"}}
	transactionService := NewTransactionService(&mockTransactionRepository{}, userRepository, accountRepository, &mockTransactionManager{})
	transaction := &domain.Transaction{
		UserID:     senderID,
		FromUserID: &senderID,
		ToUserID:   &receiverID,
		Amount:     500,
		Currency:   "SGD",
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
	userRepository := &mockUserRepository{
		user:       &domain.User{ID: senderID, Name: "Wash", Email: "wash@example.com"},
		missingIDs: map[string]bool{receiverID: true},
	}
	transactionService := NewTransactionService(&mockTransactionRepository{}, userRepository, &mockAccountRepository{}, &mockTransactionManager{})
	transaction := &domain.Transaction{
		UserID:     senderID,
		FromUserID: &senderID,
		ToUserID:   &receiverID,
		Amount:     300,
		Currency:   "SGD",
	}
	err := transactionService.Create(context.Background(), transaction)
	if err == nil {
		t.Fatal("expected error for non-existent recipient, got nil")
	}
}
