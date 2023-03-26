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

func GetChild(w http.ResponseWriter, r *http.Request) {
	serviceLocator := context.GetServiceLocator(r.Context())
	useCase := serviceLocator.GetInstance(usecase.ChildUseCaseType).(usecase.ChildUseCase)

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 0)
	if err != nil {
		log.Error("get child failure in get param: ", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	sb, err := useCase.Get(id, serviceLocator)
	if err != nil {
		log.Error("get child failure: ", err)
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

func SaveChild(w http.ResponseWriter, r *http.Request) {
	serviceLocator := context.GetServiceLocator(r.Context())
	useCase := serviceLocator.GetInstance(usecase.ChildUseCaseType).(usecase.ChildUseCase)

	var (
		child model.Child
		err   error
	)

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err = d.Decode(&child); err != nil {
		log.Error("invalid body for creation of child")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	if err = child.Validate(); err != nil || child.ValidateID() || !child.ValidateObserverUserID() {
		log.Error("validation failed for creation of child")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(
			web.NewError(http.StatusBadRequest, "validation failed for creation of child"),
		)
		return
	}

	sb, err := useCase.Save(child, serviceLocator)
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

func UpdateChild(w http.ResponseWriter, r *http.Request) {
	serviceLocator := context.GetServiceLocator(r.Context())
	useCase := serviceLocator.GetInstance(usecase.ChildUseCaseType).(usecase.ChildUseCase)

	var (
		child model.Child
		err   error
	)

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err = d.Decode(&child); err != nil {
		log.Error("invalid body for creation of child")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	if err = child.Validate(); err != nil || !child.ValidateID() {
		log.Error("validation failed for creation of child")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(
			web.NewError(http.StatusBadRequest, "validation failed for creation of child"),
		)
		return
	}

	sb, err := useCase.Update(child, serviceLocator)
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

func DeleteChild(w http.ResponseWriter, r *http.Request) {
	serviceLocator := context.GetServiceLocator(r.Context())
	useCase := serviceLocator.GetInstance(usecase.ChildUseCaseType).(usecase.ChildUseCase)

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 0)
	if err != nil {
		log.Error("get child failure in get param: ", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	err = useCase.Delete(id, serviceLocator)
	if err != nil {
		log.Error("get child failure: ", err)
		w.Header().Set("Content-Type", "application/json")
		httpStatusCode := utils.GetHTTPCodeByError(err)
		w.WriteHeader(httpStatusCode)
		_ = json.NewEncoder(w).Encode(web.NewError(httpStatusCode, err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(web.NewError(http.StatusOK, "child deleted"))
}
