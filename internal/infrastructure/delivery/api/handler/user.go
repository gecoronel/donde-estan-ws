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

func CreateObservedUser(w http.ResponseWriter, r *http.Request) {
	serviceLocator := context.GetServiceLocator(r.Context())
	useCase := serviceLocator.GetInstance(usecase.UserUseCaseType).(usecase.UserUseCase)

	var (
		u            *model.User
		observedUser model.ObservedUser
		err          error
	)

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err = d.Decode(&observedUser); err != nil {
		log.Error("invalid body for creation of observed user")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	if err = observedUser.Validate(); err != nil {
		log.Error("validation failed for creation of observed user")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	u, _ = useCase.FindByUsername(observedUser.GetUsername(), serviceLocator)
	if u != nil {
		log.Error("creation observed user failed: ", web.ErrUsernameConflict)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(web.NewError(http.StatusConflict, web.ErrUsernameConflict.Error()))
		return
	}

	u, _ = useCase.FindByEmail(observedUser.GetEmail(), serviceLocator)
	if u != nil {
		log.Error("creation observed user failed: ", web.ErrEmailConflict)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(web.NewError(http.StatusConflict, web.ErrEmailConflict.Error()))
		return
	}

	user, err := useCase.CreateObservedUser(observedUser, serviceLocator)
	if err != nil {
		log.Error("creation observed user failed: ", err)
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

func CreateObserverUser(w http.ResponseWriter, r *http.Request) {
	serviceLocator := context.GetServiceLocator(r.Context())
	useCase := serviceLocator.GetInstance(usecase.UserUseCaseType).(usecase.UserUseCase)

	var (
		u            *model.User
		observerUser model.ObserverUser
		user         *model.ObserverUser
		err          error
	)

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err = d.Decode(&observerUser); err != nil {
		log.Error("invalid body for creation of observer user")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	if err = observerUser.Validate(); err != nil {
		log.Error("validation failed for creation of observer user")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	u, _ = useCase.FindByUsername(observerUser.GetUsername(), serviceLocator)
	if u != nil {
		log.Error("creation observed user failed: ", web.ErrUsernameConflict)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(web.NewError(http.StatusConflict, web.ErrUsernameConflict.Error()))
		return
	}

	u, _ = useCase.FindByEmail(observerUser.GetEmail(), serviceLocator)
	if u != nil {
		log.Error("creation observed user failed: ", web.ErrEmailConflict)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(web.NewError(http.StatusConflict, web.ErrEmailConflict.Error()))
		return
	}

	user, err = useCase.CreateObserverUser(observerUser, serviceLocator)
	if err != nil {
		log.Error("creation observer user failed ", err)
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
