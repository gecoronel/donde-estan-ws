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

var sb = model.SchoolBus{
	ID:           "1",
	LicensePlate: "11AAA22",
	Model:        "Master",
	Brand:        "Renault",
	License:      "111",
	CreatedAt:    "2023-02-18 17:09:33",
	UpdatedAt:    "2023-02-18 17:09:33",
}

func TestUseCaseGetSchoolBus(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name              string
		mock              func() *mock_gateway.MockSchoolBusRepository
		input             string
		expectedSchoolBus *model.SchoolBus
		expectedError     error
	}{
		{
			name: "error getting school bus",
			mock: func() *mock_gateway.MockSchoolBusRepository {
				mockSchoolBusRepository := mock_gateway.NewMockSchoolBusRepository(m)
				mockSchoolBusRepository.EXPECT().Get(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockSchoolBusRepository
			},
			input:             sb.ID,
			expectedSchoolBus: nil,
			expectedError:     web.ErrInternalServerError,
		},
		{
			name: "not found error getting school bus",
			mock: func() *mock_gateway.MockSchoolBusRepository {
				mockSchoolBusRepository := mock_gateway.NewMockSchoolBusRepository(m)
				mockSchoolBusRepository.EXPECT().Get(gomock.Any()).Return(nil, nil)
				return mockSchoolBusRepository
			},
			input:             sb.ID,
			expectedSchoolBus: nil,
			expectedError:     web.ErrNotFound,
		},
		{
			name: "successful create observed user",
			mock: func() *mock_gateway.MockSchoolBusRepository {
				mockSchoolBusRepository := mock_gateway.NewMockSchoolBusRepository(m)
				mockSchoolBusRepository.EXPECT().Get(gomock.Any()).Return(&sb, nil)
				return mockSchoolBusRepository
			},
			input:             sb.ID,
			expectedSchoolBus: &sb,
			expectedError:     nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockSchoolBusRepository := test.mock()

			context := getContextSchoolBus(mockSchoolBusRepository)
			serviceLocator := ctx.GetServiceLocator(context)
			uc := serviceLocator.GetInstance(SchoolBusUseCaseType).(SchoolBusUseCase)

			u, err := uc.Get(test.input, serviceLocator)

			assert.Equalf(t, test.expectedSchoolBus, u, "Expected user %v, received %v", test.expectedSchoolBus, u)
			assert.Equalf(t, test.expectedError, err, "Expected error %v, received %d", test.expectedError, err)
		})
	}
}

func TestUseCaseSaveSchoolBus(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name              string
		mock              func() *mock_gateway.MockSchoolBusRepository
		input             model.SchoolBus
		expectedSchoolBus *model.SchoolBus
		expectedError     error
	}{
		{
			name: "error saving school bus",
			mock: func() *mock_gateway.MockSchoolBusRepository {
				mockSchoolBusRepository := mock_gateway.NewMockSchoolBusRepository(m)
				mockSchoolBusRepository.EXPECT().Save(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockSchoolBusRepository
			},
			input:             sb,
			expectedSchoolBus: nil,
			expectedError:     web.ErrInternalServerError,
		},
		{
			name: "successful save school bus",
			mock: func() *mock_gateway.MockSchoolBusRepository {
				mockSchoolBusRepository := mock_gateway.NewMockSchoolBusRepository(m)
				mockSchoolBusRepository.EXPECT().Save(gomock.Any()).Return(&sb, nil)
				return mockSchoolBusRepository
			},
			input:             sb,
			expectedSchoolBus: &sb,
			expectedError:     nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockSchoolBusRepository := test.mock()

			context := getContextSchoolBus(mockSchoolBusRepository)
			serviceLocator := ctx.GetServiceLocator(context)
			uc := serviceLocator.GetInstance(SchoolBusUseCaseType).(SchoolBusUseCase)

			u, err := uc.Save(test.input, serviceLocator)

			assert.Equalf(t, test.expectedSchoolBus, u, "Expected user %v, received %v", test.expectedSchoolBus, u)
			assert.Equalf(t, test.expectedError, err, "Expected error %v, received %d", test.expectedError, err)
		})
	}
}

