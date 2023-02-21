package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model/web"
	mock_usecase "github.com/gecoronel/donde-estan-ws/internal/bussiness/usecase/mocks"
	mock_middleware "github.com/gecoronel/donde-estan-ws/internal/infrastructure/delivery/api/middleware/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	user := model.User{
		ID:        1,
		Name:      "Juan",
		LastName:  "Perez",
		IDNumber:  "12345678",
		Username:  "jperez",
		Password:  "jperez1234",
		Email:     "jperez@mail.com",
		Enabled:   true,
		Type:      "observed",
		CreatedAt: "2022-12-10 17:49:30",
		UpdatedAt: "2022-12-10 17:49:30",
	}

	m := gomock.NewController(t)
	defer m.Finish()

	d := mock_middleware.Dependencies{}
	router := configureRoutes(d)

	t.Run("error get user", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "/where/are/they/users/:id", nil)
		if err != nil {
			t.Fatal(err)
		}
		r.URL.Query().Set("id", "invalid")
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusBadRequest, w.Code, "Expected response code %d, received %d", http.StatusOK, w.Code)
	})

	t.Run("bad request error get user", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "/where/are/they/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
		mockUserUseCase.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, web.ErrInternalServerError)
		d = mock_middleware.Dependencies{UseCase: mockUserUseCase}
		router = configureRoutes(d)
		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusInternalServerError, w.Code, "Expected response code %d, received %d", http.StatusOK, w.Code)
	})

	t.Run("successful get user", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "/where/are/they/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
		mockUserUseCase.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&user, nil)
		d = mock_middleware.Dependencies{UseCase: mockUserUseCase}
		router = configureRoutes(d)
		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusOK, w.Code, "Expected response code %d, received %d", http.StatusOK, w.Code)
	})
}

func TestLogin(t *testing.T) {
	user := model.ObservedUser{
		User: model.User{
			ID:        1,
			Name:      "Juan",
			LastName:  "Perez",
			IDNumber:  "12345678",
			Username:  "jperez",
			Password:  "jperez1234",
			Email:     "jperez@mail.com",
			Enabled:   true,
			Type:      "observed",
			CreatedAt: "2022-12-10 17:49:30",
			UpdatedAt: "2022-12-10 17:49:30",
		},
		PrivacyKey:  "juan.perez.12345678",
		CompanyName: "school bus",
		SchoolBus: model.SchoolBus{
			ID:           "1",
			LicensePlate: "11AAA222",
			Model:        "Master",
			Brand:        "Renault",
			License:      "11222",
			CreatedAt:    "2022-12-10 17:49:30",
			UpdatedAt:    "2022-12-10 17:49:30",
		},
		ObserverUsers: nil,
	}
	u := model.NewObservedUser(user)

	m := gomock.NewController(t)
	defer m.Finish()

	d := mock_middleware.Dependencies{}
	router := configureRoutes(d)

	t.Run("bad request for handler login", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest("POST", "/where/are/they/login", bytes.NewBuffer([]byte(`invalid"`)))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusBadRequest, w.Code, "Expected response code %d, received %d", http.StatusBadRequest, w.Code)
	})

	t.Run("bad request in login validation", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest("POST", "/where/are/they/login", bytes.NewBuffer([]byte(`{"username": "jperez"}`)))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusBadRequest, w.Code, "Expected response code %d, received %d", http.StatusBadRequest, w.Code)
	})

	t.Run("successful login", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest("POST", "/where/are/they/login", bytes.NewBuffer([]byte(`{"username": "jperez", "password": "jperez1234"}`)))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
		mockUserUseCase.EXPECT().Login(gomock.Any(), gomock.Any()).Return(u, nil)
		d = mock_middleware.Dependencies{UseCase: mockUserUseCase}
		router = configureRoutes(d)
		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusOK, w.Code, "Expected response code %d, received %d", http.StatusOK, w.Code)
	})

	t.Run("unsuccessful login", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest("POST", "/where/are/they/login", bytes.NewBuffer([]byte(`{"username": "jperez", "password": "jperez1234"}`)))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
		mockUserUseCase.EXPECT().Login(gomock.Any(), gomock.Any()).Return(nil, web.ErrInternalServerError)
		d = mock_middleware.Dependencies{UseCase: mockUserUseCase}
		router = configureRoutes(d)
		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusInternalServerError, w.Code, "Expected response code %d, received %d", http.StatusOK, w.Code)
	})
}

