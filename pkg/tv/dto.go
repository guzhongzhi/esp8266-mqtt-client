package tv

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
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

type HeartBeatRequest struct {
	Mac         string `json:"mac"`
	IP          string `json:"ip"`
	JsonEnabled bool   `json:"jsonEnabled"`
	WIFI        string `json:"wifi"`
	ClientId    string `json:"clientId"`
	Gateway     string `json:"gateway"`
	Relay       string `json:"relay"`
	RelayPIN    int    `json:"relayPin"`
	StatePIN    int    `json:"statePin"`
	IrPIN       int    `json:"irPin"`
	AppName     string `json:"appName"`
	Data        string `json:"data"`
	ExecutedAt  int64  `json:"executedAt"`
	IsNewBoot 	bool   `json:"isNewBoot"`
}

type Command struct {
	Cmd        string 		`json:"cmd"`
	Data       interface{}  `json:"data"`
	ExecutedAt int64  		`json:"executedAt"`
}

func (s *Command) IsTurnOff() bool {
	return s.Cmd == "off"
}

func (s *Command) ToString() string {
	return fmt.Sprintf("%s,%v,%v", s.Cmd, s.ExecutedAt, s.Data)
}

func NewCmd(cmd string, data interface{}) *Command {
	c := &Command{
		Cmd:        cmd,
		Data:       data,
		ExecutedAt: time.Now().Unix(),
	}
	return c
}
func NewIrSendCommand(data string) *Command {
	return NewCmd("irs", data)
}
func NewOffCommand() *Command {
	return NewCmd("off", "")
}

func NewOnCommand() *Command {
	return NewCmd("on", "")
}
