package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model/web"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/usecase"
	"github.com/gecoronel/donde-estan-ws/internal/infrastructure/delivery/api/context"
	"github.com/gecoronel/donde-estan-ws/internal/infrastructure/delivery/api/utils"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

func GetAddress(w http.ResponseWriter, r *http.Request) {
	serviceLocator := context.GetServiceLocator(r.Context())
	useCase := serviceLocator.GetInstance(usecase.AddressUseCaseType).(usecase.AddressUseCase)

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 0)
	if err != nil {
		log.Error("get address failure in get param: ", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	sb, err := useCase.Get(id, serviceLocator)
	if err != nil {
		log.Error("get address failure: ", err)
		w.Header().Set("Content-Type", "application/json")
		httpStatusCode := utils.GetHTTPCodeByError(err)
		w.WriteHeader(httpStatusCode)
		_ = json.NewEncoder(w).Encode(web.NewError(httpStatusCode, err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(*sb)
}

func SaveAddress(w http.ResponseWriter, r *http.Request) {
	serviceLocator := context.GetServiceLocator(r.Context())
	useCase := serviceLocator.GetInstance(usecase.AddressUseCaseType).(usecase.AddressUseCase)

	var (
		address model.Address
		err     error
	)

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err = d.Decode(&address); err != nil {
		log.Error("invalid body for creation of address")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	if err = address.Validate(); err != nil || address.ValidateID() {
		log.Error("validation failed for creation of address")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(
			web.NewError(http.StatusBadRequest, "validation failed for creation of address"),
		)
		return
	}

	sb, err := useCase.Save(address, serviceLocator)
	if err != nil {
		log.Error("creation observed user failed: ", err)
		w.Header().Set("Content-Type", "application/json")
		httpStatusCode := utils.GetHTTPCodeByError(err)
		w.WriteHeader(httpStatusCode)
		_ = json.NewEncoder(w).Encode(web.NewError(httpStatusCode, err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(sb)
}

func UpdateAddress(w http.ResponseWriter, r *http.Request) {
	serviceLocator := context.GetServiceLocator(r.Context())
	useCase := serviceLocator.GetInstance(usecase.AddressUseCaseType).(usecase.AddressUseCase)

	var (
		address model.Address
		err     error
	)

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err = d.Decode(&address); err != nil {
		log.Error("invalid body for creation of address")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	if err = address.Validate(); err != nil || !address.ValidateID() {
		log.Error("validation failed for creation of address")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(
			web.NewError(http.StatusBadRequest, "validation failed for creation of address"),
		)
		return
	}

	sb, err := useCase.Update(address, serviceLocator)
	if err != nil {
		log.Error("creation observed user failed: ", err)
		w.Header().Set("Content-Type", "application/json")
		httpStatusCode := utils.GetHTTPCodeByError(err)
		w.WriteHeader(httpStatusCode)
		_ = json.NewEncoder(w).Encode(web.NewError(httpStatusCode, err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(sb)
}

func DeleteAddress(w http.ResponseWriter, r *http.Request) {
	serviceLocator := context.GetServiceLocator(r.Context())
	useCase := serviceLocator.GetInstance(usecase.AddressUseCaseType).(usecase.AddressUseCase)

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 0)
	if err != nil {
		log.Error("get address failure in get param: ", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	err = useCase.Delete(id, serviceLocator)
	if err != nil {
		log.Error("get address failure: ", err)
		w.Header().Set("Content-Type", "application/json")
		httpStatusCode := utils.GetHTTPCodeByError(err)
		w.WriteHeader(httpStatusCode)
		_ = json.NewEncoder(w).Encode(web.NewError(httpStatusCode, err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(web.NewError(http.StatusOK, "address deleted"))
}
