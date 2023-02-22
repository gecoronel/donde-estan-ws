package middleware

import (
	"net/http"

	"github.com/gecoronel/donde-estan-ws/internal/bussiness/gateway"
	mock_gateway "github.com/gecoronel/donde-estan-ws/internal/bussiness/gateway/mocks"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/usecase"
	mock_usecase "github.com/gecoronel/donde-estan-ws/internal/bussiness/usecase/mocks"
	ctx "github.com/gecoronel/donde-estan-ws/internal/infrastructure/delivery/api/context"
	"github.com/gecoronel/donde-estan-ws/internal/infrastructure/delivery/api/middleware/ioc"
)

type Dependencies struct {
	UserUseCase      *mock_usecase.MockUserUseCase
	SchoolBusUseCase *mock_usecase.MockSchoolBusUseCase
	Repository       *mock_gateway.MockUserRepository
}

func MockIoc(d Dependencies) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Create context
			iocContext := ioc.NewContext()

			// Get context from request
			//ctx := r.Context()

			// Instantiates and resolve all dependencies
			//metricCollector := repository.NewMetricCollector(nrgin.Transaction(r.Context()))
			//configurationRepository := repository.NewConfigurationRepository()
			iocContext.Bind(gateway.UserRepositoryType).ToInstance(d.Repository)

			// Register UseCase
			//iocContext.Bind(usecase.GetConfigurationsUseCaseType).ToInstance(usecase.NewGetConfigurationsUseCase())
			iocContext.Bind(usecase.UserUseCaseType).ToInstance(d.UserUseCase)
			iocContext.Bind(usecase.SchoolBusUseCaseType).ToInstance(d.SchoolBusUseCase)

			// Register Repositories
			//iocContext.Bind(gateway.MetricCollectorType).ToInstance(metricCollector)

			// Register Services
			//iocContext.Bind(gateway.LocaleServiceType).ToInstance(service.NewLocaleService(r.Context(), metricCollector, configurationRepository))

			// Set logger in context
			//lvl := log.NewAtomicLevelAt(log.WarnLevel)
			//logger := log.NewProductionLogger(&lvl)
			//ctx = log.Context(ctx, logger)

			// Set injector in context
			injector := ioc.NewInjector(iocContext)
			contx := ctx.SetServiceLocator(r.Context(), injector)

			// New request
			newRequest := r.Clone(contx)

			next.ServeHTTP(w, newRequest)
		})
	}
}
