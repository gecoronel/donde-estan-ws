package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model/web"
	mock_middleware "github.com/gecoronel/donde-estan-ws/internal/infrastructure/delivery/api/middleware/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	d := mock_middleware.Dependencies{}
	router := configureRoutes(d)

	t.Run("successful health", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "/health", nil)
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusOK, w.Code, "Expected response code %d, received %d", http.StatusOK, w.Code)
	})
}

func configureRoutes(d mock_middleware.Dependencies) *chi.Mux {
	router := chi.NewRouter()
	router.NotFound(web.DefaultNotFoundHandler)
	router.Use(mock_middleware.MockIoc(d))
	router.Get("/health", Health)
	router.Route("/where/are/they", func(r chi.Router) {
		r.Post("/users/login", Login)
		r.Get("/users/{id}", GetUser)
		r.Post("/users/observed", CreateObservedUser)
		r.Post("/users/observer", CreateObserverUser)
		r.Put("/users/observed", UpdateObservedUser)
		r.Put("/users/observer", UpdateObserverUser)
		r.Delete("/users/observed/{id}", DeleteObservedUser)
		r.Delete("/users/observer/{id}", DeleteObserverUser)
		r.Post("/users/observer/driver", AddObservedUserInObserverUser)
		r.Delete("/users/observer/driver/{id}", DeleteObservedUserInObserverUser)

		r.Get("/school-buses/{id}", GetSchoolBus)
		r.Post("/school-buses", SaveSchoolBus)
		r.Put("/school-buses", UpdateSchoolBus)
		r.Delete("/school-buses/{id}", DeleteSchoolBus)

		r.Get("/addresses/{id}", GetAddress)
		r.Post("/addresses", SaveAddress)
		r.Put("/addresses", UpdateAddress)
		r.Delete("/addresses/{id}", DeleteAddress)

		r.Get("/children/{id}", GetChild)
		r.Post("/children", SaveChild)
		r.Put("/children", UpdateChild)
		r.Delete("/children/{id}", DeleteChild)
	})

	return router
}
