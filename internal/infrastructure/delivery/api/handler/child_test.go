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

var c = model.Child{
	ID:              1,
	Name:            "Pilar",
	LastName:        "Dominguez",
	SchoolName:      "La Salle",
	SchoolStartTime: "8:00",
	SchoolEndTime:   "12:00",
	CreatedAt:       "2023-02-18 17:09:33",
	UpdatedAt:       "2023-02-18 17:09:33",
	ObserverUserID:  uint64(10),
}

func TestGetChild(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name         string
		mock         func() *mock_usecase.MockChildUseCase
		path         string
		expectedCode int
	}{
		{
			name: "bad request error getting child",
			mock: func() *mock_usecase.MockChildUseCase {
				mockChildUseCase := mock_usecase.NewMockChildUseCase(m)
				return mockChildUseCase
			},
			path:         "/where/are/they/children/invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "error getting child",
			mock: func() *mock_usecase.MockChildUseCase {
				mockChildUseCase := mock_usecase.NewMockChildUseCase(m)
				mockChildUseCase.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockChildUseCase
			},
			path:         "/where/are/they/children/1",
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "not found error getting child",
			mock: func() *mock_usecase.MockChildUseCase {
				mockChildUseCase := mock_usecase.NewMockChildUseCase(m)
				mockChildUseCase.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, web.ErrNotFound)
				return mockChildUseCase
			},
			path:         "/where/are/they/children/1",
			expectedCode: http.StatusNotFound,
		},
		{
			name: "successful get child",
			mock: func() *mock_usecase.MockChildUseCase {
				mockChildUseCase := mock_usecase.NewMockChildUseCase(m)
				mockChildUseCase.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&c, nil)
				return mockChildUseCase
			},
			path:         "/where/are/they/children/1",
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

			d := mock_middleware.Dependencies{ChildUseCase: test.mock()}
			router := configureRoutes(d)
			router.ServeHTTP(w, r)
			assert.Equalf(t, test.expectedCode, w.Code, "Expected code %v, received %v", test.expectedCode, w.Code)
		})
	}
}

func TestSaveChild(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name         string
		mock         func() *mock_usecase.MockChildUseCase
		path         string
		body         string
		expectedCode int
	}{
		{
			name: "bad request error saving child",
			mock: func() *mock_usecase.MockChildUseCase {
				mockChildUseCase := mock_usecase.NewMockChildUseCase(m)
				return mockChildUseCase
			},
			path: "/where/are/they/children",
			body: `{
				"name":              "Pilar",
				"last_name":         "Dominguez",
				"school_name":       "La Salle",
				"school_start_time": "8:00",
				"school_end_time":   "12:00",
				"observer_user_id":  "10"
			}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "validate error saving child for include id",
			mock: func() *mock_usecase.MockChildUseCase {
				mockChildUseCase := mock_usecase.NewMockChildUseCase(m)
				return mockChildUseCase
			},
			path: "/where/are/they/children",
			body: `{
				"id":                1,
				"name":              "Pilar",
				"last_name":         "Dominguez",
				"school_name":       "La Salle",
				"school_start_time": "8:00",
				"school_end_time":   "12:00",
				"observer_user_id":  10
			}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "validate error saving child for not include observer user id",
			mock: func() *mock_usecase.MockChildUseCase {
				mockChildUseCase := mock_usecase.NewMockChildUseCase(m)
				return mockChildUseCase
			},
			path: "/where/are/they/children",
			body: `{
				"id":                1,
				"name":              "Pilar",
				"last_name":         "Dominguez",
				"school_name":       "La Salle",
				"school_start_time": "8:00",
				"school_end_time":   "12:00"
			}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "validate error saving child",
			mock: func() *mock_usecase.MockChildUseCase {
				mockChildUseCase := mock_usecase.NewMockChildUseCase(m)
				return mockChildUseCase
			},
			path: "/where/are/they/children",
			body: `{
				"name":             "Pilar",
				"observer_user_id":  10
			}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "error saving child",
			mock: func() *mock_usecase.MockChildUseCase {
				mockChildUseCase := mock_usecase.NewMockChildUseCase(m)
				mockChildUseCase.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockChildUseCase
			},
			path: "/where/are/they/children",
			body: `{
				"name":              "Pilar",
				"last_name":         "Dominguez",
				"school_name":       "La Salle",
				"school_start_time": "8:00",
				"school_end_time":   "12:00",
				"observer_user_id":  10
			}`,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "successful save child",
			mock: func() *mock_usecase.MockChildUseCase {
				mockChildUseCase := mock_usecase.NewMockChildUseCase(m)
				mockChildUseCase.EXPECT().Save(gomock.Any(), gomock.Any()).Return(&c, nil)
				return mockChildUseCase
			},
			path: "/where/are/they/children",
			body: `{
				"name":              "Pilar",
				"last_name":         "Dominguez",
				"school_name":       "La Salle",
				"school_start_time": "8:00",
				"school_end_time":   "12:00",
				"observer_user_id":  10
			}`,
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

			d := mock_middleware.Dependencies{ChildUseCase: test.mock()}
			router := configureRoutes(d)
			router.ServeHTTP(w, r)
			assert.Equalf(t, test.expectedCode, w.Code, "Expected code %v, received %v", test.expectedCode, w.Code)
		})
	}
}

