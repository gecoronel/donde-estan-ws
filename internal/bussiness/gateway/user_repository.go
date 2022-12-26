//go:generate mockgen --source=user_repository.go --destination=./internal/./infrastructure/repository/mocks/user.go

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
	Save(model.User) (*model.User, error)
	Get(uint) (*model.User, error)
	GetUsers(string, string) (*[]model.User, error)
	FindByUsername(string) (*model.User, error)
	GetObservedUser(*model.ObservedUser) (*model.IUser, error)
	GetObserverUser(*model.ObserverUser) (*model.IUser, error)
}
