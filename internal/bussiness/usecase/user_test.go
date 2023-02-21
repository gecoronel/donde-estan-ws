package usecase

import (
	"context"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/gateway"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model/web"
	"github.com/gecoronel/donde-estan-ws/internal/infrastructure/delivery/api/middleware/ioc"
	"testing"

	mock_gateway "github.com/gecoronel/donde-estan-ws/internal/bussiness/gateway/mocks"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model"
	ctx "github.com/gecoronel/donde-estan-ws/internal/infrastructure/delivery/api/context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUseCaseGet(t *testing.T) {
	user := model.User{
		ID:        1,
		Name:      "Juan",
		LastName:  "Perez",
		IDNumber:  "12345678",
		Username:  "jperez",
		Password:  "jperez1234",
		Email:     "jperez@mail.com",
		Enabled:   true,
		Type:      observed,
		CreatedAt: "2022-12-10 17:49:30",
		UpdatedAt: "2022-12-10 17:49:30",
	}

	m := gomock.NewController(t)
	defer m.Finish()

	t.Run("error login in find by username", func(t *testing.T) {
		mockUserRepository := mock_gateway.NewMockUserRepository(m)
		mockUserRepository.EXPECT().Get(gomock.Any()).Return(nil, web.ErrInternalServerError)

		context := getContext(mockUserRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(UserUseCaseType).(UserUseCase)

		u, err := uc.Get(1, serviceLocator)

		assert.Nil(t, u)
		assert.Equalf(t, web.ErrInternalServerError, err, "Expected error %v, received %d", nil, err)
	})

	t.Run("successful login for observer user", func(t *testing.T) {
		mockUserRepository := mock_gateway.NewMockUserRepository(m)
		user.Type = observer
		mockUserRepository.EXPECT().Get(gomock.Any()).Return(nil, nil)

		context := getContext(mockUserRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(UserUseCaseType).(UserUseCase)

		u, err := uc.Get(1, serviceLocator)

		assert.Nil(t, u)
		assert.Equalf(t, web.ErrNotFound, err, "Expected error %v, received %d", nil, err)
	})

	t.Run("successful login for observer user", func(t *testing.T) {
		mockUserRepository := mock_gateway.NewMockUserRepository(m)
		user.Type = observer
		mockUserRepository.EXPECT().Get(gomock.Any()).Return(&user, nil)

		context := getContext(mockUserRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(UserUseCaseType).(UserUseCase)

		u, err := uc.Get(1, serviceLocator)

		assert.Equalf(t, user, *u, "Expected response %d, received %d", user, u)
		assert.Equalf(t, nil, err, "Expected error %v, received %d", nil, err)
	})
}

func TestUseCaseLogin(t *testing.T) {
	login := model.Login{
		Username: "jperez",
		Password: "jperez1234",
	}

	user := model.User{
		ID:        1,
		Name:      "Juan",
		LastName:  "Perez",
		IDNumber:  "12345678",
		Username:  "jperez",
		Password:  "jperez1234",
		Email:     "jperez@mail.com",
		Enabled:   true,
		Type:      observed,
		CreatedAt: "2022-12-10 17:49:30",
		UpdatedAt: "2022-12-10 17:49:30",
	}

	observedUser := model.ObservedUser{
		User:        user,
		PrivacyKey:  "juan.perez.12345678",
		CompanyName: "school bus",
		SchoolBus: model.SchoolBus{
			ID:           "1",
			LicensePlate: "11AAA222",
			Model:        "Master",
			Brand:        "Renault",
			License:      "11222",
			CreatedAt:    "2022-12-10 17:49:30",
			UpdatedAt:    "2022-12-10 17:49:30",
		},
		ObserverUsers: nil,
	}

	observerUser := model.ObserverUser{
		User: user,
	}

	observedU := model.NewObservedUser(observedUser)
	observerU := model.NewObserverUser(observerUser)

	m := gomock.NewController(t)
	defer m.Finish()

	t.Run("error login in find by username", func(t *testing.T) {
		mockUserRepository := mock_gateway.NewMockUserRepository(m)
		mockUserRepository.EXPECT().FindByUsername(gomock.Any()).Return(nil, web.ErrInternalServerError)

		context := getContext(mockUserRepository)

		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(UserUseCaseType).(UserUseCase)

		uod, err := uc.Login(login, serviceLocator)

		assert.Equalf(t, nil, uod, "Expected response %d, received %d", observedU, uod)
		assert.Equalf(t, web.ErrInternalServerError, err, "Expected error %v, received %d", nil, err)
	})

	t.Run("not found error login in find by username", func(t *testing.T) {
		mockUserRepository := mock_gateway.NewMockUserRepository(m)
		mockUserRepository.EXPECT().FindByUsername(gomock.Any()).Return(nil, nil)

		context := getContext(mockUserRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(UserUseCaseType).(UserUseCase)

		uod, err := uc.Login(login, serviceLocator)

		assert.Equalf(t, nil, uod, "Expected response %d, received %d", observedU, uod)
		assert.Equalf(t, web.ErrNotFound, err, "Expected error %v, received %d", nil, err)
	})

	t.Run("error login in find by username for incorrect password", func(t *testing.T) {
		mockUserRepository := mock_gateway.NewMockUserRepository(m)
		user.Password = "incorrect"
		mockUserRepository.EXPECT().FindByUsername(gomock.Any()).Return(&user, nil)

		context := getContext(mockUserRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(UserUseCaseType).(UserUseCase)

		uod, err := uc.Login(login, serviceLocator)

		assert.Equalf(t, nil, uod, "Expected response %d, received %d", observedU, uod)
		assert.Equalf(t, web.ErrIncorrectPassword, err, "Expected error %v, received %d", nil, err)
		user.Password = "jperez1234"
	})

	t.Run("error login for observed user", func(t *testing.T) {
		mockUserRepository := mock_gateway.NewMockUserRepository(m)
		mockUserRepository.EXPECT().FindByUsername(gomock.Any()).Return(&user, nil)
		mockUserRepository.EXPECT().GetObservedUser(gomock.Any()).Return(nil, web.ErrInternalServerError)

		context := getContext(mockUserRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(UserUseCaseType).(UserUseCase)

		uod, err := uc.Login(login, serviceLocator)

		assert.Equalf(t, nil, uod, "Expected response %d, received %d", observedU, uod)
		assert.Equalf(t, web.ErrInternalServerError, err, "Expected error %v, received %d", nil, err)
	})

	t.Run("successful login for observed user", func(t *testing.T) {
		mockUserRepository := mock_gateway.NewMockUserRepository(m)
		mockUserRepository.EXPECT().GetObservedUser(gomock.Any()).Return(observedU, nil)
		mockUserRepository.EXPECT().FindByUsername(gomock.Any()).Return(&user, nil)

		context := getContext(mockUserRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(UserUseCaseType).(UserUseCase)

		uod, err := uc.Login(login, serviceLocator)

		assert.Equalf(t, observedU, uod, "Expected response %d, received %d", observedU, uod)
		assert.Equalf(t, nil, err, "Expected error %v, received %d", nil, err)
	})

	t.Run("successful login for observer user", func(t *testing.T) {
		mockUserRepository := mock_gateway.NewMockUserRepository(m)
		user.Type = observer
		mockUserRepository.EXPECT().FindByUsername(gomock.Any()).Return(&user, nil)
		mockUserRepository.EXPECT().GetObserverUser(gomock.Any()).Return(observerU, nil)

		context := getContext(mockUserRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(UserUseCaseType).(UserUseCase)

		uor, err := uc.Login(login, serviceLocator)

		assert.Equalf(t, observerU, uor, "Expected response %d, received %d", observedU, uor)
		assert.Equalf(t, nil, err, "Expected error %v, received %d", nil, err)
	})
}

func TestUseCaseCreateObservedUser(t *testing.T) {
	var observedUser = model.ObservedUser{
		User: model.User{
			ID:        1,
			Name:      "Juan",
			LastName:  "Perez",
			IDNumber:  "12345678",
			Username:  "jperez",
			Password:  "jperez1234",
			Email:     "jperez@mail.com",
			Enabled:   true,
			Type:      "observed",
			CreatedAt: "2022-12-10 17:49:30",
			UpdatedAt: "2022-12-10 17:49:30",
		},
		PrivacyKey:  "juan.perez.12345678",
		CompanyName: "school bus company",
		SchoolBus: model.SchoolBus{
			ID:           "1",
			LicensePlate: "11AAA22",
			Model:        "Master",
			Brand:        "Renault",
			License:      "111",
			CreatedAt:    "2023-02-18 17:09:33",
			UpdatedAt:    "2023-02-18 17:09:33",
		},
	}

	m := gomock.NewController(t)
	defer m.Finish()

	t.Run("error create observed user", func(t *testing.T) {
		mockUserRepository := mock_gateway.NewMockUserRepository(m)
		mockUserRepository.EXPECT().SaveObservedUser(gomock.Any()).Return(nil, web.ErrInternalServerError)

		context := getContext(mockUserRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(UserUseCaseType).(UserUseCase)

		u, err := uc.CreateObservedUser(observedUser, serviceLocator)

		assert.Nil(t, u)
		assert.Equalf(t, web.ErrInternalServerError, err, "Expected error %v, received %d", nil, err)
	})

	t.Run("successful create observed user", func(t *testing.T) {
		mockUserRepository := mock_gateway.NewMockUserRepository(m)
		mockUserRepository.EXPECT().SaveObservedUser(gomock.Any()).Return(&observedUser, nil)

		context := getContext(mockUserRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(UserUseCaseType).(UserUseCase)

		u, err := uc.CreateObservedUser(observedUser, serviceLocator)

		assert.Equalf(t, observedUser, *u, "Expected response %v, received %v", observedUser, u)
		assert.Equalf(t, nil, err, "Expected error %v, received %d", nil, err)
	})
}

func TestUseCaseCreateObserverUser(t *testing.T) {
	var observerUser = model.ObserverUser{
		User: model.User{
			ID:        1,
			Name:      "Juan",
			LastName:  "Perez",
			IDNumber:  "12345678",
			Username:  "jperez",
			Password:  "jperez1234",
			Email:     "jperez@mail.com",
			Enabled:   true,
			Type:      "observed",
			CreatedAt: "2022-12-10 17:49:30",
			UpdatedAt: "2022-12-10 17:49:30",
		},
	}

	m := gomock.NewController(t)
	defer m.Finish()

	t.Run("error create observed user", func(t *testing.T) {
		mockUserRepository := mock_gateway.NewMockUserRepository(m)
		mockUserRepository.EXPECT().SaveObserverUser(gomock.Any()).Return(nil, web.ErrInternalServerError)

		context := getContext(mockUserRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(UserUseCaseType).(UserUseCase)

		u, err := uc.CreateObserverUser(observerUser, serviceLocator)

		assert.Nil(t, u)
		assert.Equalf(t, web.ErrInternalServerError, err, "Expected error %v, received %d", nil, err)
	})

	t.Run("successful create observed user", func(t *testing.T) {
		mockUserRepository := mock_gateway.NewMockUserRepository(m)
		mockUserRepository.EXPECT().SaveObserverUser(gomock.Any()).Return(&observerUser, nil)

		context := getContext(mockUserRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(UserUseCaseType).(UserUseCase)

		u, err := uc.CreateObserverUser(observerUser, serviceLocator)

		assert.Equalf(t, observerUser, *u, "Expected response %v, received %v", observerUser, u)
		assert.Equalf(t, nil, err, "Expected error %v, received %d", nil, err)
	})
}

func TestUseCaseFindByUsername(t *testing.T) {
	var user = model.User{
		ID:        1,
		Name:      "Juan",
		LastName:  "Perez",
		IDNumber:  "12345678",
		Username:  "jperez",
		Password:  "jperez1234",
		Email:     "jperez@mail.com",
		Enabled:   true,
		Type:      "observed",
		CreatedAt: "2022-12-10 17:49:30",
		UpdatedAt: "2022-12-10 17:49:30",
	}

	m := gomock.NewController(t)
	defer m.Finish()

	t.Run("error find user by username", func(t *testing.T) {
		mockUserRepository := mock_gateway.NewMockUserRepository(m)
		mockUserRepository.EXPECT().FindByUsername(gomock.Any()).Return(nil, web.ErrInternalServerError)

		context := getContext(mockUserRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(UserUseCaseType).(UserUseCase)

		u, err := uc.FindByUsername(user.Username, serviceLocator)

		assert.Nil(t, u)
		assert.Equalf(t, web.ErrInternalServerError, err, "Expected error %v, received %d", nil, err)
	})

	t.Run("successful find user by username", func(t *testing.T) {
		mockUserRepository := mock_gateway.NewMockUserRepository(m)
		mockUserRepository.EXPECT().FindByUsername(gomock.Any()).Return(&user, nil)

		context := getContext(mockUserRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(UserUseCaseType).(UserUseCase)

		u, err := uc.FindByUsername(user.Username, serviceLocator)

		assert.Equalf(t, user, *u, "Expected response %v, received %v", user, u)
		assert.Equalf(t, nil, err, "Expected error %v, received %d", nil, err)
	})
}

func TestUseCaseFindByEmail(t *testing.T) {
	var user = model.User{
		ID:        1,
		Name:      "Juan",
		LastName:  "Perez",
		IDNumber:  "12345678",
		Username:  "jperez",
		Password:  "jperez1234",
		Email:     "jperez@mail.com",
		Enabled:   true,
		Type:      "observed",
		CreatedAt: "2022-12-10 17:49:30",
		UpdatedAt: "2022-12-10 17:49:30",
	}

	m := gomock.NewController(t)
	defer m.Finish()

	t.Run("error find user by email", func(t *testing.T) {
		mockUserRepository := mock_gateway.NewMockUserRepository(m)
		mockUserRepository.EXPECT().FindByEmail(gomock.Any()).Return(nil, web.ErrInternalServerError)

		context := getContext(mockUserRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(UserUseCaseType).(UserUseCase)

		u, err := uc.FindByEmail(user.Email, serviceLocator)

		assert.Nil(t, u)
		assert.Equalf(t, web.ErrInternalServerError, err, "Expected error %v, received %d", nil, err)
	})

	t.Run("successful find user by email", func(t *testing.T) {
		mockUserRepository := mock_gateway.NewMockUserRepository(m)
		mockUserRepository.EXPECT().FindByEmail(gomock.Any()).Return(&user, nil)

		context := getContext(mockUserRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(UserUseCaseType).(UserUseCase)

		u, err := uc.FindByEmail(user.Email, serviceLocator)

		assert.Equalf(t, user, *u, "Expected response %v, received %v", user, u)
		assert.Equalf(t, nil, err, "Expected error %v, received %d", nil, err)
	})
}

func getContext(mock *mock_gateway.MockUserRepository) context.Context {
	iocContext := ioc.NewContext()
	iocContext.Bind(gateway.UserRepositoryType).ToInstance(mock)
	iocContext.Bind(UserUseCaseType).ToInstance(NewUserUseCase())
	injector := ioc.NewInjector(iocContext)
	ctx := ctx.SetServiceLocator(context.TODO(), injector)

	return ctx
}
