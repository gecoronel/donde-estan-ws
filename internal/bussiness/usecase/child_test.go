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
		mock          func() *mock_gateway.MockChildRepository
		input         uint64
		expectedChild *model.Child
		expectedError error
	}{
		{
			name: "error getting child",
			mock: func() *mock_gateway.MockChildRepository {
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
			mock: func() *mock_gateway.MockChildRepository {
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
			mock: func() *mock_gateway.MockChildRepository {
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
			mockChildRepository := test.mock()

			context := getContextChild(mockChildRepository)
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
		mock          func() *mock_gateway.MockChildRepository
		input         model.Child
		expectedChild *model.Child
		expectedError error
	}{
		{
			name: "error saving child",
			mock: func() *mock_gateway.MockChildRepository {
				mockChildRepository := mock_gateway.NewMockChildRepository(m)
				mockChildRepository.EXPECT().Save(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockChildRepository
			},
			input:         c,
			expectedChild: nil,
			expectedError: web.ErrInternalServerError,
		},
		{
			name: "successful save child",
			mock: func() *mock_gateway.MockChildRepository {
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
			mockChildRepository := test.mock()

			context := getContextChild(mockChildRepository)
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
		mock          func() *mock_gateway.MockChildRepository
		input         model.Child
		expectedChild *model.Child
		expectedError error
	}{
		{
			name: "error updating child",
			mock: func() *mock_gateway.MockChildRepository {
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
			mock: func() *mock_gateway.MockChildRepository {
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
			mock: func() *mock_gateway.MockChildRepository {
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
			mock: func() *mock_gateway.MockChildRepository {
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
			mockChildRepository := test.mock()

			context := getContextChild(mockChildRepository)
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
		mock          func() *mock_gateway.MockChildRepository
		input         uint64
		expectedError error
	}{
		{
			name: "error updating child",
			mock: func() *mock_gateway.MockChildRepository {
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
			mock: func() *mock_gateway.MockChildRepository {
				mockChildRepository := mock_gateway.NewMockChildRepository(m)
				mockChildRepository.EXPECT().Get(gomock.Any()).Return(nil, nil)
				return mockChildRepository
			},
			input:         c.ID,
			expectedError: web.ErrNotFound,
		},
		{
			name: "error getting child",
			mock: func() *mock_gateway.MockChildRepository {
				mockChildRepository := mock_gateway.NewMockChildRepository(m)
				mockChildRepository.EXPECT().Get(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockChildRepository
			},
			input:         c.ID,
			expectedError: web.ErrInternalServerError,
		},
		{
			name: "successful update child",
			mock: func() *mock_gateway.MockChildRepository {
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
			mockChildRepository := test.mock()

			context := getContextChild(mockChildRepository)
			serviceLocator := ctx.GetServiceLocator(context)
			uc := serviceLocator.GetInstance(ChildUseCaseType).(ChildUseCase)

			err := uc.Delete(test.input, serviceLocator)

			assert.Equalf(t, test.expectedError, err, "Expected error %v, received %d", test.expectedError, err)
		})
	}
}

func getContextChild(mock *mock_gateway.MockChildRepository) context.Context {
	iocContext := ioc.NewContext()
	iocContext.Bind(gateway.ChildRepositoryType).ToInstance(mock)
	iocContext.Bind(ChildUseCaseType).ToInstance(NewChildUseCase())
	injector := ioc.NewInjector(iocContext)
	ctx := ctx.SetServiceLocator(context.TODO(), injector)

	return ctx
}
