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

var a = model.Address{
	ID:             1,
	Name:           "Casa",
	Street:         "25 de Mayo",
	Number:         "1010",
	Floor:          "1",
	Apartment:      "A",
	ZipCode:        "3000",
	City:           "Santa Fe",
	State:          "Santa Fe",
	Country:        "Argentina",
	Latitude:       "60.0000121",
	Longitude:      "-19.23423",
	CreatedAt:      "2023-02-18 17:09:33",
	UpdatedAt:      "2023-02-18 17:09:33",
	ObserverUserID: uint64(10),
}

func TestGetAddress(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name         string
		mock         func() *mock_usecase.MockAddressUseCase
		path         string
		expectedCode int
	}{
		{
			name: "bad request error getting address",
			mock: func() *mock_usecase.MockAddressUseCase {
				mockAddressUseCase := mock_usecase.NewMockAddressUseCase(m)
				return mockAddressUseCase
			},
			path:         "/where/are/they/addresses/invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "error getting address",
			mock: func() *mock_usecase.MockAddressUseCase {
				mockAddressUseCase := mock_usecase.NewMockAddressUseCase(m)
				mockAddressUseCase.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockAddressUseCase
			},
			path:         "/where/are/they/addresses/1",
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "not found error getting address",
			mock: func() *mock_usecase.MockAddressUseCase {
				mockAddressUseCase := mock_usecase.NewMockAddressUseCase(m)
				mockAddressUseCase.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, web.ErrNotFound)
				return mockAddressUseCase
			},
			path:         "/where/are/they/addresses/1",
			expectedCode: http.StatusNotFound,
		},
		{
			name: "successful get address",
			mock: func() *mock_usecase.MockAddressUseCase {
				mockAddressUseCase := mock_usecase.NewMockAddressUseCase(m)
				mockAddressUseCase.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&a, nil)
				return mockAddressUseCase
			},
			path:         "/where/are/they/addresses/1",
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

			d := mock_middleware.Dependencies{AddressUseCase: test.mock()}
			router := configureRoutes(d)
			router.ServeHTTP(w, r)
			assert.Equalf(t, test.expectedCode, w.Code, "Expected code %v, received %v", test.expectedCode, w.Code)
		})
	}
}

