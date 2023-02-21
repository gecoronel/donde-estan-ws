package handler

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Health returns a successful pong answer to all HTTP requests.
func Health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode("ok")
	if err != nil {
		log.Error("error encoding health data: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
