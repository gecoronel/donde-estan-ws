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
	UseCase    *mock_usecase.MockUserUseCase
	Repository *mock_gateway.MockUserRepository
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
			iocContext.Bind(usecase.UserUseCaseType).ToInstance(d.UseCase)

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

// HTTP middleware setting a value on the request context
/*func MyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create new context from `r` request context, and assign key `"user"`
		// to value of `"123"`
		//ctx := context.WithValue(r.Context(), "user", "123")

		// Create context
		iocContext := ioc.NewContext()

		// Instantiates and resolve all dependencies
		//metricCollector := repository.NewMetricCollector(nrgin.Transaction(r.Context()))
		//configurationRepository := repository.NewConfigurationRepository()
		iocContext.Bind(gateway.UserRepositoryType).ToInstance(repository.NewUserRepository(db, r.Context()))

		// Register UseCase
		//iocContext.Bind(usecase.GetConfigurationsUseCaseType).ToInstance(usecase.NewGetConfigurationsUseCase())
		iocContext.Bind(usecase.UserUseCaseType).ToInstance(usecase.NewUserUseCase())
		iocContext.Bind(usecase.LoginUseCaseType).ToInstance(usecase.NewLoginUseCase())

		// call the next handler in the chain, passing the response writer and
		// the updated request object with the new context value.
		//
		// note: context.Context values are nested, so any previously set
		// values will be accessible as well, and the new `"user"` key
		// will be accessible from this point forward.

		// Set injector in context
		injector := ioc.NewInjector(iocContext)
		contx := ctx.SetServiceLocator(r.Context(), injector)

		// New request
		//newRequest := r.Clone(contx)

		//next(w, newRequest)

		next.ServeHTTP(w, r.WithContext(contx))
	})
}*/

// IoCApplication Register all your useCases, repositories and services in this middleware.
// This is the first middleware. Show webapi.Start()
/*
func IoCApplication(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		context := ioc.NewContext()

		// UseCases
		context.Bind(usecase.UserUseCaseType).ToInstance(usecase.NewUserUseCase())
		context.Bind(usecase.LoginUseCaseType).ToInstance(usecase.NewLoginUseCase())

		// Repositories
		context.Bind(gateway.UserRepositoryType).ToInstance(repository.NewUserRepository(db, utils.GetContext(c)))

		injector := ioc.NewInjector(context)

		utils.SetInjector(c, injector)
	}
}
*/