func TestCreateObservedUser(t *testing.T) {
	user := model.ObservedUser{
		User: model.User{
			ID:        1,
			Name:      "Juan",
			LastName:  "Perez",
			IDNumber:  "12345678",
			Username:  "jperez",
			Password:  "jperez1234",
			Email:     "jperez@mail.com",
			Enabled:   true,
			Type:      "observed",
			CreatedAt: "2022-12-10 17:49:30",
			UpdatedAt: "2022-12-10 17:49:30",
		},
		PrivacyKey:  "juan.perez.12345678",
		CompanyName: "school bus",
		SchoolBus: model.SchoolBus{
			ID:           "1",
			LicensePlate: "11AAA222",
			Model:        "Master",
			Brand:        "Renault",
			License:      "11222",
			CreatedAt:    "2022-12-10 17:49:30",
			UpdatedAt:    "2022-12-10 17:49:30",
		},
		ObserverUsers: nil,
	}

	m := gomock.NewController(t)
	defer m.Finish()

	d := mock_middleware.Dependencies{}
	router := configureRoutes(d)

	t.Run("bad request for handler create observed user", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest("POST", "/where/are/they/users/observed", bytes.NewBuffer([]byte(
			`invalid"`)))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusBadRequest, w.Code, "Expected response code %d, received %d",
			http.StatusBadRequest, w.Code)
	})

	t.Run("bad request in login validation", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest("POST", "/where/are/they/users/observed", bytes.NewBuffer([]byte(
			`{"user": {"username": "jperez"}}`)))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusBadRequest, w.Code, "Expected response code %d, received %d",
			http.StatusBadRequest, w.Code)
	})

	t.Run("username conflict error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest("POST", "/where/are/they/users/observed", bytes.NewBuffer([]byte(
			`{
				"user": {
					"name": "Juan",
					"last_name": "Perez",
					"id_number": "12345678",
					"username": "jperez",
					"password": "jperez1234",
					"email": "jperez@mail.com",
					"enabled": true,
					"type": "observed"
				},
				"privacy_key": "juan.perez.1234",
				"company_name": "company school bus",
				"school_bus": {
					"id": "0000-0000-0004",
					"license_plate": "11AAA22",
					"model": "Master",
					"brand": "Renault",
					"license": "111"
				}
			}`,
		)))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
		mockUserUseCase.EXPECT().FindByUsername(gomock.Any(), gomock.Any()).Return(&user.User, nil)
		d = mock_middleware.Dependencies{UseCase: mockUserUseCase}
		router = configureRoutes(d)
		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusConflict, w.Code, "Expected response code %d, received %d", http.StatusOK,
			w.Code)
	})

	t.Run("email conflict error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest("POST", "/where/are/they/users/observed", bytes.NewBuffer([]byte(
			`{
				"user": {
					"name": "Juan",
					"last_name": "Perez",
					"id_number": "12345678",
					"username": "jperez",
					"password": "jperez1234",
					"email": "jperez@mail.com",
					"enabled": true,
					"type": "observed"
				},
				"privacy_key": "juan.perez.1234",
				"company_name": "company school bus",
				"school_bus": {
					"id": "0000-0000-0004",
					"license_plate": "11AAA22",
					"model": "Master",
					"brand": "Renault",
					"license": "111"
				}
			}`,
		)))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
		mockUserUseCase.EXPECT().FindByUsername(gomock.Any(), gomock.Any()).Return(nil, nil)
		mockUserUseCase.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(&user.User, nil)
		d = mock_middleware.Dependencies{UseCase: mockUserUseCase}
		router = configureRoutes(d)
		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusConflict, w.Code, "Expected response code %d, received %d", http.StatusOK,
			w.Code)
	})

	t.Run("successful creation", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest("POST", "/where/are/they/users/observed", bytes.NewBuffer([]byte(
			`{
				"user": {
					"name": "Juan",
					"last_name": "Perez",
					"id_number": "12345678",
					"username": "jperez",
					"password": "jperez1234",
					"email": "jperez@mail.com",
					"enabled": true,
					"type": "observed"
				},
				"privacy_key": "juan.perez.1234",
				"company_name": "company school bus",
				"school_bus": {
					"id": "0000-0000-0004",
					"license_plate": "11AAA22",
					"model": "Master",
					"brand": "Renault",
					"license": "111"
				}
			}`,
		)))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
		mockUserUseCase.EXPECT().FindByUsername(gomock.Any(), gomock.Any()).Return(nil, nil)
		mockUserUseCase.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(nil, nil)
		mockUserUseCase.EXPECT().CreateObservedUser(gomock.Any(), gomock.Any()).Return(&user, nil)
		d = mock_middleware.Dependencies{UseCase: mockUserUseCase}
		router = configureRoutes(d)
		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusOK, w.Code, "Expected response code %d, received %d", http.StatusOK, w.Code)
	})

	t.Run("unsuccessful creation", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest("POST", "/where/are/they/users/observed", bytes.NewBuffer([]byte(
			`{
				"user": {
					"name": "Juan",
					"last_name": "Perez",
					"id_number": "12345678",
					"username": "jperez",
					"password": "jperez1234",
					"email": "jperez@mail.com",
					"enabled": true,
					"type": "observed"
				},
				"privacy_key": "juan.perez.1234",
				"company_name": "company school bus",
				"school_bus": {
					"id": "0000-0000-0004",
					"license_plate": "11AAA22",
					"model": "Master",
					"brand": "Renault",
					"license": "111"
				}
			}`,
		)))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
		mockUserUseCase.EXPECT().FindByUsername(gomock.Any(), gomock.Any()).Return(nil, nil)
		mockUserUseCase.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(nil, nil)
		mockUserUseCase.EXPECT().CreateObservedUser(gomock.Any(), gomock.Any()).Return(nil,
			web.ErrInternalServerError)
		d = mock_middleware.Dependencies{UseCase: mockUserUseCase}
		router = configureRoutes(d)
		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusInternalServerError, w.Code, "Expected response code %d, received %d",
			http.StatusOK, w.Code)
	})
}

