package route

import (
	"gorm.io/gorm"

	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model/web"
	"github.com/gecoronel/donde-estan-ws/internal/infrastructure/delivery/api/handler"
	"github.com/gecoronel/donde-estan-ws/internal/infrastructure/delivery/api/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(db *gorm.DB) *chi.Mux { //*gin.Engine {
	/*
		router := gin.Default()
		healthyCheckGroup := router.Group("/ping")
		applicationGroup := router.Group("/where/are/they/ws")

		applicationGroup.Use(middleware.IoCApplication(db))

		healthyCheckGroup.GET("", handler.PongHandler)
		applicationGroup.POST("/login", handler.LoginHandler)
	*/

	r := chi.NewRouter()
	r.NotFound(web.DefaultNotFoundHandler)

	//r.Use(mid.RequestID)
	//r.Use(mid.RealIP)
	//r.Use(mid.Logger)
	//r.Use(mid.Recoverer)
	r.Use(middleware.Ioc(db))

	configureRoutes(r)

	return r
}

func configureRoutes(router *chi.Mux) {
	router.Get("/health", handler.Health)
	router.Route("/where/are/they", func(r chi.Router) {
		r.Get("/users/{id}", handler.GetUser)
		r.Post("/users/login", handler.Login)
		r.Post("/users/observed", handler.CreateObservedUser)
		r.Post("/users/observer", handler.CreateObserverUser)
		r.Post("/users/observer/driver", handler.AddObservedUserInObserverUser)
		r.Delete("/users/observer/driver/{id}", handler.DeleteObservedUserInObserverUser)

		r.Get("/school-bus/{id}", handler.GetSchoolBus)
		r.Post("/school-bus", handler.SaveSchoolBus)
		r.Put("/school-bus", handler.UpdateSchoolBus)
		r.Delete("/school-bus/{id}", handler.DeleteSchoolBus)

		r.Get("/address/{id}", handler.GetAddress)
		r.Post("/address", handler.SaveAddress)
		r.Put("/address", handler.UpdateAddress)
		r.Delete("/address/{id}", handler.DeleteAddress)
	})
}