func TestSaveAddress(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name         string
		mock         func() *mock_usecase.MockAddressUseCase
		path         string
		body         string
		expectedCode int
	}{
		{
			name: "bad request error saving address",
			mock: func() *mock_usecase.MockAddressUseCase {
				mockAddressUseCase := mock_usecase.NewMockAddressUseCase(m)
				return mockAddressUseCase
			},
			path: "/where/are/they/addresses",
			body: `{
				"name":             "Casa",
				"street":           "25 de Mayo",
				"number":           "1234",
				"floor":            "1",
				"apartment":        "A",
				"zip_code":         "3000",
				"city":             "Santa Fe",
				"state":            "Santa Fe",
				"country":          "Argentina",
				"latitude":         "60.0000121",
				"longitude":        "-19.23423",
				"updated_at":       "2023-02-22 11:55:14",
				"observer_user_id": "10"
			}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "validate error saving address for include id in body ",
			mock: func() *mock_usecase.MockAddressUseCase {
				mockAddressUseCase := mock_usecase.NewMockAddressUseCase(m)
				return mockAddressUseCase
			},
			path: "/where/are/they/addresses",
			body: `{
				"id":               1,
				"name":             "Casa",
				"street":           "25 de Mayo",
				"number":           "1234",
				"floor":            "1",
				"apartment":        "A",
				"zip_code":         "3000",
				"city":             "Santa Fe",
				"state":            "Santa Fe",
				"country":          "Argentina",
				"latitude":         "60.0000121",
				"longitude":        "-19.23423",
				"observer_user_id": 10
			}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "validate error saving address",
			mock: func() *mock_usecase.MockAddressUseCase {
				mockAddressUseCase := mock_usecase.NewMockAddressUseCase(m)
				return mockAddressUseCase
			},
			path: "/where/are/they/addresses",
			body: `{
				"name":             "Casa",
				"street":           "25 de Mayo",
				"country":          "Argentina",
				"latitude":         "60.0000121",
				"longitude":        "-19.23423",
				"observer_user_id": 10
			}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "error saving address",
			mock: func() *mock_usecase.MockAddressUseCase {
				mockAddressUseCase := mock_usecase.NewMockAddressUseCase(m)
				mockAddressUseCase.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockAddressUseCase
			},
			path: "/where/are/they/addresses",
			body: `{
				"name":             "Casa",
				"street":           "25 de Mayo",
				"number":           "1234",
				"floor":            "1",
				"apartment":        "A",
				"zip_code":         "3000",
				"city":             "Santa Fe",
				"state":            "Santa Fe",
				"country":          "Argentina",
				"latitude":         "60.0000121",
				"longitude":        "-19.23423",
				"updated_at":       "2023-02-22 11:55:14",
				"observer_user_id": 10
			}`,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "successful save address",
			mock: func() *mock_usecase.MockAddressUseCase {
				mockAddressUseCase := mock_usecase.NewMockAddressUseCase(m)
				mockAddressUseCase.EXPECT().Save(gomock.Any(), gomock.Any()).Return(&a, nil)
				return mockAddressUseCase
			},
			path: "/where/are/they/addresses",
			body: `{
				"name":             "Casa",
				"street":           "25 de Mayo",
				"number":           "1234",
				"floor":            "1",
				"apartment":        "A",
				"zip_code":         "3000",
				"city":             "Santa Fe",
				"state":            "Santa Fe",
				"country":          "Argentina",
				"latitude":         "60.0000121",
				"longitude":        "-19.23423",
				"updated_at":       "2023-02-22 11:55:14",
				"observer_user_id": 10
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

			d := mock_middleware.Dependencies{AddressUseCase: test.mock()}
			router := configureRoutes(d)
			router.ServeHTTP(w, r)
			assert.Equalf(t, test.expectedCode, w.Code, "Expected code %v, received %v", test.expectedCode, w.Code)
		})
	}
}

