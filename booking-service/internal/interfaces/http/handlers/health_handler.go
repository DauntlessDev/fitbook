package handlers

import (
	"encoding/json"
	"net/http"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (handler *HealthHandler) Check(writer http.ResponseWriter, request *http.Request) {
	response := map[string]string{
		"status": "ok",
	}
	json.NewEncoder(writer).Encode(response)
}
