package common

import (
	"encoding/json"
	"net/http"
	"transaction-api/internal/domain"
)

func SendCreatedResponse(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(data)
}

func SendErrorResponse(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func SendSuccessResponse(w http.ResponseWriter, data interface{}) {
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
