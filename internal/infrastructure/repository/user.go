package repository

import (
	"context"
	"database/sql"
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"

	"github.com/gcoron/donde-estan-ws/internal/bussiness/gateway"
	"github.com/gcoron/donde-estan-ws/internal/bussiness/model"
	"github.com/gcoron/donde-estan-ws/internal/bussiness/model/web"
)

func NewUserRepository(db *gorm.DB, ctx context.Context) gateway.UserRepository {
	return &UserRepository{
		DB:      db,
		context: ctx,
	}
}

// UserRepository represents the main repository for manage user.
type UserRepository struct {
	DB      *gorm.DB
	context context.Context
}

// Get obtains a user using UserRepository by ID.
func (r UserRepository) Get(id uint) (*model.User, error) {
	var user model.User
	result := r.DB.First(&user, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, web.NewError(http.StatusNotFound, "user not found")
		}
		return nil, result.Error
	}

	return &user, nil
}

// Save persists a user using UserRepository.
func (r UserRepository) Save(user model.User) (*model.User, error) {
	result := r.DB.Model(&user).Create(&user)
	return &user, result.Error
}

// FindByUsername obtains a user using UserRepository by username.
func (r UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User

	err := r.DB.
		Raw("SELECT * FROM USERS WHERE username = @username", sql.Named("username", username)).
		Row().
		Scan(&user.ID, &user.Name, &user.LastName, &user.NumberID, &user.Username, &user.Password, &user.Email)

	if err != nil {
		log.Error("error row scan")
		return nil, err
	}

	return &user, nil
}

// GetUsers obtains users using UserRepository .
func (r UserRepository) GetUsers(limit string, offset string) (*[]model.User, error) {
	var user model.User
	var users []model.User

	rows, err := r.DB.Raw(
		"SELECT * FROM USERS LIMIT @limit OFFSET @offset",
		sql.Named("limit", limit),
		sql.Named("offset", offset),
	).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Name, &user.LastName, &user.NumberID, &user.Username, &user.Password, &user.Email)
		if err != nil {
			log.Error("error rows scan")
			return &users, err
		}

		users = append(users, user)
	}

	return &users, nil
}
