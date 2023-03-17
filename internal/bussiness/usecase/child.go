//go:generate mockgen --source=child.go --destination=././mocks/child.go

package usecase

import (
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/gateway"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model/web"
)

const ChildUseCaseType = "ChildUseCase"

type (
	ChildUseCase interface {
		Get(uint64, gateway.ServiceLocator) (*model.Child, error)
		Save(model.Child, gateway.ServiceLocator) (*model.Child, error)
		Update(model.Child, gateway.ServiceLocator) (*model.Child, error)
		Delete(uint64, gateway.ServiceLocator) error
	}

	childUseCase struct{}
)

func NewChildUseCase() ChildUseCase {
	return &childUseCase{}
}

func (b childUseCase) Get(id uint64, locator gateway.ServiceLocator) (*model.Child, error) {
	repository := locator.GetInstance(gateway.ChildRepositoryType).(gateway.ChildRepository)

	Child, err := repository.Get(id)
	if err != nil {
		return nil, web.ErrInternalServerError
	}

	if Child == nil {
		return nil, web.ErrNotFound
	}

	return Child, nil
}

func (b childUseCase) Save(child model.Child, locator gateway.ServiceLocator) (
	*model.Child,
	error,
) {
	repository := locator.GetInstance(gateway.ChildRepositoryType).(gateway.ChildRepository)

	sb, err := repository.Save(child)
	if err != nil {
		return nil, err
	}

	return sb, nil
}

func (b childUseCase) Update(Child model.Child, locator gateway.ServiceLocator) (
	*model.Child,
	error,
) {
	repository := locator.GetInstance(gateway.ChildRepositoryType).(gateway.ChildRepository)

	sb, err := repository.Get(Child.ID)
	if err != nil {
		return nil, web.ErrInternalServerError
	}

	if sb == nil {
		return nil, web.ErrNotFound
	}

	sb, err = repository.Update(Child)
	if err != nil {
		return nil, web.ErrInternalServerError
	}

	return sb, nil
}

func (b childUseCase) Delete(id uint64, locator gateway.ServiceLocator) error {
	repository := locator.GetInstance(gateway.ChildRepositoryType).(gateway.ChildRepository)

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
