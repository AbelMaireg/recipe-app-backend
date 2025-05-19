package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorMessage struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func WriteError(w http.ResponseWriter, status int, errorCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	errorMessage := ErrorMessage{
		Code:        errorCode,
		Description: message,
	}
	// Serialize ErrorMessage as a JSON string
	jsonMessage, _ := json.Marshal(errorMessage)
	json.NewEncoder(w).Encode(map[string]string{
		"message": string(jsonMessage),
	})
}

func DecodeJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func EncodeJSON(w http.ResponseWriter, v any) {
	json.NewEncoder(w).Encode(v)
}
