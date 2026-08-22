package handler

import (
	"fmt"
	"net/http"
	"transaction-api/internal/authentication"
	"transaction-api/internal/common"
	"transaction-api/internal/domain"
	"transaction-api/internal/handler/dto"
	"transaction-api/internal/service"
)

type TransactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler(transactionService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

func (th *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	transaction, err := th.transactionService.GetByID(ctx, id)
	if err != nil {
		common.SendBadRequestResponse(w, err)
		return
	}
	claims, ok := authentication.ClaimsFromContext(r.Context())
	if !ok {
		common.SendUnauthorizedResponse(w, fmt.Errorf("unauthorized: missing claims"))
		return
	}
	if claims.UserID != transaction.UserID {
		common.SendUnauthorizedResponse(w, fmt.Errorf("unauthorized: user id mismatch"))
		return
	}
	common.SendSuccessResponse(w, transaction)
}

func (th *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var transactionRequest dto.TransactionRequest
	err := common.DecodeBodyToObject(w, r, &transactionRequest)
	if err != nil {
		return
	}
	claims, ok := authentication.ClaimsFromContext(r.Context())
	if !ok {
		common.SendUnauthorizedResponse(w, fmt.Errorf("unauthorized: missing claims"))
		return
	}
	if claims.UserID != transactionRequest.FromUserID {
		common.SendUnauthorizedResponse(w, fmt.Errorf("unauthorized: user id mismatch"))
		return
	}
	var transaction domain.Transaction
	transaction = transactionRequest.ToDomain()
	err = th.transactionService.Create(ctx, &transaction)
	if err != nil {
		common.SendBadRequestResponse(w, err)
		return
	}
	common.SendCreatedResponse(w, nil)
}

func (th *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	pagination, err := common.GetPagination(r)
	if err != nil {
		common.SendBadRequestResponse(w, err)
		return
	}
	userId := r.PathValue("id")
	claims, ok := authentication.ClaimsFromContext(r.Context())
	if !ok {
		common.SendUnauthorizedResponse(w, fmt.Errorf("unauthorized: missing claims"))
		return
	}
	if claims.UserID != userId {
		common.SendUnauthorizedResponse(w, fmt.Errorf("unauthorized: user id mismatch"))
		return
	}
	result, err := th.transactionService.GetByUserId(r.Context(), userId, pagination)
	if err != nil {
		common.SendBadRequestResponse(w, err)
		return
	}
	common.SendSuccessResponse(w, result)
}

func (th *TransactionHandler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("id")
	claims, ok := authentication.ClaimsFromContext(r.Context())
	if !ok {
		common.SendUnauthorizedResponse(w, fmt.Errorf("unauthorized: missing claims"))
		return
	}
	if claims.UserID != userId {
		common.SendUnauthorizedResponse(w, fmt.Errorf("unauthorized: user id mismatch"))
		return
	}
	result, err := th.transactionService.GetBalanceByUserId(r.Context(), userId)
	if err != nil {
		common.SendBadRequestResponse(w, err)
		return
	}
	common.SendSuccessResponse(w, result)
}
