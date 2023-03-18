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
		UpdateObservedUser(model.ObservedUser, gateway.ServiceLocator) (*model.ObservedUser, error)
		UpdateObserverUser(model.ObserverUser, gateway.ServiceLocator) (*model.ObserverUser, error)
		DeleteObservedUser(uint64, gateway.ServiceLocator) error
		DeleteObserverUser(uint64, gateway.ServiceLocator) error
		AddObservedUserInObserverUser(string, uint64, gateway.ServiceLocator) error
		DeleteObservedUserInObserverUser(uint64, uint64, gateway.ServiceLocator) error
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
		odu   *model.ObservedUser
		oru   *model.ObserverUser
	)

	switch user.Type {
	case observed:
		odu, err = repository.GetObservedUser(user.ID)
		iUser = model.NewObservedUser(odu)
	case observer:
		oru, err = repository.GetObserverUser(user.ID)
		iUser = model.NewObserverUser(oru)
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

	odu, err := repository.FindByUsername(observed.User.Username)
	if err != nil {
		return nil, err
	}
	if odu != nil {
		return nil, web.ErrConflict
	}

	odu, err = repository.FindByEmail(observed.User.Email)
	if err != nil {
		return nil, err
	}
	if odu != nil {
		return nil, web.ErrConflict
	}

	user, err := repository.SaveObservedUser(observed)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u userUseCase) UpdateObservedUser(observed model.ObservedUser, locator gateway.ServiceLocator) (
	*model.ObservedUser,
	error,
) {
	repository := locator.GetInstance(gateway.UserRepositoryType).(gateway.UserRepository)

	odu, err := repository.Get(observed.User.ID)
	if err != nil {
		return nil, err
	}
	if odu == nil {
		return nil, web.ErrNotFound
	}

	user, err := repository.UpdateObservedUser(observed)
	if err != nil {
		return nil, web.ErrInternalServerError
	}

	return user, nil
}

func (u userUseCase) DeleteObservedUser(id uint64, locator gateway.ServiceLocator) error {
	repository := locator.GetInstance(gateway.UserRepositoryType).(gateway.UserRepository)

	user, err := repository.Get(id)
	if err != nil {
		return web.ErrInternalServerError
	}

	if user == nil {
		return web.ErrNotFound
	}

	err = repository.DeleteObservedUser(id)
	if err != nil {
		return web.ErrInternalServerError
	}

	return nil
}

func (u userUseCase) DeleteObserverUser(id uint64, locator gateway.ServiceLocator) error {
	repository := locator.GetInstance(gateway.UserRepositoryType).(gateway.UserRepository)

	user, err := repository.Get(id)
	if err != nil {
		return web.ErrInternalServerError
	}

	if user == nil {
		return web.ErrNotFound
	}

	err = repository.DeleteObserverUser(id)
	if err != nil {
		return web.ErrInternalServerError
	}

	return nil
}

func (u userUseCase) CreateObserverUser(observer model.ObserverUser, locator gateway.ServiceLocator) (
	*model.ObserverUser,
	error,
) {
	repository := locator.GetInstance(gateway.UserRepositoryType).(gateway.UserRepository)

	oru, err := repository.FindByUsername(observer.User.Username)
	if err != nil {
		return nil, err
	}
	if oru != nil {
		return nil, web.ErrConflict
	}

	oru, err = repository.FindByEmail(observer.User.Email)
	if err != nil {
		return nil, err
	}
	if oru != nil {
		return nil, web.ErrConflict
	}

	user, err := repository.SaveObserverUser(observer)
	if err != nil {
		return nil, web.ErrInternalServerError
	}

	return user, nil
}

func (u userUseCase) UpdateObserverUser(observer model.ObserverUser, locator gateway.ServiceLocator) (
	*model.ObserverUser,
	error,
) {
	repository := locator.GetInstance(gateway.UserRepositoryType).(gateway.UserRepository)

	oru, err := repository.Get(observer.User.ID)
	if err != nil {
		return nil, err
	}
	if oru == nil {
		return nil, web.ErrNotFound
	}

	user, err := repository.UpdateObserverUser(observer)
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

func (u userUseCase) AddObservedUserInObserverUser(
	privacyKey string,
	observerUserID uint64,
	locator gateway.ServiceLocator,
) error {
	repository := locator.GetInstance(gateway.UserRepositoryType).(gateway.UserRepository)

	odu, err := repository.FindObservedUserByPrivacyKey(privacyKey)
	if err != nil {
		return web.ErrInternalServerError
	}
	if odu == nil {
		return web.ErrNotFound
	}

	oru, err := repository.GetObserverUser(observerUserID)
	if err != nil {
		return web.ErrInternalServerError
	}
	if oru == nil {
		return web.ErrNotFound
	}
	if existObservedUserInObserverUser(*oru, odu.User.ID) {
		return web.ErrConflict
	}

	err = repository.SaveObservedUserInObserverUser(odu.User.ID, observerUserID)
	if err != nil {
		return web.ErrInternalServerError
	}

	return nil
}

func (u userUseCase) DeleteObservedUserInObserverUser(
	observedUserID uint64,
	observerUserID uint64,
	locator gateway.ServiceLocator,
) error {
	repository := locator.GetInstance(gateway.UserRepositoryType).(gateway.UserRepository)

	err := repository.DeleteObservedUserInObserverUser(observedUserID, observerUserID)
	if err != nil {
		return web.ErrInternalServerError
	}

	return nil
}

func existObservedUserInObserverUser(observer model.ObserverUser, observedUserID uint64) bool {
	for _, odu := range observer.ObservedUsers {
		if odu.User.ID == observedUserID {
			return true
		}
	}

	return false
}
