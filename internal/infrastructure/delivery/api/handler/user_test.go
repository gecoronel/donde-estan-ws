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

	tests := []struct {
		name         string
		mock         func() *mock_usecase.MockUserUseCase
		path         string
		expectedCode int
	}{
		{
			name: "error getting user",
			mock: func() *mock_usecase.MockUserUseCase {
				mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
				mockUserUseCase.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockUserUseCase
			},
			path:         "/where/are/they/users/1",
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "not found error getting user",
			mock: func() *mock_usecase.MockUserUseCase {
				mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
				mockUserUseCase.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, web.ErrNotFound)
				return mockUserUseCase
			},
			path:         "/where/are/they/users/10",
			expectedCode: http.StatusNotFound,
		},
		{
			name: "bad request error getting user",
			mock: func() *mock_usecase.MockUserUseCase {
				mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
				return mockUserUseCase
			},
			path:         "/where/are/they/users/invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "successful get user",
			mock: func() *mock_usecase.MockUserUseCase {
				mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
				mockUserUseCase.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&user, nil)
				return mockUserUseCase
			},
			path:         "/where/are/they/users/1",
			expectedCode: http.StatusOK,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, err := http.NewRequest(http.MethodGet, test.path, nil)
			if err != nil {
				t.Fatal(err)
			}
			r.Header.Set("Content-Type", "application/json")

			d := mock_middleware.Dependencies{UserUseCase: test.mock()}
			router := configureRoutes(d)
			router.ServeHTTP(w, r)
			assert.Equalf(t, test.expectedCode, w.Code, "Expected code %v, received %v", test.expectedCode, w.Code)
		})
	}
}

