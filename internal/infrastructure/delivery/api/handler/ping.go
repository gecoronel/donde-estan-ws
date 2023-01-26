package handler

import (
	"net/http"

	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model/web"
)

// Pong returns a successful pong answer to all HTTP requests.
func Pong(w http.ResponseWriter, r *http.Request) { //c *gin.Context) {
	//c.String(http.StatusOK, "pong")
	web.EncodeJSON(w, "pong", http.StatusOK)
}
