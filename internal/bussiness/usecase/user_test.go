package usecase

import (
	"context"
	"testing"

	"github.com/gecoronel/donde-estan-ws/internal/bussiness/gateway"
	mock_gateway "github.com/gecoronel/donde-estan-ws/internal/bussiness/gateway/mocks"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model/web"
	ctx "github.com/gecoronel/donde-estan-ws/internal/infrastructure/delivery/api/context"
	"github.com/gecoronel/donde-estan-ws/internal/infrastructure/delivery/api/middleware/ioc"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	user = model.User{
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

	observedUser = model.ObservedUser{
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

	observerUser = model.ObserverUser{
		User: model.User{
			ID:        2,
			Name:      "Jose",
			LastName:  "Perez",
			IDNumber:  "12345678",
			Username:  "joseperez",
			Password:  "joseperez1234",
			Email:     "joseperez@mail.com",
			Enabled:   true,
			Type:      observer,
			CreatedAt: "2022-12-10 17:49:30",
			UpdatedAt: "2022-12-10 17:49:30",
		},
		ObservedUsers: []model.ObservedUser{observedUser},
	}
)

func TestUseCaseGet(t *testing.T) {
	observedUser := model.User{
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

	observerUser := observedUser
	observerUser.ID = 2
	observerUser.Type = observer

	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name          string
		mock          func() *mock_gateway.MockUserRepository
		input         uint64
		expectedUser  *model.User
		expectedError error
	}{
		{
			name: "error login in find by username",
			mock: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().Get(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockUserRepository
			},
			input:         observedUser.ID,
			expectedUser:  nil,
			expectedError: web.ErrInternalServerError,
		},
		{
			name: "not found error login for observed user",
			mock: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().Get(gomock.Any()).Return(nil, nil)
				return mockUserRepository
			},
			input:         observerUser.ID,
			expectedUser:  nil,
			expectedError: web.ErrNotFound,
		},
		{
			name: "successful login for observer user",
			mock: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().Get(gomock.Any()).Return(&observerUser, nil)
				return mockUserRepository
			},
			input:         observerUser.ID,
			expectedUser:  &observerUser,
			expectedError: nil,
		},
		{
			name: "successful login for observed user",
			mock: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().Get(gomock.Any()).Return(&observedUser, nil)
				return mockUserRepository
			},
			input:         observedUser.ID,
			expectedUser:  &observedUser,
			expectedError: nil,
		},
	}
	for _, test := range tests {
		mockUserRepository := test.mock()

		t.Run(test.name, func(t *testing.T) {
			context := getContextUser(mockUserRepository)
			serviceLocator := ctx.GetServiceLocator(context)
			uc := serviceLocator.GetInstance(UserUseCaseType).(UserUseCase)

			u, err := uc.Get(test.input, serviceLocator)

			assert.Equalf(t, test.expectedUser, u, "Expected user %v, received %v", test.expectedUser, u)
			assert.Equalf(t, test.expectedError, err, "Expected error %v, received %d", test.expectedError, err)
		})
	}
}

func TestUseCaseLogin(t *testing.T) {
	observedU := model.NewObservedUser(&observedUser)
	observerU := model.NewObserverUser(&observerUser)

	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name          string
		mock          func() *mock_gateway.MockUserRepository
		input         model.Login
		expectedUser  model.IUser
		expectedError error
	}{
		{
			name: "error login in find by username",
			mock: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().FindByUsername(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockUserRepository
			},
			input:         model.Login{Username: observedUser.User.Username, Password: observedUser.User.Password},
			expectedUser:  nil,
			expectedError: web.ErrInternalServerError,
		},
		{
			name: "not found error login for observed user",
			mock: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().FindByUsername(gomock.Any()).Return(nil, nil)
				return mockUserRepository
			},
			input:         model.Login{Username: observedUser.User.Username, Password: observedUser.User.Password},
			expectedUser:  nil,
			expectedError: web.ErrNotFound,
		},
		{
			name: "incorrect password in login for observer user",
			mock: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().FindByUsername(gomock.Any()).Return(&user, nil)
				return mockUserRepository
			},
			input:         model.Login{Username: observerUser.User.Username, Password: "incorrect"},
			expectedUser:  nil,
			expectedError: web.ErrIncorrectUserOrPassword,
		},
		{
			name: "error login for observed user",
			mock: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().FindByUsername(gomock.Any()).Return(&user, nil)
				mockUserRepository.EXPECT().GetObservedUser(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockUserRepository
			},
			input:         model.Login{Username: observedUser.User.Username, Password: observedUser.User.Password},
			expectedUser:  nil,
			expectedError: web.ErrInternalServerError,
		},
		{
			name: "error in login for observer user",
			mock: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().FindByUsername(gomock.Any()).Return(&observerUser.User, nil)
				mockUserRepository.EXPECT().GetObserverUser(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockUserRepository
			},
			input:         model.Login{Username: observerUser.User.Username, Password: observerUser.User.Password},
			expectedUser:  nil,
			expectedError: web.ErrInternalServerError,
		},
		{
			name: "successful login for observer user",
			mock: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().FindByUsername(gomock.Any()).Return(&observerUser.User, nil)
				mockUserRepository.EXPECT().GetObserverUser(gomock.Any()).Return(observerU, nil)
				return mockUserRepository
			},
			input:         model.Login{Username: observerUser.User.Username, Password: observerUser.User.Password},
			expectedUser:  observerU,
			expectedError: nil,
		},
		{
			name: "successful login for observed user",
			mock: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().FindByUsername(gomock.Any()).Return(&observedUser.User, nil)
				mockUserRepository.EXPECT().GetObservedUser(gomock.Any()).Return(observedU, nil)
				return mockUserRepository
			},
			input:         model.Login{Username: observedUser.User.Username, Password: observedUser.User.Password},
			expectedUser:  observedU,
			expectedError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUserRepository := test.mock()

			context := getContextUser(mockUserRepository)
			serviceLocator := ctx.GetServiceLocator(context)
			uc := serviceLocator.GetInstance(UserUseCaseType).(UserUseCase)

			u, err := uc.Login(test.input, serviceLocator)

			assert.Equalf(t, test.expectedUser, u, "Expected user %v, received %v", test.expectedUser, u)
			assert.Equalf(t, test.expectedError, err, "Expected error %v, received %d", test.expectedError, err)
		})
	}
}

