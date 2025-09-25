package response

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResp struct {
	Error string `json:"error"`
}

func NewErrorResp(msg string) ErrorResp {
	return ErrorResp{
		Error: msg,
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "не получилось разобрать json", http.StatusInternalServerError)
	}
}

func HandleError(w http.ResponseWriter, err error, status int, userMessage string) {
	log.Printf("[ERROR] %s: %v", userMessage, err)
	WriteJSON(w, status, NewErrorResp(userMessage))
}