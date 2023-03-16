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

var a = model.Address{
	ID:        1,
	Name:      "Casa",
	Street:    "25 de Mayo",
	Number:    "1010",
	Floor:     "1",
	Apartment: "A",
	ZipCode:   "3000",
	City:      "Santa Fe",
	State:     "Santa Fe",
	Country:   "Argentina",
	Latitude:  "60.0000121",
	Longitude: "-19.23423",
	CreatedAt: "2023-02-18 17:09:33",
	UpdatedAt: "2023-02-18 17:09:33",
}

func TestUseCaseGetAddress(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name            string
		mock            func() *mock_gateway.MockAddressRepository
		input           uint64
		expectedAddress *model.Address
		expectedError   error
	}{
		{
			name: "error getting school bus",
			mock: func() *mock_gateway.MockAddressRepository {
				mockAddressRepository := mock_gateway.NewMockAddressRepository(m)
				mockAddressRepository.EXPECT().Get(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockAddressRepository
			},
			input:           a.ID,
			expectedAddress: nil,
			expectedError:   web.ErrInternalServerError,
		},
		{
			name: "not found error getting school bus",
			mock: func() *mock_gateway.MockAddressRepository {
				mockAddressRepository := mock_gateway.NewMockAddressRepository(m)
				mockAddressRepository.EXPECT().Get(gomock.Any()).Return(nil, nil)
				return mockAddressRepository
			},
			input:           a.ID,
			expectedAddress: nil,
			expectedError:   web.ErrNotFound,
		},
		{
			name: "successful create observed user",
			mock: func() *mock_gateway.MockAddressRepository {
				mockAddressRepository := mock_gateway.NewMockAddressRepository(m)
				mockAddressRepository.EXPECT().Get(gomock.Any()).Return(&a, nil)
				return mockAddressRepository
			},
			input:           a.ID,
			expectedAddress: &a,
			expectedError:   nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockAddressRepository := test.mock()

			context := getContextAddress(mockAddressRepository)
			serviceLocator := ctx.GetServiceLocator(context)
			uc := serviceLocator.GetInstance(AddressUseCaseType).(AddressUseCase)

			u, err := uc.Get(test.input, serviceLocator)

			assert.Equalf(t, test.expectedAddress, u, "Expected user %v, received %v", test.expectedAddress, u)
			assert.Equalf(t, test.expectedError, err, "Expected error %v, received %d", test.expectedError, err)
		})
	}
}

func TestUseCaseSaveAddress(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name            string
		mock            func() *mock_gateway.MockAddressRepository
		input           model.Address
		expectedAddress *model.Address
		expectedError   error
	}{
		{
			name: "error saving school bus",
			mock: func() *mock_gateway.MockAddressRepository {
				mockAddressRepository := mock_gateway.NewMockAddressRepository(m)
				mockAddressRepository.EXPECT().Save(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockAddressRepository
			},
			input:           a,
			expectedAddress: nil,
			expectedError:   web.ErrInternalServerError,
		},
		{
			name: "successful save school bus",
			mock: func() *mock_gateway.MockAddressRepository {
				mockAddressRepository := mock_gateway.NewMockAddressRepository(m)
				mockAddressRepository.EXPECT().Save(gomock.Any()).Return(&a, nil)
				return mockAddressRepository
			},
			input:           a,
			expectedAddress: &a,
			expectedError:   nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockAddressRepository := test.mock()

			context := getContextAddress(mockAddressRepository)
			serviceLocator := ctx.GetServiceLocator(context)
			uc := serviceLocator.GetInstance(AddressUseCaseType).(AddressUseCase)

			u, err := uc.Save(test.input, serviceLocator)

			assert.Equalf(t, test.expectedAddress, u, "Expected user %v, received %v", test.expectedAddress, u)
			assert.Equalf(t, test.expectedError, err, "Expected error %v, received %d", test.expectedError, err)
		})
	}
}

