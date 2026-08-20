package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	FromUserID *string   `json:"from_user_id"`
	Currency   string    `json:"currency"`
	Amount     int64     `json:"amount"`
	Balance    int64     `json:"balance"`
	CreatedAt  time.Time `json:"created_at"`
}

func (t *Transaction) Validate() error {
	if t.UserID == "" {
		return errors.New("user_id is required")
	}
	if t.Currency == "" {
		return errors.New("currency is required")
	}
	if !t.IsSupportedCurrency(CurrencyKey(t.Currency)) {
		return errors.New("currency is not supported")
	}
	if err := uuid.Validate(t.UserID); err != nil {
		return errors.New("user_id must be a valid UUID")
	}
	if t.FromUserID != nil {
		if err := uuid.Validate(*t.FromUserID); err != nil {
			return errors.New("from_user_id must be a valid UUID")
		}
	}
	return nil
}

type CurrencyKey string

const (
	USD CurrencyKey = "USD"
	EUR CurrencyKey = "EUR"
	GBP CurrencyKey = "GBP"
	JPY CurrencyKey = "JPY"
	CHF CurrencyKey = "CHF"
	CAD CurrencyKey = "CAD"
	AUD CurrencyKey = "AUD"
	NZD CurrencyKey = "NZD"
	CNY CurrencyKey = "CNY"
	HKD CurrencyKey = "HKD"
	SGD CurrencyKey = "SGD"
	INR CurrencyKey = "INR"
	BRL CurrencyKey = "BRL"
	MXN CurrencyKey = "MXN"
	ARS CurrencyKey = "ARS"
	CLP CurrencyKey = "CLP"
	COP CurrencyKey = "COP"
	PEN CurrencyKey = "PEN"
	ZAR CurrencyKey = "ZAR"
	SEK CurrencyKey = "SEK"
	NOK CurrencyKey = "NOK"
	DKK CurrencyKey = "DKK"
	PLN CurrencyKey = "PLN"
	CZK CurrencyKey = "CZK"
	HUF CurrencyKey = "HUF"
	TRY CurrencyKey = "TRY"
	AED CurrencyKey = "AED"
	SAR CurrencyKey = "SAR"
	ILS CurrencyKey = "ILS"
	KRW CurrencyKey = "KRW"
	THB CurrencyKey = "THB"
	IDR CurrencyKey = "IDR"
	MYR CurrencyKey = "MYR"
	PHP CurrencyKey = "PHP"
	VND CurrencyKey = "VND"
	RUB CurrencyKey = "RUB"
	UAH CurrencyKey = "UAH"
)

var SupportedCurrencies = map[CurrencyKey]struct{}{SGD: {}}

func (t *Transaction) IsSupportedCurrency(c CurrencyKey) bool {
	_, ok := SupportedCurrencies[c]
	return ok
}