func TestUpdateAddress(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name         string
		mock         func() *mock_usecase.MockAddressUseCase
		path         string
		body         string
		expectedCode int
	}{
		{
			name: "bad request error updating address",
			mock: func() *mock_usecase.MockAddressUseCase {
				mockAddressUseCase := mock_usecase.NewMockAddressUseCase(m)
				return mockAddressUseCase
			},
			path: "/where/are/they/addresses",
			body: `{
				"id":               1,
				"name":             "Casa",
				"street":           "25 de Mayo",
				"number":           "1234",
				"floor":            "1",
				"apartment":        "A",
				"zip_code":         "3000",
				"city":             "Santa Fe",
				"state":            "Santa Fe",
				"country":          "Argentina",
				"latitude":         "60.0000121",
				"longitude":        "-19.23423",
				"updated_at":       "2023-02-22 11:55:14",
				"observer_user_id": "10"
			}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "validate error saving address for not include id in body ",
			mock: func() *mock_usecase.MockAddressUseCase {
				mockAddressUseCase := mock_usecase.NewMockAddressUseCase(m)
				return mockAddressUseCase
			},
			path: "/where/are/they/addresses",
			body: `{
			"name":             "Casa",
			"street":           "25 de Mayo",
			"number":           "1234",
			"floor":            "1",
			"apartment":        "A",
			"zip_code":         "3000",
			"city":             "Santa Fe",
			"state":            "Santa Fe",
			"country":          "Argentina",
			"latitude":         "60.0000121",
			"longitude":        "-19.23423",
			"updated_at":       "2023-02-22 11:55:14",
			"observer_user_id": 10
		}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "validate error updating address",
			mock: func() *mock_usecase.MockAddressUseCase {
				mockAddressUseCase := mock_usecase.NewMockAddressUseCase(m)
				return mockAddressUseCase
			},
			path: "/where/are/they/addresses",
			body: `{
			"id":               1,
			"name":             "Casa",
			"zip_code":         "3000",
			"city":             "Santa Fe",
			"state":            "Santa Fe",
			"country":          "Argentina",
			"latitude":         "60.0000121",
			"longitude":        "-19.23423",
			"updated_at":       "2023-02-22 11:55:14",
			"observer_user_id": 10
		}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "error updating address",
			mock: func() *mock_usecase.MockAddressUseCase {
				mockAddressUseCase := mock_usecase.NewMockAddressUseCase(m)
				mockAddressUseCase.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, web.ErrInternalServerError)
				return mockAddressUseCase
			},
			path: "/where/are/they/addresses",
			body: `{
			"id":               1,
			"name":             "Casa",
			"street":           "25 de Mayo",
			"number":           "1234",
			"floor":            "1",
			"apartment":        "A",
			"zip_code":         "3000",
			"city":             "Santa Fe",
			"state":            "Santa Fe",
			"country":          "Argentina",
			"latitude":         "60.0000121",
			"longitude":        "-19.23423",
			"updated_at":       "2023-02-22 11:55:14",
			"observer_user_id": 10
		}`,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "successful update address",
			mock: func() *mock_usecase.MockAddressUseCase {
				mockAddressUseCase := mock_usecase.NewMockAddressUseCase(m)
				mockAddressUseCase.EXPECT().Update(gomock.Any(), gomock.Any()).Return(&a, nil)
				return mockAddressUseCase
			},
			path: "/where/are/they/addresses",
			body: `{
				"id":               1,
				"name":             "Casa",
				"street":           "25 de Mayo",
				"number":           "1234",
				"floor":            "1",
				"apartment":        "A",
				"zip_code":         "3000",
				"city":             "Santa Fe",
				"state":            "Santa Fe",
				"country":          "Argentina",
				"latitude":         "60.0000121",
				"longitude":        "-19.23423",
				"updated_at":       "2023-02-22 11:55:14",
				"observer_user_id": 10
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

			d := mock_middleware.Dependencies{AddressUseCase: test.mock()}
			router := configureRoutes(d)
			router.ServeHTTP(w, r)
			assert.Equalf(t, test.expectedCode, w.Code, "Expected code %v, received %v", test.expectedCode, w.Code)
		})
	}
}

func TestDeleteAddress(t *testing.T) {
	m := gomock.NewController(t)
	defer m.Finish()

	tests := []struct {
		name         string
		mock         func() *mock_usecase.MockAddressUseCase
		path         string
		expectedCode int
	}{
		{
			name: "bad request error deleting address",
			mock: func() *mock_usecase.MockAddressUseCase {
				mockAddressUseCase := mock_usecase.NewMockAddressUseCase(m)
				return mockAddressUseCase
			},
			path:         "/where/are/they/addresses/invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "error deleting address",
			mock: func() *mock_usecase.MockAddressUseCase {
				mockAddressUseCase := mock_usecase.NewMockAddressUseCase(m)
				mockAddressUseCase.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(web.ErrInternalServerError)
				return mockAddressUseCase
			},
			path:         "/where/are/they/addresses/1",
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "not found error deleting address",
			mock: func() *mock_usecase.MockAddressUseCase {
				mockAddressUseCase := mock_usecase.NewMockAddressUseCase(m)
				mockAddressUseCase.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(web.ErrNotFound)
				return mockAddressUseCase
			},
			path:         "/where/are/they/addresses/1",
			expectedCode: http.StatusNotFound,
		},
		{
			name: "successful delete address",
			mock: func() *mock_usecase.MockAddressUseCase {
				mockAddressUseCase := mock_usecase.NewMockAddressUseCase(m)
				mockAddressUseCase.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
				return mockAddressUseCase
			},
			path:         "/where/are/they/addresses/1",
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

			d := mock_middleware.Dependencies{AddressUseCase: test.mock()}
			router := configureRoutes(d)
			router.ServeHTTP(w, r)
			assert.Equalf(t, test.expectedCode, w.Code, "Expected code %v, received %v", test.expectedCode, w.Code)
		})
	}
}
