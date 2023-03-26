//go:generate mockgen --source=school_bus_repository.go --destination=././mocks/school_bus.go

// Package gateway representing the invocation to outer layers, needed by application business logic: external services, data repositories, events, etc.
package gateway

import "github.com/gecoronel/donde-estan-ws/internal/bussiness/model"

// SchoolBusRepositoryType define IoC key for user repository
const SchoolBusRepositoryType = "SchoolBusRepository"

// SchoolBusRepository is an interface that provides the necessary methods for the school bus repository.
type SchoolBusRepository interface {
	Get(uint64) (*model.SchoolBus, error)
	Save(model.SchoolBus) (*model.SchoolBus, error)
	Update(model.SchoolBus) (*model.SchoolBus, error)
	Delete(uint64) error
}
