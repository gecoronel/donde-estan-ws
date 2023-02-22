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

	t.Run("error getting school bus", func(t *testing.T) {
		mockSchoolBusRepository := mock_gateway.NewMockSchoolBusRepository(m)
		mockSchoolBusRepository.EXPECT().Get(gomock.Any()).Return(nil, web.ErrInternalServerError)

		context := getContextSchoolBus(mockSchoolBusRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(SchoolBusUseCaseType).(SchoolBusUseCase)

		u, err := uc.Get("1", serviceLocator)

		assert.Nil(t, u)
		assert.Equalf(t, web.ErrInternalServerError, err, "Expected error %v, received %d", nil, err)
	})

	t.Run("not found error getting school bus", func(t *testing.T) {
		mockSchoolBusRepository := mock_gateway.NewMockSchoolBusRepository(m)
		mockSchoolBusRepository.EXPECT().Get(gomock.Any()).Return(nil, nil)

		context := getContextSchoolBus(mockSchoolBusRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(SchoolBusUseCaseType).(SchoolBusUseCase)

		u, err := uc.Get("1", serviceLocator)

		assert.Nil(t, u)
		assert.Equalf(t, web.ErrNotFound, err, "Expected error %v, received %d", nil, err)
	})

	t.Run("successful get school bus", func(t *testing.T) {
		mockSchoolBusRepository := mock_gateway.NewMockSchoolBusRepository(m)
		mockSchoolBusRepository.EXPECT().Get(gomock.Any()).Return(&sb, nil)

		context := getContextSchoolBus(mockSchoolBusRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(SchoolBusUseCaseType).(SchoolBusUseCase)

		u, err := uc.Get("1", serviceLocator)

		assert.Equalf(t, sb, *u, "Expected response %v, received %v", sb, u)
		assert.Equalf(t, nil, err, "Expected error %v, received %d", nil, err)
	})
}

func TestUseCaseSaveSchoolBus(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	t.Run("error saving school bus", func(t *testing.T) {
		mockSchoolBusRepository := mock_gateway.NewMockSchoolBusRepository(m)
		mockSchoolBusRepository.EXPECT().Save(gomock.Any()).Return(nil, web.ErrInternalServerError)

		context := getContextSchoolBus(mockSchoolBusRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(SchoolBusUseCaseType).(SchoolBusUseCase)

		bus, err := uc.Save(sb, serviceLocator)

		assert.Nil(t, bus)
		assert.Equalf(t, web.ErrInternalServerError, err, "Expected error %v, received %d", nil, err)
	})

	t.Run("successful save school bus", func(t *testing.T) {
		mockSchoolBusRepository := mock_gateway.NewMockSchoolBusRepository(m)
		mockSchoolBusRepository.EXPECT().Save(gomock.Any()).Return(&sb, nil)

		context := getContextSchoolBus(mockSchoolBusRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(SchoolBusUseCaseType).(SchoolBusUseCase)

		bus, err := uc.Save(sb, serviceLocator)

		assert.Equalf(t, sb, *bus, "Expected response %v, received %v", sb, bus)
		assert.Equalf(t, nil, err, "Expected error %v, received %d", nil, err)
	})
}

func TestUseCaseUpdateSchoolBus(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	t.Run("error updating school bus", func(t *testing.T) {
		mockSchoolBusRepository := mock_gateway.NewMockSchoolBusRepository(m)
		mockSchoolBusRepository.EXPECT().Get(gomock.Any()).Return(&sb, nil)
		mockSchoolBusRepository.EXPECT().Update(gomock.Any()).Return(nil, web.ErrInternalServerError)

		context := getContextSchoolBus(mockSchoolBusRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(SchoolBusUseCaseType).(SchoolBusUseCase)

		bus, err := uc.Update(sb, serviceLocator)

		assert.Nil(t, bus)
		assert.Equalf(t, web.ErrInternalServerError, err, "Expected error %v, received %d", nil, err)
	})

	t.Run("error selecting school bus", func(t *testing.T) {
		mockSchoolBusRepository := mock_gateway.NewMockSchoolBusRepository(m)
		mockSchoolBusRepository.EXPECT().Get(gomock.Any()).Return(nil, web.ErrInternalServerError)

		context := getContextSchoolBus(mockSchoolBusRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(SchoolBusUseCaseType).(SchoolBusUseCase)

		bus, err := uc.Update(sb, serviceLocator)

		assert.Nil(t, bus)
		assert.Equalf(t, web.ErrInternalServerError, err, "Expected error %v, received %d", nil, err)
	})

	t.Run("not found error updating school bus", func(t *testing.T) {
		mockSchoolBusRepository := mock_gateway.NewMockSchoolBusRepository(m)
		mockSchoolBusRepository.EXPECT().Get(gomock.Any()).Return(nil, nil)

		context := getContextSchoolBus(mockSchoolBusRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(SchoolBusUseCaseType).(SchoolBusUseCase)

		bus, err := uc.Update(sb, serviceLocator)

		assert.Nil(t, bus)
		assert.Equalf(t, web.ErrNotFound, err, "Expected error %v, received %d", nil, err)
	})

	t.Run("successful update school bus", func(t *testing.T) {
		sb.Brand = "Fiat"
		mockSchoolBusRepository := mock_gateway.NewMockSchoolBusRepository(m)
		mockSchoolBusRepository.EXPECT().Get(gomock.Any()).Return(&sb, nil)
		mockSchoolBusRepository.EXPECT().Update(gomock.Any()).Return(&sb, nil)

		context := getContextSchoolBus(mockSchoolBusRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(SchoolBusUseCaseType).(SchoolBusUseCase)

		bus, err := uc.Update(sb, serviceLocator)

		assert.Equalf(t, "Fiat", bus.Brand, "Expected response %v, received %v", "Fiat", bus.Brand)
		assert.Equalf(t, sb, *bus, "Expected response %v, received %v", sb, bus)
		assert.Equalf(t, nil, err, "Expected error %v, received %d", nil, err)
	})
}

func TestUseCaseDeleteSchoolBus(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	t.Run("error updating school bus", func(t *testing.T) {
		mockSchoolBusRepository := mock_gateway.NewMockSchoolBusRepository(m)
		mockSchoolBusRepository.EXPECT().Get(gomock.Any()).Return(&sb, nil)
		mockSchoolBusRepository.EXPECT().Delete(gomock.Any()).Return(web.ErrInternalServerError)

		context := getContextSchoolBus(mockSchoolBusRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(SchoolBusUseCaseType).(SchoolBusUseCase)

		err := uc.Delete(sb.ID, serviceLocator)

		assert.Equalf(t, web.ErrInternalServerError, err, "Expected error %v, received %d", nil, err)
	})

	t.Run("error selecting school bus", func(t *testing.T) {
		mockSchoolBusRepository := mock_gateway.NewMockSchoolBusRepository(m)
		mockSchoolBusRepository.EXPECT().Get(gomock.Any()).Return(nil, web.ErrInternalServerError)

		context := getContextSchoolBus(mockSchoolBusRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(SchoolBusUseCaseType).(SchoolBusUseCase)

		err := uc.Delete(sb.ID, serviceLocator)

		assert.Equalf(t, web.ErrInternalServerError, err, "Expected error %v, received %d", nil, err)
	})

	t.Run("not found error updating school bus", func(t *testing.T) {
		mockSchoolBusRepository := mock_gateway.NewMockSchoolBusRepository(m)
		mockSchoolBusRepository.EXPECT().Get(gomock.Any()).Return(nil, nil)

		context := getContextSchoolBus(mockSchoolBusRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(SchoolBusUseCaseType).(SchoolBusUseCase)

		err := uc.Delete(sb.ID, serviceLocator)

		assert.Equalf(t, web.ErrNotFound, err, "Expected error %v, received %d", nil, err)
	})

	t.Run("successful update school bus", func(t *testing.T) {
		mockSchoolBusRepository := mock_gateway.NewMockSchoolBusRepository(m)
		mockSchoolBusRepository.EXPECT().Get(gomock.Any()).Return(&sb, nil)
		mockSchoolBusRepository.EXPECT().Delete(gomock.Any()).Return(nil)

		context := getContextSchoolBus(mockSchoolBusRepository)
		serviceLocator := ctx.GetServiceLocator(context)
		uc := serviceLocator.GetInstance(SchoolBusUseCaseType).(SchoolBusUseCase)

		err := uc.Delete(sb.ID, serviceLocator)

		assert.Equalf(t, nil, err, "Expected error %v, received %d", nil, err)
	})
}

func getContextSchoolBus(mock *mock_gateway.MockSchoolBusRepository) context.Context {
	iocContext := ioc.NewContext()
	iocContext.Bind(gateway.SchoolBusRepositoryType).ToInstance(mock)
	iocContext.Bind(SchoolBusUseCaseType).ToInstance(NewSchoolBusUseCase())
	injector := ioc.NewInjector(iocContext)
	ctx := ctx.SetServiceLocator(context.TODO(), injector)

	return ctx
}
