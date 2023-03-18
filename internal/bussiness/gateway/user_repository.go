//go:generate mockgen --source=user_repository.go --destination=././mocks/user.go

// Package gateway representing the invocation to outer layers, needed by application business logic: external services, data repositories, events, etc.
package gateway

import (
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model"
	_ "github.com/golang/mock/mockgen/model"
)

// UserRepositoryType define IoC key for user repository
const UserRepositoryType = "UserRepository"

// UserRepository is an interface that provides the necessary methods for the user repository.
type UserRepository interface {
	Get(uint64) (*model.User, error)
	GetUsers(string, string) (*[]model.User, error)
	FindByUsername(string) (*model.User, error)
	FindByEmail(string) (*model.User, error)
	GetObservedUser(uint64) (*model.ObservedUser, error)
	SaveObservedUser(model.ObservedUser) (*model.ObservedUser, error)
	UpdateObservedUser(model.ObservedUser) (*model.ObservedUser, error)
	DeleteObservedUser(uint64) error
	GetObserverUser(uint64) (*model.ObserverUser, error)
	SaveObserverUser(model.ObserverUser) (*model.ObserverUser, error)
	UpdateObserverUser(model.ObserverUser) (*model.ObserverUser, error)
	DeleteObserverUser(uint64) error
	FindObservedUserByPrivacyKey(string) (*model.ObservedUser, error)
	SaveObservedUserInObserverUser(uint64, uint64) error
	DeleteObservedUserInObserverUser(uint64, uint64) error
}
