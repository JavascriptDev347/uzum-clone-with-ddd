package response

import (
	"encoding/json"
	"net/http"
)

type Envelope struct {
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, payload Envelope) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func Success(w http.ResponseWriter, status int, data any) {
	WriteJSON(w, status, Envelope{Data: data})
}

func Error(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, Envelope{Error: message})
}
