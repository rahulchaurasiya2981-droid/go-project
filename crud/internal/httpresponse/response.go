package httpresponse

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type httpResponseSuccess struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type httpResponseError struct {
	Success bool                  `json:"success"`
	Error   httpResponseErrorData `json:"error"`
}

type httpResponseErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// writeJSON writes any response as JSON.
func writeJSON(w http.ResponseWriter, statusCode int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error(
			"Failed to encode HTTP response",
			"action", "RESPONSE_ENCODING_ERROR",
			"error", err,
		)
	}
}

func Success(
	w http.ResponseWriter,
	statusCode int,
	message string,
	data any,
) {
	response := httpResponseSuccess{
		Success: true,
		Message: message,
		Data:    data,
	}

	writeJSON(w, statusCode, response)
}

func Error(
	w http.ResponseWriter,
	statusCode int,
	code string,
	message string,
) {
	response := httpResponseError{
		Success: false,
		Error: httpResponseErrorData{
			Code:    code,
			Message: message,
		},
	}

	writeJSON(w, statusCode, response)
}
