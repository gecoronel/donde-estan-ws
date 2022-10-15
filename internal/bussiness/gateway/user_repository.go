//go:generate mockgen -destination ../.././infrastructure/repository/mocks/user.go -package mock . UserRepository

// Package gateway representing the invocation to outer layers, needed by application business logic: external services, data repositories, events, etc.
package gateway

import (
	"github.com/gcoron/donde-estan-ws/internal/bussiness/model"
	_ "github.com/golang/mock/mockgen/model"
)

// UserRepositoryType define IoC key for user repository
const UserRepositoryType = "UserRepository"

// UserRepository is an interface that provides the necessary methods for the user repository.
type UserRepository interface {
	Save(user model.User) (*model.User, error)
	Get(ID uint) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
}
