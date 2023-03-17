//go:generate mockgen --source=child_repository.go --destination=././mocks/child.go

// Package gateway representing the invocation to outer layers, needed by application business logic: external services, data repositories, events, etc.
package gateway

import "github.com/gecoronel/donde-estan-ws/internal/bussiness/model"

// ChildRepositoryType define IoC key for user repository
const ChildRepositoryType = "ChildRepository"

// ChildRepository is an interface that provides the necessary methods for the address repository.
type ChildRepository interface {
	Get(uint64) (*model.Child, error)
	Save(model.Child) (*model.Child, error)
	Update(model.Child) (*model.Child, error)
	Delete(uint64) error
}
