package tv

import (
	"camera360.com/tv/pkg/controller"
	"code.aliyun.com/MIG-server/micro-base/orm/mongo"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"sort"
)

func NewApi() *Api {
	a := &Api{}
	return a
}

type Api struct {
	controller.Controller
	AppName string
}

func (s *Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	s.AppName = vars["appId"]
	s.SetRequest(r)
	s.SetResponse(w)
	s.Dispatch(s)
}

func (s *Api) Dashboard() error {
	request := s.NewAppRequest()
	tmp := fmt.Sprintf(indexHTML, request.AppId)
	s.Response().Write([]byte(tmp))
	return nil
}

func (s *Api) Users() error {
	request := s.NewAppRequest()
	users := request.GetApp().GetUsers()
	sort.SliceStable(users, func(i, j int) bool {
		return users[i].Mac < users[j].Mac
	})
	return s.WriteJSON(users)
}

func (s *Api) SendIr() error {
	code := s.String("code")
	s.NewAppRequest().GetApp().SendMessage(NewIrSendCommand(code))
	return nil
}

func (s *Api) DeviceSendIr() error {
	r := s.NewAppRequest()
	app := r.GetApp()
	code := s.String("code")
	mac := s.String("mac")

	if mac != "" {
		app.SendMessageToUser(mac, NewIrSendCommand(code))
	} else {
		fmt.Println("send message to application user failure, user do not existing")
	}
	return nil
}

func (s *Api) DeviceSendMessage() error {
	appRequest := s.NewAppRequest()
	app := appRequest.GetApp()
	cmd := s.String("cmd")
	mac := s.String("mac")
	if mac != "" {
		app.SendMessageToUser(mac, NewCmd(cmd, ""))
	} else {
		fmt.Println("invalid user mac:", mac)
	}
	s.WriteJSON("OK")
	return nil
}

func (s *Api) SendMessage() error {
	appRequest := s.NewAppRequest()
	app := appRequest.GetApp()
	topic := s.String("topic")
	cmd := s.String("cmd")
	data := s.String("data")
	if topic != "" {
		app.SendMessageToTopic(topic, NewCmd(cmd, data))
	} else {
		app.SendMessage(NewCmd(cmd, data))
	}
	return s.WriteJSON("OK")
}

func (s *Api) DeviceList() error {
	request := s.NewAppRequest()
	device, err := NewDevice(s.Request().Context())
	if err != nil {
		return err
	}
	pager, _ := device.GetCollection().Where(mongo.M{
		"appName": request.AppId,
	}).GetPager(1, 1000)
	s.WriteJSON(pager)
	return nil
}

func (s *Api) DeviceDelete() error {
	id := s.String("id")
	d, _ := NewDevice(s.Request().Context())
	d.Load(id)
	if d.HasId() {
		d.Delete()
	}
	s.WriteStatusData(true, http.StatusOK, "OK")
	return nil
}

func (s *Api) DeviceSave() error {

	devicePO := &DevicePO{}
	body, _ := ioutil.ReadAll(s.Request().Body)
	json.Unmarshal(body, devicePO)

	d, _ := NewDevice(context.Background())
	d.LoadByMac(devicePO.Mac)
	if d.HasId() == false {
		return errors.New("invalid device")
	}
	d.GetPlainObject().Name = devicePO.Name
	d.GetPlainObject().ModeId = devicePO.ModeId
	d.GetPlainObject().RelayTriggeredByLowLevel = devicePO.RelayTriggeredByLowLevel
	isRelayPinChanged := d.GetPlainObject().HasCustomRelayPin != devicePO.HasCustomRelayPin

	fmt.Println("devicePO.HasCustomRelayPin",devicePO.HasCustomRelayPin,devicePO.CustomRelayPin)
	d.GetPlainObject().CustomRelayPin = devicePO.CustomRelayPin
	d.GetPlainObject().HasCustomRelayPin = devicePO.HasCustomRelayPin
	d.Save()
	app := NewApp(devicePO.Name,NewAppNameOption(devicePO.AppName))
	app.AddUser( d.GetPlainObject())

	if isRelayPinChanged {
		if d.GetPlainObject().HasCustomRelayPin {
			app.SendMessageToUser(d.GetPlainObject().Mac,NewCmd("setRelayPIN",d.GetPlainObject().CustomRelayPin))
		} else {
			//use default relay pin 5
			app.SendMessageToUser(d.GetPlainObject().Mac,NewCmd("setRelayPIN",5))
		}
	}

	s.WriteJSON("OK")
	return nil
}

func (s *Api) NewAppRequest() *AppRequest {
	return NewAppRequest(s.Request())
}
