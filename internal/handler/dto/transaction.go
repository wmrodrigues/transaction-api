package dto

import (
	"errors"
	"transaction-api/internal/domain"
)

type Transaction struct {
	FromUserID string  `json:"fromUserId"`
	ToUserID   *string `json:"toUserId"`
	Amount     int64   `json:"amount"`
	Currency   string  `json:"currency"`
}

func (t *Transaction) Validate() error {
	if t.Amount == 0 {
		return errors.New("amount is required")
	}
	if *t.ToUserID != "" && t.Amount < 0 {
		return errors.New("amount must be positive")
	}
	return nil
}

func (t *Transaction) ToDomain() domain.Transaction {
	return domain.Transaction{
		UserID:     t.FromUserID,
		FromUserID: &t.FromUserID,
		ToUserID:   t.ToUserID,
		Currency:   t.Currency,
		Amount:     t.Amount,
	}
}
