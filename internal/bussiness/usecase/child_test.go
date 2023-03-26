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

var c = model.Child{
	ID:              1,
	Name:            "Pilar",
	LastName:        "Dominguez",
	SchoolName:      "La Salle",
	SchoolStartTime: "8:00",
	SchoolEndTime:   "12:00",
	CreatedAt:       "2023-02-18 17:09:33",
	UpdatedAt:       "2023-02-18 17:09:33",
	ObserverUserID:  uint64(10),
}

func TestUseCaseGetChild(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name          string
		mockChild     func() *mock_gateway.MockChildRepository
		mockUser      func() *mock_gateway.MockUserRepository
		input         uint64
		expectedChild *model.Child
		expectedError error
	}{
		{
			name: "error getting child",
			mockUser: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				return mockUserRepository
			},
			mockChild: func() *mock_gateway.MockChildRepository {
				mockChildRepository := mock_gateway.NewMockChildRepository(m)
				mockChildRepository.EXPECT().Get(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockChildRepository
			},
			input:         c.ID,
			expectedChild: nil,
			expectedError: web.ErrInternalServerError,
		},
		{
			name: "not found error getting child",
			mockUser: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				return mockUserRepository
			},
			mockChild: func() *mock_gateway.MockChildRepository {
				mockChildRepository := mock_gateway.NewMockChildRepository(m)
				mockChildRepository.EXPECT().Get(gomock.Any()).Return(nil, nil)
				return mockChildRepository
			},
			input:         c.ID,
			expectedChild: nil,
			expectedError: web.ErrNotFound,
		},
		{
			name: "successful create child",
			mockUser: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				return mockUserRepository
			},
			mockChild: func() *mock_gateway.MockChildRepository {
				mockChildRepository := mock_gateway.NewMockChildRepository(m)
				mockChildRepository.EXPECT().Get(gomock.Any()).Return(&c, nil)
				return mockChildRepository
			},
			input:         c.ID,
			expectedChild: &c,
			expectedError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUserRepository := test.mockUser()
			mockChildRepository := test.mockChild()

			context := getContextChild(mockUserRepository, mockChildRepository)
			serviceLocator := ctx.GetServiceLocator(context)
			uc := serviceLocator.GetInstance(ChildUseCaseType).(ChildUseCase)

			u, err := uc.Get(test.input, serviceLocator)

			assert.Equalf(t, test.expectedChild, u, "Expected user %v, received %v", test.expectedChild, u)
			assert.Equalf(t, test.expectedError, err, "Expected error %v, received %d", test.expectedError, err)
		})
	}
}

func TestUseCaseSaveChild(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name          string
		mockUser      func() *mock_gateway.MockUserRepository
		mockChild     func() *mock_gateway.MockChildRepository
		input         model.Child
		expectedChild *model.Child
		expectedError error
	}{
		{
			name: "error saving child",
			mockUser: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().Get(gomock.Any()).Return(&observerUser.User, nil)
				return mockUserRepository
			},
			mockChild: func() *mock_gateway.MockChildRepository {
				mockChildRepository := mock_gateway.NewMockChildRepository(m)
				mockChildRepository.EXPECT().Save(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockChildRepository
			},
			input:         c,
			expectedChild: nil,
			expectedError: web.ErrInternalServerError,
		},
		{
			name: "error save child for not found observer user",
			mockUser: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().Get(gomock.Any()).Return(nil, nil)
				return mockUserRepository
			},
			mockChild: func() *mock_gateway.MockChildRepository {
				mockChildRepository := mock_gateway.NewMockChildRepository(m)
				return mockChildRepository
			},
			input:         c,
			expectedChild: nil,
			expectedError: web.ErrNotFound,
		},
		{
			name: "error save child for get observer user",
			mockUser: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().Get(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockUserRepository
			},
			mockChild: func() *mock_gateway.MockChildRepository {
				mockChildRepository := mock_gateway.NewMockChildRepository(m)
				return mockChildRepository
			},
			input:         c,
			expectedChild: nil,
			expectedError: web.ErrInternalServerError,
		},
		{
			name: "error save child for incorrect type of user",
			mockUser: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().Get(gomock.Any()).Return(&observedUser.User, nil)
				return mockUserRepository
			},
			mockChild: func() *mock_gateway.MockChildRepository {
				mockChildRepository := mock_gateway.NewMockChildRepository(m)
				return mockChildRepository
			},
			input:         c,
			expectedChild: nil,
			expectedError: web.ErrNotFound,
		},
		{
			name: "successful save child",
			mockUser: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				mockUserRepository.EXPECT().Get(gomock.Any()).Return(&observerUser.User, nil)
				return mockUserRepository
			},
			mockChild: func() *mock_gateway.MockChildRepository {
				mockChildRepository := mock_gateway.NewMockChildRepository(m)
				mockChildRepository.EXPECT().Save(gomock.Any()).Return(&c, nil)
				return mockChildRepository
			},
			input:         c,
			expectedChild: &c,
			expectedError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUserRepository := test.mockUser()
			mockChildRepository := test.mockChild()

			context := getContextChild(mockUserRepository, mockChildRepository)
			serviceLocator := ctx.GetServiceLocator(context)
			uc := serviceLocator.GetInstance(ChildUseCaseType).(ChildUseCase)

			u, err := uc.Save(test.input, serviceLocator)

			assert.Equalf(t, test.expectedChild, u, "Expected user %v, received %v", test.expectedChild, u)
			assert.Equalf(t, test.expectedError, err, "Expected error %v, received %d", test.expectedError, err)
		})
	}
}

