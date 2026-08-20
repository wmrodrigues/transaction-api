package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"transaction-api/internal/domain"
)

const (
	ContentTypeHeader = "Content-Type"
	ContentTypeJSON   = "application/json"
)

func SendCreatedResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set(ContentTypeHeader, ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(data)
}

func SendErrorResponse(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func SendSuccessResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set(ContentTypeHeader, ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func SendBadRequestResponse(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func DecodeBodyToObject[O domain.User](w http.ResponseWriter, r *http.Request, object *O) error {
	if err := json.NewDecoder(r.Body).Decode(&object); err != nil {
		SendBadRequestResponse(w, err)
		return err
	}
	return nil
}

func GetPagination(r *http.Request) (domain.Pagination, error) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		return domain.Pagination{}, fmt.Errorf("error parsing page: %w", err)
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		return domain.Pagination{}, fmt.Errorf("error parsing page_size: %w", err)
	}
	return domain.Pagination{Page: page, PageSize: pageSize}, nil
}
