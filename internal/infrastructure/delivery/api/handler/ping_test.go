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

func TestPing(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	d := mock_middleware.Dependencies{}
	router := configureRoutes(d)

	t.Run("successful ping", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "/ping", nil)
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
	router.Get("/ping", Pong)
	router.Route("/where/are/they", func(r chi.Router) {
		r.Post("/login", Login)
		r.Get("/users/{id}", Get)
		r.Post("/users/observed", CreateObservedUser)
		r.Post("/users/observer", CreateObserverUser)
	})

	return router
}
