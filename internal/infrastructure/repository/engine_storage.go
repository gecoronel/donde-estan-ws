//go:generate mockgen -destination ../.././infrastructure/repository/mocks/engine_storage.go -package mock . DB

package repository

import (
	"errors"
	"github.com/gcoron/donde-estan-ws/internal/bussiness/model"
)

type DB interface {
	Save(key string, value interface{}) error
	Get(ID uint) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	Exists(username string) bool
	//FlushAll() error
}

type EngineStorage struct {
	db DB
}

func NewEngineStorage(db DB) *EngineStorage {
	return &EngineStorage{db: db}
}

func (s *EngineStorage) GetUser(ID uint) (*model.User, error) {
	user, err := s.db.Get(ID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *EngineStorage) FindByUsername(username string) (*model.User, error) {
	user, err := s.db.FindByUsername(username)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *EngineStorage) CreateUser(user *model.User) error {
	if s.db.Exists(user.Username) {
		return errors.New("already_exists")
	}

	if err := s.db.Save(user.Username, user); err != nil {
		return errors.New("persistence error")
	}
	return nil
}