//go:generate mockgen --source=address.go --destination=././mocks/address.go

package usecase

import (
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/gateway"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model/web"
)

const AddressUseCaseType = "AddressUseCase"

type (
	AddressUseCase interface {
		Get(uint64, gateway.ServiceLocator) (*model.Address, error)
		Save(model.Address, gateway.ServiceLocator) (*model.Address, error)
		Update(model.Address, gateway.ServiceLocator) (*model.Address, error)
		Delete(uint64, gateway.ServiceLocator) error
	}

	addressUseCase struct{}
)

func NewAddressUseCase() AddressUseCase {
	return &addressUseCase{}
}

func (b addressUseCase) Get(id uint64, locator gateway.ServiceLocator) (*model.Address, error) {
	repository := locator.GetInstance(gateway.AddressRepositoryType).(gateway.AddressRepository)

	Address, err := repository.Get(id)
	if err != nil {
		return nil, web.ErrInternalServerError
	}

	if Address == nil {
		return nil, web.ErrNotFound
	}

	return Address, nil
}

func (b addressUseCase) Save(address model.Address, locator gateway.ServiceLocator) (
	*model.Address,
	error,
) {
	repository := locator.GetInstance(gateway.AddressRepositoryType).(gateway.AddressRepository)

	sb, err := repository.Save(address)
	if err != nil {
		return nil, err
	}

	return sb, nil
}

func (b addressUseCase) Update(Address model.Address, locator gateway.ServiceLocator) (
	*model.Address,
	error,
) {
	repository := locator.GetInstance(gateway.AddressRepositoryType).(gateway.AddressRepository)

	sb, err := repository.Get(Address.ID)
	if err != nil {
		return nil, web.ErrInternalServerError
	}

	if sb == nil {
		return nil, web.ErrNotFound
	}

	sb, err = repository.Update(Address)
	if err != nil {
		return nil, web.ErrInternalServerError
	}

	return sb, nil
}

func (b addressUseCase) Delete(id uint64, locator gateway.ServiceLocator) error {
	repository := locator.GetInstance(gateway.AddressRepositoryType).(gateway.AddressRepository)

	sb, err := repository.Get(id)
	if err != nil {
		return web.ErrInternalServerError
	}

	if sb == nil {
		return web.ErrNotFound
	}

	err = repository.Delete(id)
	if err != nil {
		return web.ErrInternalServerError
	}

	return nil
}
