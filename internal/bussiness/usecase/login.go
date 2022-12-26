package usecase

import (
	"github.com/gcoron/donde-estan-ws/internal/bussiness/gateway"
	"github.com/gcoron/donde-estan-ws/internal/bussiness/model"
	"github.com/gcoron/donde-estan-ws/internal/bussiness/model/web"
	"strings"
)

const (
	LoginUseCaseType = "LoginUseCase"
	observed         = "observed"
	observer         = "observer"
)

type (
	LoginUseCase interface {
		Login(model.Login, gateway.ServiceLocator) (*model.IUser, error)
	}

	loginUseCase struct{}
)

func NewLoginUseCase() LoginUseCase {
	return &loginUseCase{}
}

func (l loginUseCase) Login(login model.Login, locator gateway.ServiceLocator) (*model.IUser, error) {
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

	var (
		u   *model.IUser
		odu model.ObservedUser
		oru model.ObserverUser
	)

	odu.User = *user
	oru.User = *user

	switch user.Type {
	case observed:
		u, err = repository.GetObservedUser(&odu)
	case observer:
		u, err = repository.GetObserverUser(&oru)
	}

	if err != nil {
		return nil, err
	}

	return u, nil
}
