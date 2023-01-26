package conn

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDBConnection(logLevel logger.LogLevel, database model.Database) (*gorm.DB, error) {
	dbUsername := database.User
	dbPassword := database.Password
	dbHost := database.Host
	dbName := database.Name
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", dbUsername, dbPassword, dbHost, dbName)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logLevel,
			Colorful:      true,
		},
	)

	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
