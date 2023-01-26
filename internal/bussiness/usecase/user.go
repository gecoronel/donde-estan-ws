//go:generate mockgen --source=user.go --destination=././mocks/user.go

package usecase

import (
	"strings"

	"github.com/gecoronel/donde-estan-ws/internal/bussiness/gateway"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model/web"
)

const (
	observed = "observed"
	observer = "observer"

	UserUseCaseType = "UserUseCase"
)

type (
	UserUseCase interface {
		Get(uint64, gateway.ServiceLocator) (*model.User, error)
		Login(model.Login, gateway.ServiceLocator) (model.IUser, error)
	}

	userUseCase struct{}
)

func NewUserUseCase() UserUseCase {
	return &userUseCase{}
}

func (u userUseCase) Get(userID uint64, locator gateway.ServiceLocator) (*model.User, error) {
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

func (u userUseCase) Login(login model.Login, locator gateway.ServiceLocator) (model.IUser, error) {
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
		iUser model.IUser
		odu   model.ObservedUser
		oru   model.ObserverUser
	)

	odu.User = *user
	oru.User = *user

	switch user.Type {
	case observed:
		iUser, err = repository.GetObservedUser(&odu)
	case observer:
		iUser, err = repository.GetObserverUser(&oru)
	}

	if err != nil {
		return nil, web.ErrInternalServerError
	}

	return iUser, nil
}
