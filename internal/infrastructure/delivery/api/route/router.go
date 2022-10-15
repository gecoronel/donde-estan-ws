package route

import (
	"github.com/gcoron/donde-estan-ws/internal/infrastructure/delivery/api/middleware"
	"gorm.io/gorm"

	"github.com/gcoron/donde-estan-ws/internal/bussiness/model/web"
	"github.com/gcoron/donde-estan-ws/internal/infrastructure/delivery/api/handler"
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
	router.Route("/where/are/they/ws", func(r chi.Router) {
		router.Post("/users", handler.Login)
	})
	router.Post("/users", handler.Login)
}
