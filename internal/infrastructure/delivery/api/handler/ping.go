package handler

import (
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model/web"
)

// Pong returns a successful pong answer to all HTTP requests.
func Pong(w http.ResponseWriter, r *http.Request) { //c *gin.Context) {
	//c.String(http.StatusOK, "pong")
	if err := web.EncodeJSON(w, "pong", http.StatusOK); err != nil {
		log.Error("error encoding pong data: ", err)
	}
}
