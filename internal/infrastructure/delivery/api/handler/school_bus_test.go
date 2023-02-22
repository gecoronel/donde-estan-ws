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

var sb = model.SchoolBus{
	ID:           "1",
	LicensePlate: "11AAA22",
	Model:        "Master",
	Brand:        "Renault",
	License:      "111",
	CreatedAt:    "2023-02-18 17:09:33",
	UpdatedAt:    "2023-02-18 17:09:33",
}

func TestGetSchoolBus(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	d := mock_middleware.Dependencies{}
	router := configureRoutes(d)

	t.Run("error getting school bus", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "/where/are/they/school-bus/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		r.Header.Set("Content-Type", "application/json")
		mockSchoolBusUseCase := mock_usecase.NewMockSchoolBusUseCase(m)
		mockSchoolBusUseCase.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, web.ErrInternalServerError)
		d = mock_middleware.Dependencies{SchoolBusUseCase: mockSchoolBusUseCase}
		router = configureRoutes(d)
		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusInternalServerError, w.Code, "Expected response code %d, received %d",
			http.StatusInternalServerError, w.Code)
	})

	t.Run("not found error getting school bus", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "/where/are/they/school-bus/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		mockSchoolBusUseCase := mock_usecase.NewMockSchoolBusUseCase(m)
		mockSchoolBusUseCase.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, web.ErrNotFound)
		d = mock_middleware.Dependencies{SchoolBusUseCase: mockSchoolBusUseCase}
		router = configureRoutes(d)
		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusNotFound, w.Code, "Expected response code %d, received %d",
			http.StatusOK, w.Code)
	})

	t.Run("successful get school bus", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "/where/are/they/school-bus/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		mockSchoolBusUseCase := mock_usecase.NewMockSchoolBusUseCase(m)
		mockSchoolBusUseCase.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&sb, nil)
		d = mock_middleware.Dependencies{SchoolBusUseCase: mockSchoolBusUseCase}
		router = configureRoutes(d)
		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusOK, w.Code, "Expected response code %d, received %d", http.StatusOK, w.Code)
	})
}

func TestSaveSchoolBus(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	d := mock_middleware.Dependencies{}
	router := configureRoutes(d)

	t.Run("bad request error saving school bus", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodPost, "/where/are/they/school-bus", bytes.NewBuffer([]byte(`{
			"id": 12345,
			"license_plate": "11AAA55",
			"model": "Fiat",
			"brand": "Ducato",
			"license": "555",
			"updated_at": "2023-02-22 11:55:14"
		}`)))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusBadRequest, w.Code, "Expected response code %d, received %d",
			http.StatusOK, w.Code)
	})

	t.Run("validate error saving school bus", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodPost, "/where/are/they/school-bus", bytes.NewBuffer([]byte(`{
			"id": "0000-0000-0005",
			"model": "Fiat",
			"brand": "Ducato",
			"license": "555",
			"updated_at": "2023-02-22 11:55:14"
		}`)))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusBadRequest, w.Code, "Expected response code %d, received %d",
			http.StatusBadRequest, w.Code)
	})

	t.Run("error getting school bus", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodPost, "/where/are/they/school-bus", bytes.NewBuffer([]byte(`{
			"id": "0000-0000-0005",
			"license_plate": "11AAA55",
			"model": "Fiat",
			"brand": "Ducato",
			"license": "555",
			"updated_at": "2023-02-22 11:55:14"
		}`)))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		mockSchoolBusUseCase := mock_usecase.NewMockSchoolBusUseCase(m)
		mockSchoolBusUseCase.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil, web.ErrInternalServerError)
		d = mock_middleware.Dependencies{SchoolBusUseCase: mockSchoolBusUseCase}
		router = configureRoutes(d)
		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusInternalServerError, w.Code, "Expected response code %d, received %d",
			http.StatusInternalServerError, w.Code)
	})

	t.Run("successful get school bus", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodPost, "/where/are/they/school-bus", bytes.NewBuffer([]byte(`{
			"id": "0000-0000-0005",
			"license_plate": "11AAA55",
			"model": "Fiat",
			"brand": "Ducato",
			"license": "555",
			"updated_at": "2023-02-22 11:55:14"
		}`)))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		mockSchoolBusUseCase := mock_usecase.NewMockSchoolBusUseCase(m)
		mockSchoolBusUseCase.EXPECT().Save(gomock.Any(), gomock.Any()).Return(&sb, nil)
		d = mock_middleware.Dependencies{SchoolBusUseCase: mockSchoolBusUseCase}
		router = configureRoutes(d)
		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusOK, w.Code, "Expected response code %d, received %d", http.StatusOK, w.Code)
	})
}

