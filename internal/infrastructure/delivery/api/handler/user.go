package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model/web"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/usecase"
	"github.com/gecoronel/donde-estan-ws/internal/infrastructure/delivery/api/context"
	"github.com/gecoronel/donde-estan-ws/internal/infrastructure/delivery/api/utils"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	serviceLocator := context.GetServiceLocator(r.Context())
	useCase := serviceLocator.GetInstance(usecase.UserUseCaseType).(usecase.UserUseCase)

	userID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 0)
	if err != nil {
		log.Error("get user failure in get param: ", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	user, err := useCase.Get(userID, serviceLocator)
	if err != nil {
		log.Error("get user failure: ", err)
		w.Header().Set("Content-Type", "application/json")
		httpStatusCode := utils.GetHTTPCodeByError(err)
		w.WriteHeader(httpStatusCode)
		_ = json.NewEncoder(w).Encode(web.NewError(httpStatusCode, err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
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
		_ = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	if err := login.Validate(); err != nil {
		log.Error("invalid body")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	user, err := useCase.Login(login, serviceLocator)
	if err != nil {
		log.Error("login failure: ", err)
		w.Header().Set("Content-Type", "application/json")
		httpStatusCode := utils.GetHTTPCodeByError(err)
		w.WriteHeader(httpStatusCode)
		_ = json.NewEncoder(w).Encode(web.NewError(httpStatusCode, err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}

func CreateObservedUser(w http.ResponseWriter, r *http.Request) {
	serviceLocator := context.GetServiceLocator(r.Context())
	useCase := serviceLocator.GetInstance(usecase.UserUseCaseType).(usecase.UserUseCase)

	var (
		observedUser model.ObservedUser
		err          error
	)

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err = d.Decode(&observedUser); err != nil {
		log.Error("invalid body for creation of observed user")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	if err = observedUser.Validate(); err != nil {
		log.Error("validation failed for creation of observed user")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	user, err := useCase.CreateObservedUser(observedUser, serviceLocator)
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
	_ = json.NewEncoder(w).Encode(user)
}

func UpdateObservedUser(w http.ResponseWriter, r *http.Request) {
	serviceLocator := context.GetServiceLocator(r.Context())
	useCase := serviceLocator.GetInstance(usecase.UserUseCaseType).(usecase.UserUseCase)

	var (
		observedUser model.ObservedUser
		err          error
	)

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err = d.Decode(&observedUser); err != nil {
		log.Error("invalid body for update of observed user")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	if err = observedUser.Validate(); err != nil || !observedUser.User.ValidateID() {
		log.Error("validation failed for update of observed user")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(
			web.NewError(http.StatusBadRequest, "validation failed for update of observed user"),
		)
		return
	}

	user, err := useCase.UpdateObservedUser(observedUser, serviceLocator)
	if err != nil {
		log.Error("update observed user failed: ", err)
		w.Header().Set("Content-Type", "application/json")
		httpStatusCode := utils.GetHTTPCodeByError(err)
		w.WriteHeader(httpStatusCode)
		_ = json.NewEncoder(w).Encode(web.NewError(httpStatusCode, err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}

func CreateObserverUser(w http.ResponseWriter, r *http.Request) {
	serviceLocator := context.GetServiceLocator(r.Context())
	useCase := serviceLocator.GetInstance(usecase.UserUseCaseType).(usecase.UserUseCase)

	var (
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
		_ = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	if err = observerUser.Validate(); err != nil {
		log.Error("validation failed for creation of observer user")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	user, err = useCase.CreateObserverUser(observerUser, serviceLocator)
	if err != nil {
		log.Error("creation observer user failed ", err)
		w.Header().Set("Content-Type", "application/json")
		httpStatusCode := utils.GetHTTPCodeByError(err)
		w.WriteHeader(httpStatusCode)
		_ = json.NewEncoder(w).Encode(web.NewError(httpStatusCode, err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}

func UpdateObserverUser(w http.ResponseWriter, r *http.Request) {
	serviceLocator := context.GetServiceLocator(r.Context())
	useCase := serviceLocator.GetInstance(usecase.UserUseCaseType).(usecase.UserUseCase)

	var (
		observerUser model.ObserverUser
		err          error
	)

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err = d.Decode(&observerUser); err != nil {
		log.Error("invalid body for update of observer user")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	if err = observerUser.Validate(); err != nil || !observerUser.User.ValidateID() {
		log.Error("validation failed for update of observer user")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(
			web.NewError(http.StatusBadRequest, "validation failed for update of observer user"),
		)
		return
	}

	user, err := useCase.UpdateObserverUser(observerUser, serviceLocator)
	if err != nil {
		log.Error("update observer user failed: ", err)
		w.Header().Set("Content-Type", "application/json")
		httpStatusCode := utils.GetHTTPCodeByError(err)
		w.WriteHeader(httpStatusCode)
		_ = json.NewEncoder(w).Encode(web.NewError(httpStatusCode, err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}

func AddObservedUserInObserverUser(w http.ResponseWriter, r *http.Request) {
	serviceLocator := context.GetServiceLocator(r.Context())
	useCase := serviceLocator.GetInstance(usecase.UserUseCaseType).(usecase.UserUseCase)

	var (
		req model.AddDriverReq
		err error
	)

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err = d.Decode(&req); err != nil {
		log.Error("invalid body for creation of observed user")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	if err = req.Validate(); err != nil {
		log.Error("validation failed for add observed user in observer user")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	err = useCase.AddObservedUserInObserverUser(req.PrivacyKey, req.ObserverUserID, serviceLocator)
	if err != nil {
		log.Error("creation observed user in observer user failed: ", err)
		w.Header().Set("Content-Type", "application/json")
		httpStatusCode := utils.GetHTTPCodeByError(err)
		w.WriteHeader(httpStatusCode)
		_ = json.NewEncoder(w).Encode(web.NewError(httpStatusCode, err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(web.NewError(http.StatusOK, "observed user added"))
}

func DeleteObservedUserInObserverUser(w http.ResponseWriter, r *http.Request) {
	serviceLocator := context.GetServiceLocator(r.Context())
	useCase := serviceLocator.GetInstance(usecase.UserUseCaseType).(usecase.UserUseCase)

	ids := strings.Split(chi.URLParam(r, "id"), "_")
	if len(ids) != 2 {
		log.Error("invalid ids for delete of observed user")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, "invalid ids"))
		return
	}

	observedUserID, err := strconv.ParseUint(ids[0], 10, 0)
	if err != nil {
		log.Error("invalid ids for delete of observed user")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, "invalid ids"))
		return
	}
	observerUserID, err := strconv.ParseUint(ids[1], 10, 0)
	if err != nil {
		log.Error("invalid ids for delete of observed user")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, "invalid ids"))
		return
	}

	err = useCase.DeleteObservedUserInObserverUser(observedUserID, observerUserID, serviceLocator)
	if err != nil {
		log.Error("creation observed user in observer user failed: ", err)
		w.Header().Set("Content-Type", "application/json")
		httpStatusCode := utils.GetHTTPCodeByError(err)
		w.WriteHeader(httpStatusCode)
		_ = json.NewEncoder(w).Encode(web.NewError(httpStatusCode, err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(web.NewError(http.StatusOK, "observed user added"))
}