func TestUpdateChild(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name         string
		mock         func() *mock_usecase.MockChildUseCase
		path         string
		body         string
		expectedCode int
	}{
		{
			name: "bad request error updating child",
			mock: func() *mock_usecase.MockChildUseCase {
				mockChildUseCase := mock_usecase.NewMockChildUseCase(m)
				return mockChildUseCase
			},
			path: "/where/are/they/children",
			body: `{
				"id":                1,
				"name":              "Pilar",
				"last_name":         "Dominguez",
				"school_name":       "La Salle",
				"school_start_time": "8:00",
				"school_end_time":   "12:00",
				"created_at":        "2023-02-18 17:09:33",
				"updated_at":        "2023-02-18 17:09:33",
				"observer_user_id":  "10"
			}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "validate error saving child for not include id field",
			mock: func() *mock_usecase.MockChildUseCase {
				mockChildUseCase := mock_usecase.NewMockChildUseCase(m)
				return mockChildUseCase
			},
			path: "/where/are/they/children",
			body: `{
				"name":              "Pilar",
				"last_name":         "Dominguez",
				"school_name":       "La Salle",
				"school_start_time": "8:00",
				"school_end_time":   "12:00",
				"created_at":        "2023-02-18 17:09:33",
				"updated_at":        "2023-02-18 17:09:33",
				"observer_user_id":  10
			}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "validate error updating child",
			mock: func() *mock_usecase.MockChildUseCase {
				mockChildUseCase := mock_usecase.NewMockChildUseCase(m)
				return mockChildUseCase
			},
			path: "/where/are/they/children",
			body: `{
				"id":                1,
				"name":              "Pilar",
				"school_end_time":   "12:00",
				"created_at":        "2023-02-18 17:09:33",
				"updated_at":        "2023-02-18 17:09:33",
				"observer_user_id":  10
			}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "error updating child",
			mock: func() *mock_usecase.MockChildUseCase {
				mockChildUseCase := mock_usecase.NewMockChildUseCase(m)
				mockChildUseCase.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockChildUseCase
			},
			path: "/where/are/they/children",
			body: `{
				"id":                1,
				"name":              "Pilar",
				"last_name":         "Dominguez",
				"school_name":       "La Salle",
				"school_start_time": "8:00",
				"school_end_time":   "12:00",
				"created_at":        "2023-02-18 17:09:33",
				"updated_at":        "2023-02-18 17:09:33",
				"observer_user_id":  10
			}`,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "successful update child",
			mock: func() *mock_usecase.MockChildUseCase {
				mockChildUseCase := mock_usecase.NewMockChildUseCase(m)
				mockChildUseCase.EXPECT().Update(gomock.Any(), gomock.Any()).Return(&c, nil)
				return mockChildUseCase
			},
			path: "/where/are/they/children",
			body: `{
				"id":                1,
				"name":              "Pilar",
				"last_name":         "Dominguez",
				"school_name":       "La Salle",
				"school_start_time": "8:00",
				"school_end_time":   "12:00",
				"created_at":        "2023-02-18 17:09:33",
				"updated_at":        "2023-02-18 17:09:33",
				"observer_user_id":  10
			}`,
			expectedCode: http.StatusOK,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, err := http.NewRequest(http.MethodPut, test.path, bytes.NewBuffer([]byte(test.body)))
			if err != nil {
				t.Fatal(err)
			}
			r.Header.Set("Content-Type", "application/json")

			d := mock_middleware.Dependencies{ChildUseCase: test.mock()}
			router := configureRoutes(d)
			router.ServeHTTP(w, r)
			assert.Equalf(t, test.expectedCode, w.Code, "Expected code %v, received %v", test.expectedCode, w.Code)
		})
	}
}

func TestDeleteChild(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name         string
		mock         func() *mock_usecase.MockChildUseCase
		path         string
		expectedCode int
	}{
		{
			name: "bad request error deleting child",
			mock: func() *mock_usecase.MockChildUseCase {
				mockChildUseCase := mock_usecase.NewMockChildUseCase(m)
				return mockChildUseCase
			},
			path:         "/where/are/they/children/invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "error deleting child",
			mock: func() *mock_usecase.MockChildUseCase {
				mockChildUseCase := mock_usecase.NewMockChildUseCase(m)
				mockChildUseCase.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(web.ErrInternalServerError)
				return mockChildUseCase
			},
			path:         "/where/are/they/children/1",
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "not found error deleting child",
			mock: func() *mock_usecase.MockChildUseCase {
				mockChildUseCase := mock_usecase.NewMockChildUseCase(m)
				mockChildUseCase.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(web.ErrNotFound)
				return mockChildUseCase
			},
			path:         "/where/are/they/children/1",
			expectedCode: http.StatusNotFound,
		},
		{
			name: "successful delete child",
			mock: func() *mock_usecase.MockChildUseCase {
				mockChildUseCase := mock_usecase.NewMockChildUseCase(m)
				mockChildUseCase.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
				return mockChildUseCase
			},
			path:         "/where/are/they/children/1",
			expectedCode: http.StatusOK,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, err := http.NewRequest(http.MethodDelete, test.path, nil)
			if err != nil {
				t.Fatal(err)
			}
			r.Header.Set("Content-Type", "application/json")

			d := mock_middleware.Dependencies{ChildUseCase: test.mock()}
			router := configureRoutes(d)
			router.ServeHTTP(w, r)
			assert.Equalf(t, test.expectedCode, w.Code, "Expected code %v, received %v", test.expectedCode, w.Code)
		})
	}
}
