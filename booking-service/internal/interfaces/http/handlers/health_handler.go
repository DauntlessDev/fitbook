package handlers

import (
	"encoding/json"
	"net/http"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(w http.ResponseWriter, req *http.Request) {
	response := map[string]string{
		"status": "ok",
	}
	json.NewEncoder(w).Encode(response)
}
