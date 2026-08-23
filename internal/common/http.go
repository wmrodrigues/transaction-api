package common

import (
	"encoding/json"
	"net/http"
	"strconv"
	"transaction-api/internal/domain"
	"transaction-api/internal/handler/dto"
)

const (
	ContentTypeHeader = "Content-Type"
	ContentTypeJSON   = "application/json"
	DefaultPageSize   = 10
	DefaultPage       = 0
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SendCreatedResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set(ContentTypeHeader, ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func SendInternalServerErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Set(ContentTypeHeader, ContentTypeJSON)
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
}

func SendSuccessResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set(ContentTypeHeader, ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func SendBadRequestResponse(w http.ResponseWriter, err error) {
	w.Header().Set(ContentTypeHeader, ContentTypeJSON)
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()})
}

func SendUnauthorizedResponse(w http.ResponseWriter, err error) {
	w.Header().Set(ContentTypeHeader, ContentTypeJSON)
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Code: http.StatusUnauthorized, Message: err.Error()})
}

func DecodeBodyToObject[O dto.UserRequest | dto.TransactionRequest | dto.AuthRequest](w http.ResponseWriter, r *http.Request, object *O) error {
	if err := json.NewDecoder(r.Body).Decode(&object); err != nil {
		SendBadRequestResponse(w, err)
		return err
	}
	return nil
}

func GetPagination(r *http.Request) (domain.Pagination, error) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = DefaultPage
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		pageSize = DefaultPageSize
	}
	return domain.Pagination{Page: page, PageSize: pageSize}, nil
}
