package handler

import (
	"net/http"
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
	user, err := th.transactionService.GetByID(ctx, id)
	if err != nil {
		common.SendBadRequestResponse(w, err)
		return
	}
	common.SendSuccessResponse(w, user)
}

func (th *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var transactionRequest dto.Transaction
	err := common.DecodeBodyToObject(w, r, &transactionRequest)
	if err != nil {
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
	result, err := th.transactionService.GetByUserId(r.Context(), userId, pagination)
	if err != nil {
		common.SendBadRequestResponse(w, err)
		return
	}
	common.SendSuccessResponse(w, result)
}