func TestLogin(t *testing.T) {
	observedUser := model.ObservedUser{
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
	odu := model.NewObservedUser(&observedUser)

	observerUser := model.ObserverUser{
		User: model.User{
			ID:        2,
			Name:      "Jose",
			LastName:  "Perez",
			IDNumber:  "12345678",
			Username:  "joseperez",
			Password:  "joseperez1234",
			Email:     "joseperez@mail.com",
			Enabled:   true,
			Type:      "observer",
			CreatedAt: "2022-12-10 17:49:30",
			UpdatedAt: "2022-12-10 17:49:30",
		},
	}
	oru := model.NewObserverUser(&observerUser)

	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name         string
		mock         func() *mock_usecase.MockUserUseCase
		path         string
		body         string
		expectedCode int
	}{
		{
			name: "bad request error for handler login",
			mock: func() *mock_usecase.MockUserUseCase {
				mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
				return mockUserUseCase
			},
			path:         "/where/are/they/users/login",
			body:         `invalid`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "bad request error in validation for handler login",
			mock: func() *mock_usecase.MockUserUseCase {
				mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
				return mockUserUseCase
			},
			path:         "/where/are/they/users/login",
			body:         `{"username": "jperez"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "not found error in handler login",
			mock: func() *mock_usecase.MockUserUseCase {
				mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
				mockUserUseCase.EXPECT().Login(gomock.Any(), gomock.Any()).Return(nil, web.ErrNotFound)
				return mockUserUseCase
			},
			path:         "/where/are/they/users/login",
			body:         `{"username": "jp", "password": "jperez1234"}`,
			expectedCode: http.StatusNotFound,
		},
		{
			name: "successful login",
			mock: func() *mock_usecase.MockUserUseCase {
				mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
				mockUserUseCase.EXPECT().Login(gomock.Any(), gomock.Any()).Return(odu, nil)
				return mockUserUseCase
			},
			path:         "/where/are/they/users/login",
			body:         `{"username": "jperez", "password": "jperez1234"}`,
			expectedCode: http.StatusOK,
		},
		{
			name: "successful login",
			mock: func() *mock_usecase.MockUserUseCase {
				mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
				mockUserUseCase.EXPECT().Login(gomock.Any(), gomock.Any()).Return(oru, nil)
				return mockUserUseCase
			},
			path:         "/where/are/they/users/login",
			body:         `{"username": "joseperez", "password": "joseperez1234"}`,
			expectedCode: http.StatusOK,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, err := http.NewRequest(http.MethodPost, test.path, bytes.NewBuffer([]byte(test.body)))
			if err != nil {
				t.Fatal(err)
			}
			r.Header.Set("Content-Type", "application/json")

			d := mock_middleware.Dependencies{UserUseCase: test.mock()}
			router := configureRoutes(d)
			router.ServeHTTP(w, r)
			assert.Equalf(t, test.expectedCode, w.Code, "Expected code %v, received %v", test.expectedCode, w.Code)
		})
	}
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

	body := `{
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
			}`

	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name         string
		mock         func() *mock_usecase.MockUserUseCase
		path         string
		body         string
		expectedCode int
	}{
		{
			name: "bad request for handler create observed user",
			mock: func() *mock_usecase.MockUserUseCase {
				mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
				return mockUserUseCase
			},
			path:         "/where/are/they/users/observed",
			body:         `invalid`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "bad request in login validation",
			mock: func() *mock_usecase.MockUserUseCase {
				mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
				return mockUserUseCase
			},
			path:         "/where/are/they/users/observed",
			body:         `{"user": {"username": "jperez"}}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "username conflict error",
			mock: func() *mock_usecase.MockUserUseCase {
				mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
				mockUserUseCase.EXPECT().FindByUsername(gomock.Any(), gomock.Any()).Return(&user.User, nil)
				return mockUserUseCase
			},
			path:         "/where/are/they/users/observed",
			body:         body,
			expectedCode: http.StatusConflict,
		},
		{
			name: "email conflict error",
			mock: func() *mock_usecase.MockUserUseCase {
				mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
				mockUserUseCase.EXPECT().FindByUsername(gomock.Any(), gomock.Any()).Return(nil, nil)
				mockUserUseCase.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(&user.User, nil)
				return mockUserUseCase
			},
			path:         "/where/are/they/users/observed",
			body:         body,
			expectedCode: http.StatusConflict,
		},
		{
			name: "successful creation",
			mock: func() *mock_usecase.MockUserUseCase {
				mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
				mockUserUseCase.EXPECT().FindByUsername(gomock.Any(), gomock.Any()).Return(nil, nil)
				mockUserUseCase.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(nil, nil)
				mockUserUseCase.EXPECT().CreateObservedUser(gomock.Any(), gomock.Any()).Return(&user, nil)
				return mockUserUseCase
			},
			path:         "/where/are/they/users/observed",
			body:         body,
			expectedCode: http.StatusOK,
		},
		{
			name: "unsuccessful creation",
			mock: func() *mock_usecase.MockUserUseCase {
				mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
				mockUserUseCase.EXPECT().FindByUsername(gomock.Any(), gomock.Any()).Return(nil, nil)
				mockUserUseCase.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(nil, nil)
				mockUserUseCase.EXPECT().CreateObservedUser(gomock.Any(), gomock.Any()).Return(nil,
					web.ErrInternalServerError)
				return mockUserUseCase
			},
			path:         "/where/are/they/users/observed",
			body:         body,
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, err := http.NewRequest(http.MethodPost, test.path, bytes.NewBuffer([]byte(test.body)))
			if err != nil {
				t.Fatal(err)
			}
			r.Header.Set("Content-Type", "application/json")

			d := mock_middleware.Dependencies{UserUseCase: test.mock()}
			router := configureRoutes(d)
			router.ServeHTTP(w, r)
			assert.Equalf(t, test.expectedCode, w.Code, "Expected code %v, received %v", test.expectedCode, w.Code)
		})
	}
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

	body := `{
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
			}`

	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name         string
		mock         func() *mock_usecase.MockUserUseCase
		path         string
		body         string
		expectedCode int
	}{
		{
			name: "bad request for handler create observed user",
			mock: func() *mock_usecase.MockUserUseCase {
				mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
				return mockUserUseCase
			},
			path:         "/where/are/they/users/observer",
			body:         `invalid`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "bad request in login validation",
			mock: func() *mock_usecase.MockUserUseCase {
				mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
				return mockUserUseCase
			},
			path:         "/where/are/they/users/observer",
			body:         `{"user": {"username": "jperez"}}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "username conflict error",
			mock: func() *mock_usecase.MockUserUseCase {
				mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
				mockUserUseCase.EXPECT().FindByUsername(gomock.Any(), gomock.Any()).Return(&user.User, nil)
				return mockUserUseCase
			},
			path:         "/where/are/they/users/observer",
			body:         body,
			expectedCode: http.StatusConflict,
		},
		{
			name: "email conflict error",
			mock: func() *mock_usecase.MockUserUseCase {
				mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
				mockUserUseCase.EXPECT().FindByUsername(gomock.Any(), gomock.Any()).Return(nil, nil)
				mockUserUseCase.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(&user.User, nil)
				return mockUserUseCase
			},
			path:         "/where/are/they/users/observer",
			body:         body,
			expectedCode: http.StatusConflict,
		},
		{
			name: "successful creation",
			mock: func() *mock_usecase.MockUserUseCase {
				mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
				mockUserUseCase.EXPECT().FindByUsername(gomock.Any(), gomock.Any()).Return(nil, nil)
				mockUserUseCase.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(nil, nil)
				mockUserUseCase.EXPECT().CreateObserverUser(gomock.Any(), gomock.Any()).Return(&user, nil)
				return mockUserUseCase
			},
			path:         "/where/are/they/users/observer",
			body:         body,
			expectedCode: http.StatusOK,
		},
		{
			name: "unsuccessful creation",
			mock: func() *mock_usecase.MockUserUseCase {
				mockUserUseCase := mock_usecase.NewMockUserUseCase(m)
				mockUserUseCase.EXPECT().FindByUsername(gomock.Any(), gomock.Any()).Return(nil, nil)
				mockUserUseCase.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(nil, nil)
				mockUserUseCase.EXPECT().CreateObserverUser(gomock.Any(), gomock.Any()).Return(nil,
					web.ErrInternalServerError)
				return mockUserUseCase
			},
			path:         "/where/are/they/users/observer",
			body:         body,
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, err := http.NewRequest(http.MethodPost, test.path, bytes.NewBuffer([]byte(test.body)))
			if err != nil {
				t.Fatal(err)
			}
			r.Header.Set("Content-Type", "application/json")

			d := mock_middleware.Dependencies{UserUseCase: test.mock()}
			router := configureRoutes(d)
			router.ServeHTTP(w, r)
			assert.Equalf(t, test.expectedCode, w.Code, "Expected code %v, received %v", test.expectedCode, w.Code)
		})
	}
}
