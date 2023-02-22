package handler

import (
	"encoding/json"
	"net/http"
)

// Health returns a successful pong answer to all HTTP requests.
func Health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("ok")
}