func TestUseCaseUpdateAddress(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name            string
		mock            func() *mock_gateway.MockAddressRepository
		input           model.Address
		expectedAddress *model.Address
		expectedError   error
	}{
		{
			name: "error updating school bus",
			mock: func() *mock_gateway.MockAddressRepository {
				mockAddressRepository := mock_gateway.NewMockAddressRepository(m)
				mockAddressRepository.EXPECT().Get(gomock.Any()).Return(&a, nil)
				mockAddressRepository.EXPECT().Update(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockAddressRepository
			},
			input:           a,
			expectedAddress: nil,
			expectedError:   web.ErrInternalServerError,
		},
		{
			name: "not found error updating school bus",
			mock: func() *mock_gateway.MockAddressRepository {
				mockAddressRepository := mock_gateway.NewMockAddressRepository(m)
				mockAddressRepository.EXPECT().Get(gomock.Any()).Return(nil, nil)
				return mockAddressRepository
			},
			input:           a,
			expectedAddress: nil,
			expectedError:   web.ErrNotFound,
		},
		{
			name: "error getting school bus",
			mock: func() *mock_gateway.MockAddressRepository {
				mockAddressRepository := mock_gateway.NewMockAddressRepository(m)
				mockAddressRepository.EXPECT().Get(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockAddressRepository
			},
			input:           a,
			expectedAddress: nil,
			expectedError:   web.ErrInternalServerError,
		},
		{
			name: "successful update school bus",
			mock: func() *mock_gateway.MockAddressRepository {
				mockAddressRepository := mock_gateway.NewMockAddressRepository(m)
				mockAddressRepository.EXPECT().Get(gomock.Any()).Return(&a, nil)
				mockAddressRepository.EXPECT().Update(gomock.Any()).Return(&a, nil)
				return mockAddressRepository
			},
			input:           a,
			expectedAddress: &a,
			expectedError:   nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockAddressRepository := test.mock()

			context := getContextAddress(mockAddressRepository)
			serviceLocator := ctx.GetServiceLocator(context)
			uc := serviceLocator.GetInstance(AddressUseCaseType).(AddressUseCase)

			u, err := uc.Update(test.input, serviceLocator)

			assert.Equalf(t, test.expectedAddress, u, "Expected user %v, received %v", test.expectedAddress, u)
			assert.Equalf(t, test.expectedError, err, "Expected error %v, received %d", test.expectedError, err)
		})
	}
}

func TestUseCaseDeleteAddress(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name          string
		mock          func() *mock_gateway.MockAddressRepository
		input         uint64
		expectedError error
	}{
		{
			name: "error updating school bus",
			mock: func() *mock_gateway.MockAddressRepository {
				mockAddressRepository := mock_gateway.NewMockAddressRepository(m)
				mockAddressRepository.EXPECT().Get(gomock.Any()).Return(&a, nil)
				mockAddressRepository.EXPECT().Delete(gomock.Any()).Return(web.ErrInternalServerError)
				return mockAddressRepository
			},
			input:         a.ID,
			expectedError: web.ErrInternalServerError,
		},
		{
			name: "not found error updating school bus",
			mock: func() *mock_gateway.MockAddressRepository {
				mockAddressRepository := mock_gateway.NewMockAddressRepository(m)
				mockAddressRepository.EXPECT().Get(gomock.Any()).Return(nil, nil)
				return mockAddressRepository
			},
			input:         a.ID,
			expectedError: web.ErrNotFound,
		},
		{
			name: "error getting school bus",
			mock: func() *mock_gateway.MockAddressRepository {
				mockAddressRepository := mock_gateway.NewMockAddressRepository(m)
				mockAddressRepository.EXPECT().Get(gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockAddressRepository
			},
			input:         a.ID,
			expectedError: web.ErrInternalServerError,
		},
		{
			name: "successful update school bus",
			mock: func() *mock_gateway.MockAddressRepository {
				mockAddressRepository := mock_gateway.NewMockAddressRepository(m)
				mockAddressRepository.EXPECT().Get(gomock.Any()).Return(&a, nil)
				mockAddressRepository.EXPECT().Delete(gomock.Any()).Return(nil)
				return mockAddressRepository
			},
			input:         a.ID,
			expectedError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockAddressRepository := test.mock()

			context := getContextAddress(mockAddressRepository)
			serviceLocator := ctx.GetServiceLocator(context)
			uc := serviceLocator.GetInstance(AddressUseCaseType).(AddressUseCase)

			err := uc.Delete(test.input, serviceLocator)

			assert.Equalf(t, test.expectedError, err, "Expected error %v, received %d", test.expectedError, err)
		})
	}
}

func getContextAddress(mock *mock_gateway.MockAddressRepository) context.Context {
	iocContext := ioc.NewContext()
	iocContext.Bind(gateway.AddressRepositoryType).ToInstance(mock)
	iocContext.Bind(AddressUseCaseType).ToInstance(NewAddressUseCase())
	injector := ioc.NewInjector(iocContext)
	ctx := ctx.SetServiceLocator(context.TODO(), injector)

	return ctx
}