func TestUseCaseCreateObservedUser(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name          string
		mock          func() *mock_gateway.MockUserRepository
		input         model.ObservedUser
		expectedUser  *model.ObservedUser
		expectedError error
	}{
		{
			name: "error create observed user",
			mock: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().SaveObservedUser(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockUserRepository
			},
			input:         observedUser,
			expectedUser:  nil,
			expectedError: web.ErrInternalServerError,
		},
		{
			name: "successful create observed user",
			mock: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().SaveObservedUser(gomock.Any()).Return(&observedUser, nil)
				return mockUserRepository
			},
			input:         observedUser,
			expectedUser:  &observedUser,
			expectedError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUserRepository := test.mock()

			context := getContextUser(mockUserRepository)
			serviceLocator := ctx.GetServiceLocator(context)
			uc := serviceLocator.GetInstance(UserUseCaseType).(UserUseCase)

			u, err := uc.CreateObservedUser(test.input, serviceLocator)

			assert.Equalf(t, test.expectedUser, u, "Expected user %v, received %v", test.expectedUser, u)
			assert.Equalf(t, test.expectedError, err, "Expected error %v, received %d", test.expectedError, err)
		})
	}
}

func TestUseCaseCreateObserverUser(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name          string
		mock          func() *mock_gateway.MockUserRepository
		input         model.ObserverUser
		expectedUser  *model.ObserverUser
		expectedError error
	}{
		{
			name: "error create observed user",
			mock: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().SaveObserverUser(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockUserRepository
			},
			input:         observerUser,
			expectedUser:  nil,
			expectedError: web.ErrInternalServerError,
		},
		{
			name: "successful create observed user",
			mock: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().SaveObserverUser(gomock.Any()).Return(&observerUser, nil)
				return mockUserRepository
			},
			input:         observerUser,
			expectedUser:  &observerUser,
			expectedError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUserRepository := test.mock()

			context := getContextUser(mockUserRepository)
			serviceLocator := ctx.GetServiceLocator(context)
			uc := serviceLocator.GetInstance(UserUseCaseType).(UserUseCase)

			u, err := uc.CreateObserverUser(test.input, serviceLocator)

			assert.Equalf(t, test.expectedUser, u, "Expected user %v, received %v", test.expectedUser, u)
			assert.Equalf(t, test.expectedError, err, "Expected error %v, received %d", test.expectedError, err)
		})
	}
}

func TestUseCaseFindByUsername(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name          string
		mock          func() *mock_gateway.MockUserRepository
		input         string
		expectedUser  *model.User
		expectedError error
	}{
		{
			name: "error create observed user",
			mock: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().FindByUsername(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockUserRepository
			},
			input:         user.Username,
			expectedUser:  nil,
			expectedError: web.ErrInternalServerError,
		},
		{
			name: "successful create observed user",
			mock: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().FindByUsername(gomock.Any()).Return(&user, nil)
				return mockUserRepository
			},
			input:         user.Username,
			expectedUser:  &user,
			expectedError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUserRepository := test.mock()

			context := getContextUser(mockUserRepository)
			serviceLocator := ctx.GetServiceLocator(context)
			uc := serviceLocator.GetInstance(UserUseCaseType).(UserUseCase)

			u, err := uc.FindByUsername(test.input, serviceLocator)

			assert.Equalf(t, test.expectedUser, u, "Expected user %v, received %v", test.expectedUser, u)
			assert.Equalf(t, test.expectedError, err, "Expected error %v, received %d", test.expectedError, err)
		})
	}
}

