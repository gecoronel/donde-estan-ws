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
		r.Put("/users/observed", handler.UpdateObservedUser)
		r.Put("/users/observer", handler.UpdateObserverUser)
		r.Delete("/users/observed/{id}", handler.DeleteObservedUser)
		r.Delete("/users/observer/{id}", handler.DeleteObserverUser)
		r.Post("/users/observer/driver", handler.AddObservedUserInObserverUser)
		r.Delete("/users/observer/driver/{id}", handler.DeleteObservedUserInObserverUser)

		r.Get("/school-buses/{id}", handler.GetSchoolBus)
		r.Post("/school-buses", handler.SaveSchoolBus)
		r.Put("/school-buses", handler.UpdateSchoolBus)
		r.Delete("/school-buses/{id}", handler.DeleteSchoolBus)

		r.Get("/addresses/{id}", handler.GetAddress)
		r.Post("/addresses", handler.SaveAddress)
		r.Put("/addresses", handler.UpdateAddress)
		r.Delete("/addresses/{id}", handler.DeleteAddress)

		r.Get("/children/{id}", handler.GetChild)
		r.Post("/children", handler.SaveChild)
		r.Put("/children", handler.UpdateChild)
		r.Delete("/children/{id}", handler.DeleteChild)
	})
}
