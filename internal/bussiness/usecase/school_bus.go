//go:generate mockgen --source=school_bus.go --destination=././mocks/school_bus.go

package usecase

import (
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/gateway"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model/web"
)

const SchoolBusUseCaseType = "SchoolBusUseCase"

type (
	SchoolBusUseCase interface {
		Get(string, gateway.ServiceLocator) (*model.SchoolBus, error)
		Save(model.SchoolBus, gateway.ServiceLocator) (*model.SchoolBus, error)
		Update(model.SchoolBus, gateway.ServiceLocator) (*model.SchoolBus, error)
		Delete(string, gateway.ServiceLocator) error
	}

	schoolBusUseCase struct{}
)

func NewSchoolBusUseCase() SchoolBusUseCase {
	return &schoolBusUseCase{}
}

func (b schoolBusUseCase) Get(id string, locator gateway.ServiceLocator) (*model.SchoolBus, error) {
	repository := locator.GetInstance(gateway.SchoolBusRepositoryType).(gateway.SchoolBusRepository)

	schoolBus, err := repository.Get(id)
	if err != nil {
		return nil, web.ErrInternalServerError
	}

	if schoolBus == nil {
		return nil, web.ErrNotFound
	}

	return schoolBus, nil
}

func (b schoolBusUseCase) Save(schoolBus model.SchoolBus, locator gateway.ServiceLocator) (
	*model.SchoolBus,
	error,
) {
	repository := locator.GetInstance(gateway.SchoolBusRepositoryType).(gateway.SchoolBusRepository)

	sb, err := repository.Save(schoolBus)
	if err != nil {
		return nil, err
	}

	return sb, nil
}

func (b schoolBusUseCase) Update(schoolBus model.SchoolBus, locator gateway.ServiceLocator) (
	*model.SchoolBus,
	error,
) {
	repository := locator.GetInstance(gateway.SchoolBusRepositoryType).(gateway.SchoolBusRepository)

	sb, err := repository.Get(schoolBus.ID)
	if err != nil {
		return nil, web.ErrInternalServerError
	}

	if sb == nil {
		return nil, web.ErrNotFound
	}

	sb, err = repository.Update(schoolBus)
	if err != nil {
		return nil, web.ErrInternalServerError
	}

	return sb, nil
}

func (b schoolBusUseCase) Delete(id string, locator gateway.ServiceLocator) error {
	repository := locator.GetInstance(gateway.SchoolBusRepositoryType).(gateway.SchoolBusRepository)

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
