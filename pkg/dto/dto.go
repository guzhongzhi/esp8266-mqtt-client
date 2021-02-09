package dto

import (
	"github.com/gorilla/mux"
	"net/http"
)

type AppRequest struct {
	AppId string
}

func NewAppRequest(r *http.Request) *AppRequest {
	vars := mux.Vars(r)
	appId, ok := vars["appId"]
	if !ok {
		appId = "guz"
	}
	newRequest := &AppRequest{
		AppId: appId,
	}
	return newRequest
}

type BeatRequest struct {
	Mac        string `json:"m"`
	IP         string `json:"i"`
	WIFI       string `json:"w"`
	ClientId   string `json:"cid"`
	Gateway    string `json:"g"`
	Relay      string `json:"r"`
	RelayPIN   int    `json:"rp"`
	StatePIN   int    `json:"sp"`
	IrPIN      int    `json:"irp"`
	App        string `json:"a"`
	Data       string `json:"d"`
	ExecutedAt int64  `json:"e"`
	IsNewBoot  bool   `json:"b"`
	Version    string `json:"v"`
	Command    string `json:"c"`
}