func TestCreateObserverUser(t *testing.T) {
	user := model.ObserverUser{
		User: model.User{
			ID:        2,
			Name:      "Maria",
			LastName:  "Dominguez",
			IDNumber:  "12345678",
			Username:  "mdominguez",
			Password:  "mdominguez1234",
			Email:     "mdominguez@mail.com",
			Enabled:   true,
			Type:      "observer",
			CreatedAt: "2022-12-10 17:49:30",
			UpdatedAt: "2022-12-10 17:49:30",
		},
	}

	m := gomock.NewController(t)
	defer m.Finish()

	d := mock_middleware.Dependencies{}
	router := configureRoutes(d)

	t.Run("bad request for handler create observer user", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest("POST", "/where/are/they/users/observer", bytes.NewBuffer([]byte(
			`invalid"`)))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusBadRequest, w.Code, "Expected response code %d, received %d",
			http.StatusBadRequest, w.Code)
	})

	t.Run("bad request in validation", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest("POST", "/where/are/they/users/observer", bytes.NewBuffer([]byte(
			`{"user": {"username": "mdominguez"}}`)))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusBadRequest, w.Code, "Expected response code %d, received %d",
			http.StatusBadRequest, w.Code)
	})

	t.Run("username conflict error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest("POST", "/where/are/they/users/observer", bytes.NewBuffer([]byte(
			`{
				"user": {
					"name": "Maria",
					"last_name": "Dominguez",
					"id_number": "12345678",
					"username": "mdominguez",
					"password": "mdominguez1234",
					"email": "mdominguez@mail.com",
					"enabled": true,
					"type": "observer"
				}
			}`,
		)))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
		mockUserUseCase.EXPECT().FindByUsername(gomock.Any(), gomock.Any()).Return(&user.User, nil)
		d = mock_middleware.Dependencies{UseCase: mockUserUseCase}
		router = configureRoutes(d)
		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusConflict, w.Code, "Expected response code %d, received %d", http.StatusOK,
			w.Code)
	})

	t.Run("email conflict error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest("POST", "/where/are/they/users/observer", bytes.NewBuffer([]byte(
			`{
				"user": {
					"name": "Maria",
					"last_name": "Dominguez",
					"id_number": "12345678",
					"username": "mdominguez",
					"password": "mdominguez1234",
					"email": "mdominguez@mail.com",
					"enabled": true,
					"type": "observer"
				}
			}`,
		)))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
		mockUserUseCase.EXPECT().FindByUsername(gomock.Any(), gomock.Any()).Return(nil, nil)
		mockUserUseCase.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(&user.User, nil)
		d = mock_middleware.Dependencies{UseCase: mockUserUseCase}
		router = configureRoutes(d)
		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusConflict, w.Code, "Expected response code %d, received %d", http.StatusOK,
			w.Code)
	})

	t.Run("successful creation", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest("POST", "/where/are/they/users/observer", bytes.NewBuffer([]byte(
			`{
				"user": {
					"name": "Maria",
					"last_name": "Dominguez",
					"id_number": "12345678",
					"username": "mdominguez",
					"password": "mdominguez1234",
					"email": "mdominguez@mail.com",
					"enabled": true,
					"type": "observer"
				}
			}`,
		)))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
		mockUserUseCase.EXPECT().FindByUsername(gomock.Any(), gomock.Any()).Return(nil, nil)
		mockUserUseCase.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(nil, nil)
		mockUserUseCase.EXPECT().CreateObserverUser(gomock.Any(), gomock.Any()).Return(&user, nil)
		d = mock_middleware.Dependencies{UseCase: mockUserUseCase}
		router = configureRoutes(d)
		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusOK, w.Code, "Expected response code %d, received %d", http.StatusOK, w.Code)
	})

	t.Run("unsuccessful creation", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest("POST", "/where/are/they/users/observer", bytes.NewBuffer([]byte(
			`{
				"user": {
					"name": "Maria",
					"last_name": "Dominguez",
					"id_number": "12345678",
					"username": "mdominguez",
					"password": "mdominguez1234",
					"email": "mdominguez@mail.com",
					"enabled": true,
					"type": "observer"
				}
			}`,
		)))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
		mockUserUseCase.EXPECT().FindByUsername(gomock.Any(), gomock.Any()).Return(nil, nil)
		mockUserUseCase.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(nil, nil)
		mockUserUseCase.EXPECT().CreateObserverUser(gomock.Any(), gomock.Any()).Return(nil,
			web.ErrInternalServerError)
		d = mock_middleware.Dependencies{UseCase: mockUserUseCase}
		router = configureRoutes(d)
		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusInternalServerError, w.Code, "Expected response code %d, received %d",
			http.StatusOK, w.Code)
	})
}