func TestUseCaseUpdateChild(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name          string
		mockUser      func() *mock_gateway.MockUserRepository
		mockChild     func() *mock_gateway.MockChildRepository
		input         model.Child
		expectedChild *model.Child
		expectedError error
	}{
		{
			name: "error updating child",
			mockUser: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				return mockUserRepository
			},
			mockChild: func() *mock_gateway.MockChildRepository {
				mockChildRepository := mock_gateway.NewMockChildRepository(m)
				mockChildRepository.EXPECT().Get(gomock.Any()).Return(&c, nil)
				mockChildRepository.EXPECT().Update(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockChildRepository
			},
			input:         c,
			expectedChild: nil,
			expectedError: web.ErrInternalServerError,
		},
		{
			name: "not found error updating child",
			mockUser: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				return mockUserRepository
			},
			mockChild: func() *mock_gateway.MockChildRepository {
				mockChildRepository := mock_gateway.NewMockChildRepository(m)
				mockChildRepository.EXPECT().Get(gomock.Any()).Return(nil, nil)
				return mockChildRepository
			},
			input:         c,
			expectedChild: nil,
			expectedError: web.ErrNotFound,
		},
		{
			name: "error getting child",
			mockUser: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				return mockUserRepository
			},
			mockChild: func() *mock_gateway.MockChildRepository {
				mockChildRepository := mock_gateway.NewMockChildRepository(m)
				mockChildRepository.EXPECT().Get(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockChildRepository
			},
			input:         c,
			expectedChild: nil,
			expectedError: web.ErrInternalServerError,
		},
		{
			name: "successful update child",
			mockUser: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				return mockUserRepository
			},
			mockChild: func() *mock_gateway.MockChildRepository {
				mockChildRepository := mock_gateway.NewMockChildRepository(m)
				mockChildRepository.EXPECT().Get(gomock.Any()).Return(&c, nil)
				mockChildRepository.EXPECT().Update(gomock.Any()).Return(&c, nil)
				return mockChildRepository
			},
			input:         c,
			expectedChild: &c,
			expectedError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUserRepository := test.mockUser()
			mockChildRepository := test.mockChild()

			context := getContextChild(mockUserRepository, mockChildRepository)
			serviceLocator := ctx.GetServiceLocator(context)
			uc := serviceLocator.GetInstance(ChildUseCaseType).(ChildUseCase)

			u, err := uc.Update(test.input, serviceLocator)

			assert.Equalf(t, test.expectedChild, u, "Expected user %v, received %v", test.expectedChild, u)
			assert.Equalf(t, test.expectedError, err, "Expected error %v, received %d", test.expectedError, err)
		})
	}
}

func TestUseCaseDeleteChild(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name          string
		mockUser      func() *mock_gateway.MockUserRepository
		mockChild     func() *mock_gateway.MockChildRepository
		input         uint64
		expectedError error
	}{
		{
			name: "error updating child",
			mockUser: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				return mockUserRepository
			},
			mockChild: func() *mock_gateway.MockChildRepository {
				mockChildRepository := mock_gateway.NewMockChildRepository(m)
				mockChildRepository.EXPECT().Get(gomock.Any()).Return(&c, nil)
				mockChildRepository.EXPECT().Delete(gomock.Any()).Return(web.ErrInternalServerError)
				return mockChildRepository
			},
			input:         c.ID,
			expectedError: web.ErrInternalServerError,
		},
		{
			name: "not found error updating child",
			mockUser: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				return mockUserRepository
			},
			mockChild: func() *mock_gateway.MockChildRepository {
				mockChildRepository := mock_gateway.NewMockChildRepository(m)
				mockChildRepository.EXPECT().Get(gomock.Any()).Return(nil, nil)
				return mockChildRepository
			},
			input:         c.ID,
			expectedError: web.ErrNotFound,
		},
		{
			name: "error getting child",
			mockUser: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				return mockUserRepository
			},
			mockChild: func() *mock_gateway.MockChildRepository {
				mockChildRepository := mock_gateway.NewMockChildRepository(m)
				mockChildRepository.EXPECT().Get(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockChildRepository
			},
			input:         c.ID,
			expectedError: web.ErrInternalServerError,
		},
		{
			name: "successful update child",
			mockUser: func() *mock_gateway.MockUserRepository {
				mockUserRepository := mock_gateway.NewMockUserRepository(m)
				return mockUserRepository
			},
			mockChild: func() *mock_gateway.MockChildRepository {
				mockChildRepository := mock_gateway.NewMockChildRepository(m)
				mockChildRepository.EXPECT().Get(gomock.Any()).Return(&c, nil)
				mockChildRepository.EXPECT().Delete(gomock.Any()).Return(nil)
				return mockChildRepository
			},
			input:         c.ID,
			expectedError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUserRepository := test.mockUser()
			mockChildRepository := test.mockChild()

			context := getContextChild(mockUserRepository, mockChildRepository)
			serviceLocator := ctx.GetServiceLocator(context)
			uc := serviceLocator.GetInstance(ChildUseCaseType).(ChildUseCase)

			err := uc.Delete(test.input, serviceLocator)

			assert.Equalf(t, test.expectedError, err, "Expected error %v, received %d", test.expectedError, err)
		})
	}
}

func getContextChild(mockUser *mock_gateway.MockUserRepository, mockChild *mock_gateway.MockChildRepository) context.Context {
	iocContext := ioc.NewContext()
	iocContext.Bind(gateway.UserRepositoryType).ToInstance(mockUser)
	iocContext.Bind(gateway.ChildRepositoryType).ToInstance(mockChild)
	iocContext.Bind(UserUseCaseType).ToInstance(NewUserUseCase())
	iocContext.Bind(ChildUseCaseType).ToInstance(NewChildUseCase())
	injector := ioc.NewInjector(iocContext)
	ctx := ctx.SetServiceLocator(context.TODO(), injector)

	return ctx
}