func TestUpdateSchoolBus(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	d := mock_middleware.Dependencies{}
	router := configureRoutes(d)

	t.Run("bad request error updating school bus", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodPut, "/where/are/they/school-bus", bytes.NewBuffer([]byte(`{
			"id": 12345,
			"license_plate": "11AAA55",
			"model": "Fiat",
			"brand": "Ducato",
			"license": "555",
			"updated_at": "2023-02-22 11:55:14"
		}`)))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusBadRequest, w.Code, "Expected response code %d, received %d",
			http.StatusOK, w.Code)
	})

	t.Run("validate error updating school bus", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodPut, "/where/are/they/school-bus", bytes.NewBuffer([]byte(`{
			"id": "0000-0000-0005",
			"model": "Fiat",
			"brand": "Ducato",
			"license": "555",
			"updated_at": "2023-02-22 11:55:14"
		}`)))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusBadRequest, w.Code, "Expected response code %d, received %d",
			http.StatusBadRequest, w.Code)
	})

	t.Run("error updating school bus", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodPut, "/where/are/they/school-bus", bytes.NewBuffer([]byte(`{
			"id": "0000-0000-0005",
			"license_plate": "11AAA55",
			"model": "Fiat",
			"brand": "Ducato",
			"license": "555",
			"updated_at": "2023-02-22 11:55:14"
		}`)))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		mockSchoolBusUseCase := mock_usecase.NewMockSchoolBusUseCase(m)
		mockSchoolBusUseCase.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, web.ErrInternalServerError)
		d = mock_middleware.Dependencies{SchoolBusUseCase: mockSchoolBusUseCase}
		router = configureRoutes(d)
		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusInternalServerError, w.Code, "Expected response code %d, received %d",
			http.StatusInternalServerError, w.Code)
	})

	t.Run("successful update school bus", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodPut, "/where/are/they/school-bus", bytes.NewBuffer([]byte(`{
			"id": "0000-0000-0005",
			"license_plate": "11AAA55",
			"model": "Fiat",
			"brand": "Ducato",
			"license": "555",
			"updated_at": "2023-02-22 11:55:14"
		}`)))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		mockSchoolBusUseCase := mock_usecase.NewMockSchoolBusUseCase(m)
		mockSchoolBusUseCase.EXPECT().Update(gomock.Any(), gomock.Any()).Return(&sb, nil)
		d = mock_middleware.Dependencies{SchoolBusUseCase: mockSchoolBusUseCase}
		router = configureRoutes(d)
		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusOK, w.Code, "Expected response code %d, received %d", http.StatusOK, w.Code)
	})
}

func TestDeleteSchoolBus(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	d := mock_middleware.Dependencies{}
	router := configureRoutes(d)

	t.Run("error deleting school bus", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodDelete, "/where/are/they/school-bus/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		r.Header.Set("Content-Type", "application/json")
		mockSchoolBusUseCase := mock_usecase.NewMockSchoolBusUseCase(m)
		mockSchoolBusUseCase.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(web.ErrInternalServerError)
		d = mock_middleware.Dependencies{SchoolBusUseCase: mockSchoolBusUseCase}
		router = configureRoutes(d)
		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusInternalServerError, w.Code, "Expected response code %d, received %d",
			http.StatusInternalServerError, w.Code)
	})

	t.Run("not found error deleting school bus", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodDelete, "/where/are/they/school-bus/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		mockSchoolBusUseCase := mock_usecase.NewMockSchoolBusUseCase(m)
		mockSchoolBusUseCase.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(web.ErrNotFound)
		d = mock_middleware.Dependencies{SchoolBusUseCase: mockSchoolBusUseCase}
		router = configureRoutes(d)
		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusNotFound, w.Code, "Expected response code %d, received %d",
			http.StatusOK, w.Code)
	})

	t.Run("successful delete school bus", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodDelete, "/where/are/they/school-bus/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")

		mockSchoolBusUseCase := mock_usecase.NewMockSchoolBusUseCase(m)
		mockSchoolBusUseCase.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
		d = mock_middleware.Dependencies{SchoolBusUseCase: mockSchoolBusUseCase}
		router = configureRoutes(d)
		router.ServeHTTP(w, r)
		assert.Equalf(t, http.StatusOK, w.Code, "Expected response code %d, received %d", http.StatusOK, w.Code)
	})
}
