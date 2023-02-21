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
	router.Get("/ping", handler.Pong)
	router.Route("/where/are/they", func(r chi.Router) {
		r.Get("/users/{id}", handler.Get)
		r.Post("/login", handler.Login)
		r.Post("/users/observed", handler.CreateObservedUser)
		r.Post("/users/observer", handler.CreateObserverUser)
	})
}