func TestUseCaseFindByEmail(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name          string
		mock          func() *mock_gateway.MockUserRepository
		input         string
		expectedUser  *model.User
		expectedError error
	}{
		{
			name: "error finding user by email",
			mock: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().FindByEmail(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockUserRepository
			},
			input:         user.Email,
			expectedUser:  nil,
			expectedError: web.ErrInternalServerError,
		},
		{
			name: "successful find user by email",
			mock: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().FindByEmail(gomock.Any()).Return(&user, nil)
				return mockUserRepository
			},
			input:         user.Email,
			expectedUser:  &user,
			expectedError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUserRepository := test.mock()

			context := getContextUser(mockUserRepository)
			serviceLocator := ctx.GetServiceLocator(context)
			uc := serviceLocator.GetInstance(UserUseCaseType).(UserUseCase)

			u, err := uc.FindByEmail(test.input, serviceLocator)

			assert.Equalf(t, test.expectedUser, u, "Expected user %v, received %v", test.expectedUser, u)
			assert.Equalf(t, test.expectedError, err, "Expected error %v, received %d", test.expectedError, err)
		})
	}
}

func TestUseCaseAddObservedUserInObserverUser(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name           string
		mock           func() *mock_gateway.MockUserRepository
		privacyKey     string
		observerUserID uint64
		expectedError  error
	}{
		{
			name: "error creating observed user in observer user",
			mock: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().FindObservedUserByPrivacyKey(gomock.Any()).Return(nil,
					web.ErrInternalServerError)
				return mockUserRepository
			},
			privacyKey:     observedUser.PrivacyKey,
			observerUserID: 2,
			expectedError:  web.ErrInternalServerError,
		},
		{
			name: "not found error creating observed user in observer user",
			mock: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().FindObservedUserByPrivacyKey(gomock.Any()).Return(nil, nil)
				return mockUserRepository
			},
			privacyKey:     observedUser.PrivacyKey,
			observerUserID: 2,
			expectedError:  web.ErrNotFound,
		},
		{
			name: "conflict error creating observed user in observer user",
			mock: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().FindObservedUserByPrivacyKey(gomock.Any()).Return(&observedUser, nil)
				mockUserRepository.EXPECT().GetObserverUser(gomock.Any()).Return(&observerUser, nil)
				return mockUserRepository
			},
			privacyKey:     observedUser.PrivacyKey,
			observerUserID: 2,
			expectedError:  web.ErrConflict,
		},
		{
			name: "error getting observer user",
			mock: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().FindObservedUserByPrivacyKey(gomock.Any()).Return(&observedUser, nil)
				mockUserRepository.EXPECT().GetObserverUser(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockUserRepository
			},
			privacyKey:     observedUser.PrivacyKey,
			observerUserID: 2,
			expectedError:  web.ErrInternalServerError,
		},
		{
			name: "not found error getting observer user",
			mock: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().FindObservedUserByPrivacyKey(gomock.Any()).Return(&observedUser, nil)
				mockUserRepository.EXPECT().GetObserverUser(gomock.Any()).Return(nil, nil)
				return mockUserRepository
			},
			privacyKey:     observedUser.PrivacyKey,
			observerUserID: 2,
			expectedError:  web.ErrNotFound,
		},
		{
			name: "unsuccessful create observed user in observer user",
			mock: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().FindObservedUserByPrivacyKey(gomock.Any()).Return(&observedUser, nil)
				mockUserRepository.EXPECT().GetObserverUser(gomock.Any()).Return(&model.ObserverUser{User: user}, nil)
				mockUserRepository.EXPECT().SaveObservedUserInObserverUser(gomock.Any(), gomock.Any()).
					Return(web.ErrInternalServerError)
				return mockUserRepository
			},
			privacyKey:     observedUser.PrivacyKey,
			observerUserID: 2,
			expectedError:  web.ErrInternalServerError,
		},
		{
			name: "successful create observed user in observer user",
			mock: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().FindObservedUserByPrivacyKey(gomock.Any()).Return(&observedUser, nil)
				mockUserRepository.EXPECT().GetObserverUser(gomock.Any()).Return(&model.ObserverUser{User: user}, nil)
				mockUserRepository.EXPECT().SaveObservedUserInObserverUser(gomock.Any(), gomock.Any()).Return(nil)
				return mockUserRepository
			},
			privacyKey:     observedUser.PrivacyKey,
			observerUserID: 2,
			expectedError:  nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUserRepository := test.mock()

			context := getContextUser(mockUserRepository)
			serviceLocator := ctx.GetServiceLocator(context)
			uc := serviceLocator.GetInstance(UserUseCaseType).(UserUseCase)

			err := uc.AddObservedUserInObserverUser(test.privacyKey, test.observerUserID, serviceLocator)

			assert.Equalf(t, test.expectedError, err, "Expected error %v, received %d", test.expectedError, err)
		})
	}
}

func getContextUser(mock *mock_gateway.MockUserRepository) context.Context {
	iocContext := ioc.NewContext()
	iocContext.Bind(gateway.UserRepositoryType).ToInstance(mock)
	iocContext.Bind(UserUseCaseType).ToInstance(NewUserUseCase())
	injector := ioc.NewInjector(iocContext)
	ctx := ctx.SetServiceLocator(context.TODO(), injector)

	return ctx
}
