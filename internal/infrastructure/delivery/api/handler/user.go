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
		if err = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error())); err != nil {
			log.Error("error encoding error data: ", err)
		}
		return
	}

	user, err := useCase.Get(userID, serviceLocator)
	if err != nil {
		log.Error("get user failure: ", err)
		w.Header().Set("Content-Type", "application/json")
		httpStatusCode := utils.GetHTTPCodeByError(err)
		w.WriteHeader(httpStatusCode)
		if err = json.NewEncoder(w).Encode(web.NewError(httpStatusCode, err.Error())); err != nil {
			log.Error("error encoding error data: ", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(user); err != nil {
		log.Error("error encoding user data: ", err)
		return
	}
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
		if err = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error())); err != nil {
			log.Error("error encoding error data: ", err)
		}
		return
	}

	if err := login.Validate(); err != nil {
		log.Error("invalid body")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error())); err != nil {
			log.Error("error encoding error data: ", err)
		}
		return
	}

	user, err := useCase.Login(login, serviceLocator)
	if err != nil {
		log.Error("login failure: ", err)
		w.Header().Set("Content-Type", "application/json")
		httpStatusCode := utils.GetHTTPCodeByError(err)
		w.WriteHeader(httpStatusCode)
		if err = json.NewEncoder(w).Encode(web.NewError(httpStatusCode, err.Error())); err != nil {
			log.Error("error encoding error data: ", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(user); err != nil {
		log.Error("error encoding user data: ", err)
		return
	}
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
		if err = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error())); err != nil {
			log.Error("error encoding error data: ", err)
		}
		return
	}

	if err = observedUser.Validate(); err != nil {
		log.Error("validation failed for creation of observed user")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error())); err != nil {
			log.Error("error encoding error data: ", err)
		}
		return
	}

	u, _ = useCase.FindByUsername(observedUser.GetUsername(), serviceLocator)
	if u != nil {
		log.Error("creation observed user failed: ", web.ErrConflict)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		if err = json.NewEncoder(w).Encode(web.NewError(http.StatusConflict, web.ErrConflict.Error())); err != nil {
			log.Error("error encoding error data: ", err)
		}
		return
	}

	u, _ = useCase.FindByEmail(observedUser.GetEmail(), serviceLocator)
	if u != nil {
		log.Error("creation observed user failed: ", web.ErrConflict)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		if err = json.NewEncoder(w).Encode(web.NewError(http.StatusConflict, web.ErrConflict.Error())); err != nil {
			log.Error("error encoding error data: ", err)
		}
		return
	}

	user, err := useCase.CreateObservedUser(observedUser, serviceLocator)
	if err != nil {
		log.Error("creation observed user failed: ", err)
		w.Header().Set("Content-Type", "application/json")
		httpStatusCode := utils.GetHTTPCodeByError(err)
		w.WriteHeader(httpStatusCode)
		if err = json.NewEncoder(w).Encode(web.NewError(httpStatusCode, err.Error())); err != nil {
			log.Error("error encoding error data: ", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(user); err != nil {
		log.Error("error encoding user data: ", err)
	}
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
		if err = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error())); err != nil {
			log.Error("error encoding error data: ", err)
		}
		return
	}

	if err = observerUser.Validate(); err != nil {
		log.Error("validation failed for creation of observer user")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error())); err != nil {
			log.Error("error encoding error data: ", err)
		}
		return
	}

	u, _ = useCase.FindByUsername(observerUser.GetUsername(), serviceLocator)
	if u != nil {
		log.Error("creation observed user failed: ", web.ErrConflict)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		if err = json.NewEncoder(w).Encode(web.NewError(http.StatusConflict, web.ErrConflict.Error())); err != nil {
			log.Error("error encoding error data: ", err)
		}
		return
	}

	u, _ = useCase.FindByEmail(observerUser.GetEmail(), serviceLocator)
	if u != nil {
		log.Error("creation observed user failed: ", web.ErrConflict)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		if err = json.NewEncoder(w).Encode(web.NewError(http.StatusConflict, web.ErrConflict.Error())); err != nil {
			log.Error("error encoding error data: ", err)
		}
		return
	}

	user, err = useCase.CreateObserverUser(observerUser, serviceLocator)
	if err != nil {
		log.Error("creation observer user failed ", err)
		w.Header().Set("Content-Type", "application/json")
		httpStatusCode := utils.GetHTTPCodeByError(err)
		w.WriteHeader(httpStatusCode)
		if err = json.NewEncoder(w).Encode(web.NewError(httpStatusCode, err.Error())); err != nil {
			log.Error("error encoding error data: ", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(user); err != nil {
		log.Error("error encoding user data: ", err)
	}
}
