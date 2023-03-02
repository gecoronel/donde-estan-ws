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
		FindByUsername(string, gateway.ServiceLocator) (*model.User, error)
		FindByEmail(string, gateway.ServiceLocator) (*model.User, error)
		CreateObservedUser(model.ObservedUser, gateway.ServiceLocator) (*model.ObservedUser, error)
		CreateObserverUser(model.ObserverUser, gateway.ServiceLocator) (*model.ObserverUser, error)
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
		return nil, web.ErrIncorrectUserOrPassword
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
		iUser, err = repository.GetObservedUser(odu.User.ID)
	case observer:
		iUser, err = repository.GetObserverUser(oru.User.ID)
	}

	if err != nil {
		return nil, web.ErrInternalServerError
	}

	return iUser, nil
}

func (u userUseCase) CreateObservedUser(observed model.ObservedUser, locator gateway.ServiceLocator) (
	*model.ObservedUser,
	error,
) {
	repository := locator.GetInstance(gateway.UserRepositoryType).(gateway.UserRepository)

	user, err := repository.SaveObservedUser(observed)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u userUseCase) CreateObserverUser(observer model.ObserverUser, locator gateway.ServiceLocator) (
	*model.ObserverUser,
	error,
) {
	repository := locator.GetInstance(gateway.UserRepositoryType).(gateway.UserRepository)

	user, err := repository.SaveObserverUser(observer)
	if err != nil {
		return nil, web.ErrInternalServerError
	}

	return user, nil
}

func (u userUseCase) FindByUsername(username string, locator gateway.ServiceLocator) (*model.User, error) {
	repository := locator.GetInstance(gateway.UserRepositoryType).(gateway.UserRepository)

	user, err := repository.FindByUsername(username)
	if err != nil {
		return nil, web.ErrInternalServerError
	}

	return user, nil
}

func (u userUseCase) FindByEmail(email string, locator gateway.ServiceLocator) (*model.User, error) {
	repository := locator.GetInstance(gateway.UserRepositoryType).(gateway.UserRepository)

	user, err := repository.FindByEmail(email)
	if err != nil {
		return nil, web.ErrInternalServerError
	}

	return user, nil
}
