package tv

import (
	"github.com/gorilla/mux"
	"net/http"
)

type AppRequest struct {
	AppId string
}

func (s *AppRequest) GetApp() App {
	return NewApp(s.AppId)
}

func NewAppRequest(r *http.Request) *AppRequest {
	vars := mux.Vars(r)
	appId, ok := vars["appId"]
	if !ok {
		appId = "camera360"
	}
	newRequest := &AppRequest{
		AppId: appId,
	}
	return newRequest
}
