package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// HandleError - Error handler terpusat dengan logging
func HandleError(w http.ResponseWriter, statusCode int, err error, message string) {
	log.Printf("Error [%d]: %s - %v", statusCode, message, err)

	response := ErrorResponse{
		Error:   message,
		Message: err.Error(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// HandleValidationError - Validasi error handler
func HandleValidationError(w http.ResponseWriter, message string) {
	response := ErrorResponse{
		Error: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(response)
}

// HandleSuccess - Success response handler
func HandleSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	response := SuccessResponse{
		Message: message,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
