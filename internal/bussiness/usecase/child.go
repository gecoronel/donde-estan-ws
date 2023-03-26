//go:generate mockgen --source=child.go --destination=././mocks/child.go

package usecase

import (
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/gateway"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model/web"
	log "github.com/sirupsen/logrus"
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

func (buc childUseCase) Get(id uint64, locator gateway.ServiceLocator) (*model.Child, error) {
	repository := locator.GetInstance(gateway.ChildRepositoryType).(gateway.ChildRepository)

	child, err := repository.Get(id)
	if err != nil {
		return nil, web.ErrInternalServerError
	}

	if child == nil {
		return nil, web.ErrNotFound
	}

	return child, nil
}

func (buc childUseCase) Save(child model.Child, locator gateway.ServiceLocator) (
	*model.Child,
	error,
) {
	userRepository := locator.GetInstance(gateway.UserRepositoryType).(gateway.UserRepository)
	childRepository := locator.GetInstance(gateway.ChildRepositoryType).(gateway.ChildRepository)

	u, err := userRepository.Get(child.ObserverUserID)
	if err != nil {
		return nil, err
	}
	if u == nil || u.Type != observer {
		log.Error("incorrect observer user for save child")
		return nil, web.ErrNotFound
	}

	c, err := childRepository.Save(child)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (buc childUseCase) Update(child model.Child, locator gateway.ServiceLocator) (
	*model.Child,
	error,
) {
	repository := locator.GetInstance(gateway.ChildRepositoryType).(gateway.ChildRepository)

	c, err := repository.Get(child.ID)
	if err != nil {
		return nil, web.ErrInternalServerError
	}

	if c == nil {
		return nil, web.ErrNotFound
	}

	c, err = repository.Update(child)
	if err != nil {
		return nil, web.ErrInternalServerError
	}

	return c, nil
}

func (buc childUseCase) Delete(id uint64, locator gateway.ServiceLocator) error {
	repository := locator.GetInstance(gateway.ChildRepositoryType).(gateway.ChildRepository)

	c, err := repository.Get(id)
	if err != nil {
		return web.ErrInternalServerError
	}

	if c == nil {
		return web.ErrNotFound
	}

	err = repository.Delete(id)
	if err != nil {
		return web.ErrInternalServerError
	}

	return nil
}
