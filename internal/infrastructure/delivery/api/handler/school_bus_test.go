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

	tests := []struct {
		name         string
		mock         func() *mock_usecase.MockSchoolBusUseCase
		path         string
		expectedCode int
	}{
		{
			name: "error getting school bus",
			mock: func() *mock_usecase.MockSchoolBusUseCase {
				mockSchoolBusUseCase := mock_usecase.NewMockSchoolBusUseCase(m)
				mockSchoolBusUseCase.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockSchoolBusUseCase
			},
			path:         "/where/are/they/school-buses/1",
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "not found error getting school bus",
			mock: func() *mock_usecase.MockSchoolBusUseCase {
				mockSchoolBusUseCase := mock_usecase.NewMockSchoolBusUseCase(m)
				mockSchoolBusUseCase.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, web.ErrNotFound)
				return mockSchoolBusUseCase
			},
			path:         "/where/are/they/school-buses/1",
			expectedCode: http.StatusNotFound,
		},
		{
			name: "successful get school bus",
			mock: func() *mock_usecase.MockSchoolBusUseCase {
				mockSchoolBusUseCase := mock_usecase.NewMockSchoolBusUseCase(m)
				mockSchoolBusUseCase.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&sb, nil)
				return mockSchoolBusUseCase
			},
			path:         "/where/are/they/school-buses/1",
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

			d := mock_middleware.Dependencies{SchoolBusUseCase: test.mock()}
			router := configureRoutes(d)
			router.ServeHTTP(w, r)
			assert.Equalf(t, test.expectedCode, w.Code, "Expected code %v, received %v", test.expectedCode, w.Code)
		})
	}
}

