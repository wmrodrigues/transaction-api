package handler

import (
	"net/http"
	"time"
	"transaction-api/internal/common"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (hh *HealthHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	common.SendSuccessResponse(w, map[string]string{"status": "ok", "timestamp": now.Format(time.RFC3339Nano)})
}
