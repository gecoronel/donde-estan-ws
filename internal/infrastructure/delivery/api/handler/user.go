package handler

import (
	"encoding/json"
	"github.com/gcoron/donde-estan-ws/internal/bussiness/model/web"
	"github.com/gcoron/donde-estan-ws/internal/infrastructure/delivery/api/context"
	"github.com/gcoron/donde-estan-ws/internal/infrastructure/delivery/api/utils"
	"io"
	"net/http"

	"github.com/gcoron/donde-estan-ws/internal/bussiness/model"
	"github.com/gcoron/donde-estan-ws/internal/bussiness/usecase"
	log "github.com/sirupsen/logrus"
)

func Login(w http.ResponseWriter, r *http.Request) { //c *gin.Context) {
	serviceLocator := context.GetServiceLocator(r.Context())
	useCase := serviceLocator.GetInstance(usecase.LoginUseCaseType).(usecase.LoginUseCase)

	var bodyBytes []byte
	var login model.Login

	if r.Body != nil {
		bodyBytes, _ = io.ReadAll(r.Body)
	}

	if len(bodyBytes) <= 0 {
		log.Error("body is nil. ")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, "body is empty"))
		return
	}

	if err := json.Unmarshal(bodyBytes, &login); err != nil {
		log.Error("post return Login unmarshall returns an error. ", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(web.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	user, err := useCase.Login(login, serviceLocator)

	if err != nil {
		log.Error("login failure. ", err)
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