func TestSaveSchoolBus(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name         string
		mock         func() *mock_usecase.MockSchoolBusUseCase
		path         string
		body         string
		expectedCode int
	}{
		{
			name: "bad request error saving school bus",
			mock: func() *mock_usecase.MockSchoolBusUseCase {
				mockSchoolBusUseCase := mock_usecase.NewMockSchoolBusUseCase(m)
				return mockSchoolBusUseCase
			},
			path: "/where/are/they/school-buses",
			body: `{
			"id": 12345,
			"license_plate": "11AAA55",
			"model": "Fiat",
			"brand": "Ducato",
			"license": "555",
			"updated_at": "2023-02-22 11:55:14"
		}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "validate error saving school bus",
			mock: func() *mock_usecase.MockSchoolBusUseCase {
				mockSchoolBusUseCase := mock_usecase.NewMockSchoolBusUseCase(m)
				return mockSchoolBusUseCase
			},
			path: "/where/are/they/school-buses",
			body: `{
			"id": "0000-0000-0005",
			"model": "Fiat",
			"brand": "Ducato",
			"license": "555",
			"updated_at": "2023-02-22 11:55:14"
		}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "error saving school bus",
			mock: func() *mock_usecase.MockSchoolBusUseCase {
				mockSchoolBusUseCase := mock_usecase.NewMockSchoolBusUseCase(m)
				mockSchoolBusUseCase.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockSchoolBusUseCase
			},
			path: "/where/are/they/school-buses",
			body: `{
			"id": "0000-0000-0005",
			"license_plate": "11AAA55",
			"model": "Fiat",
			"brand": "Ducato",
			"license": "555",
			"updated_at": "2023-02-22 11:55:14"
		}`,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "successful save school bus",
			mock: func() *mock_usecase.MockSchoolBusUseCase {
				mockSchoolBusUseCase := mock_usecase.NewMockSchoolBusUseCase(m)
				mockSchoolBusUseCase.EXPECT().Save(gomock.Any(), gomock.Any()).Return(&sb, nil)
				return mockSchoolBusUseCase
			},
			path: "/where/are/they/school-buses",
			body: `{
			"id": "0000-0000-0005",
			"license_plate": "11AAA55",
			"model": "Fiat",
			"brand": "Ducato",
			"license": "555",
			"updated_at": "2023-02-22 11:55:14"
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

			d := mock_middleware.Dependencies{SchoolBusUseCase: test.mock()}
			router := configureRoutes(d)
			router.ServeHTTP(w, r)
			assert.Equalf(t, test.expectedCode, w.Code, "Expected code %v, received %v", test.expectedCode, w.Code)
		})
	}
}

func TestUpdateSchoolBus(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name         string
		mock         func() *mock_usecase.MockSchoolBusUseCase
		path         string
		body         string
		expectedCode int
	}{
		{
			name: "bad request error updating school bus",
			mock: func() *mock_usecase.MockSchoolBusUseCase {
				mockSchoolBusUseCase := mock_usecase.NewMockSchoolBusUseCase(m)
				return mockSchoolBusUseCase
			},
			path: "/where/are/they/school-buses",
			body: `{
			"id": 12345,
			"license_plate": "11AAA55",
			"model": "Fiat",
			"brand": "Ducato",
			"license": "555",
			"updated_at": "2023-02-22 11:55:14"
		}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "validate error updating school bus",
			mock: func() *mock_usecase.MockSchoolBusUseCase {
				mockSchoolBusUseCase := mock_usecase.NewMockSchoolBusUseCase(m)
				return mockSchoolBusUseCase
			},
			path: "/where/are/they/school-buses",
			body: `{
			"id": "0000-0000-0005",
			"model": "Fiat",
			"brand": "Ducato",
			"license": "555",
			"updated_at": "2023-02-22 11:55:14"
		}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "error updating school bus",
			mock: func() *mock_usecase.MockSchoolBusUseCase {
				mockSchoolBusUseCase := mock_usecase.NewMockSchoolBusUseCase(m)
				mockSchoolBusUseCase.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockSchoolBusUseCase
			},
			path: "/where/are/they/school-buses",
			body: `{
			"id": "0000-0000-0005",
			"license_plate": "11AAA55",
			"model": "Fiat",
			"brand": "Ducato",
			"license": "555",
			"updated_at": "2023-02-22 11:55:14"
		}`,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "successful update school bus",
			mock: func() *mock_usecase.MockSchoolBusUseCase {
				mockSchoolBusUseCase := mock_usecase.NewMockSchoolBusUseCase(m)
				mockSchoolBusUseCase.EXPECT().Update(gomock.Any(), gomock.Any()).Return(&sb, nil)
				return mockSchoolBusUseCase
			},
			path: "/where/are/they/school-buses",
			body: `{
			"id": "0000-0000-0005",
			"license_plate": "11AAA55",
			"model": "Fiat",
			"brand": "Ducato",
			"license": "555",
			"updated_at": "2023-02-22 11:55:14"
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

			d := mock_middleware.Dependencies{SchoolBusUseCase: test.mock()}
			router := configureRoutes(d)
			router.ServeHTTP(w, r)
			assert.Equalf(t, test.expectedCode, w.Code, "Expected code %v, received %v", test.expectedCode, w.Code)
		})
	}
}

func TestDeleteSchoolBus(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name         string
		mock         func() *mock_usecase.MockSchoolBusUseCase
		path         string
		expectedCode int
	}{
		{
			name: "error deleting school bus",
			mock: func() *mock_usecase.MockSchoolBusUseCase {
				mockSchoolBusUseCase := mock_usecase.NewMockSchoolBusUseCase(m)
				mockSchoolBusUseCase.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(web.ErrInternalServerError)
				return mockSchoolBusUseCase
			},
			path:         "/where/are/they/school-buses/1",
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "not found error deleting school bus",
			mock: func() *mock_usecase.MockSchoolBusUseCase {
				mockSchoolBusUseCase := mock_usecase.NewMockSchoolBusUseCase(m)
				mockSchoolBusUseCase.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(web.ErrNotFound)
				return mockSchoolBusUseCase
			},
			path:         "/where/are/they/school-buses/1",
			expectedCode: http.StatusNotFound,
		},
		{
			name: "successful delete school bus",
			mock: func() *mock_usecase.MockSchoolBusUseCase {
				mockSchoolBusUseCase := mock_usecase.NewMockSchoolBusUseCase(m)
				mockSchoolBusUseCase.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
				return mockSchoolBusUseCase
			},
			path:         "/where/are/they/school-buses/1",
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

			d := mock_middleware.Dependencies{SchoolBusUseCase: test.mock()}
			router := configureRoutes(d)
			router.ServeHTTP(w, r)
			assert.Equalf(t, test.expectedCode, w.Code, "Expected code %v, received %v", test.expectedCode, w.Code)
		})
	}
}
