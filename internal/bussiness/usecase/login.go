package usecase

import (
	"github.com/gcoron/donde-estan-ws/internal/bussiness/gateway"
	"github.com/gcoron/donde-estan-ws/internal/bussiness/model"
	"github.com/gcoron/donde-estan-ws/internal/bussiness/model/web"
	"strings"
)

const LoginUseCaseType = "LoginUseCase"

type (
	LoginUseCase interface {
		Login(model.Login, gateway.ServiceLocator) (*model.User, error)
	}

	loginUseCase struct{}
)

func NewLoginUseCase() LoginUseCase {
	return &loginUseCase{}
}

func (u loginUseCase) Login(login model.Login, locator gateway.ServiceLocator) (*model.User, error) {
	repository := locator.GetInstance(gateway.UserRepositoryType).(gateway.UserRepository)
	user, err := repository.FindByUsername(login.Username)
	if err != nil {
		return nil, web.ErrInternalServerError
	}

	if user == nil {
		return nil, web.ErrNotFound
	}

	if !strings.EqualFold(user.Password, login.Password) {
		return nil, web.ErrIncorrectPassword
	}

	return user, nil
}
