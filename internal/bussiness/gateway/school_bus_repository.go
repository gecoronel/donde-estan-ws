//go:generate mockgen --source=user_repository.go --destination=././mocks/user.go

// Package gateway representing the invocation to outer layers, needed by application business logic: external services, data repositories, events, etc.
package gateway

import (
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model"
	_ "github.com/golang/mock/mockgen/model"
)

// SchoolBusRepositoryType define IoC key for user repository
const SchoolBusRepositoryType = "SchoolBusRepository"

// SchoolBusRepository is an interface that provides the necessary methods for the school bus repository.
type SchoolBusRepository interface {
	Get(string) (*model.SchoolBus, error)
	Save(model.SchoolBus) (*model.SchoolBus, error)
	FindByID(string) (*model.SchoolBus, error)
}
