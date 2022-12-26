package usecase

import (
	"github.com/gcoron/donde-estan-ws/internal/bussiness/gateway"
	"github.com/gcoron/donde-estan-ws/internal/bussiness/model"
	"github.com/gcoron/donde-estan-ws/internal/bussiness/model/web"
)

const UserUseCaseType = "UserUseCase"

type (
	UserUseCase interface {
		Get(uint, gateway.ServiceLocator) (*model.User, error)
	}

	userUseCase struct{}
)

func NewUserUseCase() UserUseCase {
	return &userUseCase{}
}

func (u userUseCase) Get(userID uint, locator gateway.ServiceLocator) (*model.User, error) {
	repository := locator.GetInstance(gateway.UserRepositoryType).(gateway.UserRepository)
	user, err := repository.Get(userID)
	if err != nil {
		return nil, web.ErrInternalServerError
	}

	if user == nil {
		return nil, web.ErrNotFound
	}

	return user, nil
}
