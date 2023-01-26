package api

import (
	"net/http"
	"os"

	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model"
	"github.com/gecoronel/donde-estan-ws/internal/infrastructure/delivery/api/conn"
	"github.com/gecoronel/donde-estan-ws/internal/infrastructure/delivery/api/route"
	log "github.com/sirupsen/logrus"

	"gorm.io/gorm/logger"
)

const (
	ExitCodeOK = iota
	ExitCodeFailToCreateDBConnection
)

// StartApp Start app
func StartApp() {
	startServer()
}

func startServer() {
	dbConnection, err := conn.GetDBConnection(logger.Silent, model.Database{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Name:     os.Getenv("DB_NAME"),
	})
	if err != nil {
		log.Fatal("couldn't establish a connection with the database", err)
		os.Exit(ExitCodeFailToCreateDBConnection)
	}

	router := route.NewRouter(dbConnection)

	log.Info("server start")
	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))

	log.Info("server exit", err)
	os.Exit(ExitCodeOK)
}
