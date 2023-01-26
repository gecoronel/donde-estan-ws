package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"

	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model/web"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/usecase"
	"github.com/gecoronel/donde-estan-ws/internal/infrastructure/delivery/api/context"
	"github.com/gecoronel/donde-estan-ws/internal/infrastructure/delivery/api/utils"
	log "github.com/sirupsen/logrus"
)

func Get(w http.ResponseWriter, r *http.Request) {
	serviceLocator := context.GetServiceLocator(r.Context())
	useCase := serviceLocator.GetInstance(usecase.UserUseCaseType).(usecase.UserUseCase)

	userID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 0)
	if err != nil {
		log.Error("get user failure in get param: ", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	user, err := useCase.Get(userID, serviceLocator)
	if err != nil {
		log.Error("get user failure: ", err)
		w.Header().Set("Content-Type", "application/json")
		httpStatusCode := utils.GetHTTPCodeByError(err)
		w.WriteHeader(httpStatusCode)
		json.NewEncoder(w).Encode(web.NewError(httpStatusCode, err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func Login(w http.ResponseWriter, r *http.Request) {
	serviceLocator := context.GetServiceLocator(r.Context())
	useCase := serviceLocator.GetInstance(usecase.UserUseCaseType).(usecase.UserUseCase)

	var login model.Login

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&login); err != nil {
		log.Error("body is wrong. ")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	if err := login.Validate(); err != nil {
		log.Error("invalid body")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	user, err := useCase.Login(login, serviceLocator)
	if err != nil {
		log.Error("login failure: ", err)
		w.Header().Set("Content-Type", "application/json")
		httpStatusCode := utils.GetHTTPCodeByError(err)
		w.WriteHeader(httpStatusCode)
		json.NewEncoder(w).Encode(web.NewError(httpStatusCode, err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