func TestUseCaseUpdateSchoolBus(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name              string
		mock              func() *mock_gateway.MockSchoolBusRepository
		input             model.SchoolBus
		expectedSchoolBus *model.SchoolBus
		expectedError     error
	}{
		{
			name: "error updating school bus",
			mock: func() *mock_gateway.MockSchoolBusRepository {
				mockSchoolBusRepository := mock_gateway.NewMockSchoolBusRepository(m)
				mockSchoolBusRepository.EXPECT().Get(gomock.Any()).Return(&sb, nil)
				mockSchoolBusRepository.EXPECT().Update(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockSchoolBusRepository
			},
			input:             sb,
			expectedSchoolBus: nil,
			expectedError:     web.ErrInternalServerError,
		},
		{
			name: "not found error updating school bus",
			mock: func() *mock_gateway.MockSchoolBusRepository {
				mockSchoolBusRepository := mock_gateway.NewMockSchoolBusRepository(m)
				mockSchoolBusRepository.EXPECT().Get(gomock.Any()).Return(nil, nil)
				return mockSchoolBusRepository
			},
			input:             sb,
			expectedSchoolBus: nil,
			expectedError:     web.ErrNotFound,
		},
		{
			name: "error getting school bus",
			mock: func() *mock_gateway.MockSchoolBusRepository {
				mockSchoolBusRepository := mock_gateway.NewMockSchoolBusRepository(m)
				mockSchoolBusRepository.EXPECT().Get(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockSchoolBusRepository
			},
			input:             sb,
			expectedSchoolBus: nil,
			expectedError:     web.ErrInternalServerError,
		},
		{
			name: "successful update school bus",
			mock: func() *mock_gateway.MockSchoolBusRepository {
				mockSchoolBusRepository := mock_gateway.NewMockSchoolBusRepository(m)
				mockSchoolBusRepository.EXPECT().Get(gomock.Any()).Return(&sb, nil)
				mockSchoolBusRepository.EXPECT().Update(gomock.Any()).Return(&sb, nil)
				return mockSchoolBusRepository
			},
			input:             sb,
			expectedSchoolBus: &sb,
			expectedError:     nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockSchoolBusRepository := test.mock()

			context := getContextSchoolBus(mockSchoolBusRepository)
			serviceLocator := ctx.GetServiceLocator(context)
			uc := serviceLocator.GetInstance(SchoolBusUseCaseType).(SchoolBusUseCase)

			u, err := uc.Update(test.input, serviceLocator)

			assert.Equalf(t, test.expectedSchoolBus, u, "Expected user %v, received %v", test.expectedSchoolBus, u)
			assert.Equalf(t, test.expectedError, err, "Expected error %v, received %d", test.expectedError, err)
		})
	}
}

func TestUseCaseDeleteSchoolBus(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name          string
		mock          func() *mock_gateway.MockSchoolBusRepository
		input         string
		expectedError error
	}{
		{
			name: "error updating school bus",
			mock: func() *mock_gateway.MockSchoolBusRepository {
				mockSchoolBusRepository := mock_gateway.NewMockSchoolBusRepository(m)
				mockSchoolBusRepository.EXPECT().Get(gomock.Any()).Return(&sb, nil)
				mockSchoolBusRepository.EXPECT().Delete(gomock.Any()).Return(web.ErrInternalServerError)
				return mockSchoolBusRepository
			},
			input:         sb.ID,
			expectedError: web.ErrInternalServerError,
		},
		{
			name: "not found error updating school bus",
			mock: func() *mock_gateway.MockSchoolBusRepository {
				mockSchoolBusRepository := mock_gateway.NewMockSchoolBusRepository(m)
				mockSchoolBusRepository.EXPECT().Get(gomock.Any()).Return(nil, nil)
				return mockSchoolBusRepository
			},
			input:         sb.ID,
			expectedError: web.ErrNotFound,
		},
		{
			name: "error getting school bus",
			mock: func() *mock_gateway.MockSchoolBusRepository {
				mockSchoolBusRepository := mock_gateway.NewMockSchoolBusRepository(m)
				mockSchoolBusRepository.EXPECT().Get(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockSchoolBusRepository
			},
			input:         sb.ID,
			expectedError: web.ErrInternalServerError,
		},
		{
			name: "successful update school bus",
			mock: func() *mock_gateway.MockSchoolBusRepository {
				mockSchoolBusRepository := mock_gateway.NewMockSchoolBusRepository(m)
				mockSchoolBusRepository.EXPECT().Get(gomock.Any()).Return(&sb, nil)
				mockSchoolBusRepository.EXPECT().Delete(gomock.Any()).Return(nil)
				return mockSchoolBusRepository
			},
			input:         sb.ID,
			expectedError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockSchoolBusRepository := test.mock()

			context := getContextSchoolBus(mockSchoolBusRepository)
			serviceLocator := ctx.GetServiceLocator(context)
			uc := serviceLocator.GetInstance(SchoolBusUseCaseType).(SchoolBusUseCase)

			err := uc.Delete(test.input, serviceLocator)

			assert.Equalf(t, test.expectedError, err, "Expected error %v, received %d", test.expectedError, err)
		})
	}
}

func getContextSchoolBus(mock *mock_gateway.MockSchoolBusRepository) context.Context {
	iocContext := ioc.NewContext()
	iocContext.Bind(gateway.SchoolBusRepositoryType).ToInstance(mock)
	iocContext.Bind(SchoolBusUseCaseType).ToInstance(NewSchoolBusUseCase())
	injector := ioc.NewInjector(iocContext)
	ctx := ctx.SetServiceLocator(context.TODO(), injector)

	return ctx
}
