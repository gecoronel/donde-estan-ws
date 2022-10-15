package handler

import (
	"github.com/gcoron/donde-estan-ws/internal/bussiness/model/web"
	"net/http"
)

// Pong returns a successful pong answer to all HTTP requests.
func Pong(w http.ResponseWriter, r *http.Request) { //c *gin.Context) {
	//c.String(http.StatusOK, "pong")
	web.EncodeJSON(w, "pong", http.StatusOK)
}
