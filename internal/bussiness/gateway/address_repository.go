//go:generate mockgen --source=address_repository.go --destination=././mocks/address.go

// Package gateway representing the invocation to outer layers, needed by application business logic: external services, data repositories, events, etc.
package gateway

import "github.com/gecoronel/donde-estan-ws/internal/bussiness/model"

// AddressRepositoryType define IoC key for user repository
const AddressRepositoryType = "AddressRepository"

// AddressRepository is an interface that provides the necessary methods for the address repository.
type AddressRepository interface {
	Get(uint64) (*model.Address, error)
	Save(model.Address) (*model.Address, error)
	Update(model.Address) (*model.Address, error)
	Delete(uint64) error
}
