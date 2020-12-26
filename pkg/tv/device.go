package tv

import (
	"code.aliyun.com/MIG-server/micro-base/orm/mongo"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RelayPin *int

type DevicePO struct {
	Id                       primitive.ObjectID `json:"id" bson:"_id"`
	AppName                  string             `json:"appName" bson:"appName"`
	Name                     string             `json:"name" bson:"name"` //设备名
	IP                       string             `json:"ip" bson:"ip"`
	WIFI                     string             `json:"wifi" bson:"wifi"`
	Relay                    string             `json:"relay" bson:"relay"`
	RelayPin				 int     			`json:"RelayPin" bson:"RelayPin"`
	CustomRelayPin 			 int 				`json:"CustomRelayPin" bson:"CustomRelayPin"`
	HasCustomRelayPin		 bool 				`json:"HasCustomRelayPin" bson:"HasCustomRelayPin"`
	RelayTriggeredByLowLevel bool               `json:"relayTriggeredByLowLevel" bson:"relayTriggeredByLowLevel"`
	Mac                      string             `json:"mac" bson:"mac"`
	ModeId                   []string           `json:"modeId" bson:"modeId"` //遥控板
	ConnectedAt              int64              `json:"connectedAt" bson:"connectedAt"`
	HeartbeatAt              int64              `json:"heartbeatAt" bson:"heartbeatAt"`
	ExecutedAt               int64              `json:"executedAt" bson:"executedAt"` //最后执行的指令
}

func (s *DevicePO) ValidateMessages() map[string]string {
	return nil
}

func (s *DevicePO) GetTopic() string {
	return s.Mac
}

func NewDevice(ctx context.Context) (*Device, error) {
	d := &Device{}
	err := d.initialize(ctx)
	if err != nil {
		return nil, err
	}
	return d, nil
}

type Device struct {
	mongo.Object
}

func (s *Device) LoadByMac(mac string) error {
	return s.LoadByCondition(mongo.M{
		"mac": mac,
	})
}

func (s *Device) initialize(ctx context.Context) (error) {
	return s.InitializeWithContext(ctx, &DevicePO{}, "Id", "device", mongo.ConfigNodeNameOption("tvads"))
}

func (s *Device) GetPlainObject() *DevicePO {
	return s.Data.(*DevicePO)
}

func loadUsers(appName string) map[string]*DevicePO {
	device, _ := NewDevice(context.Background())
	collection := device.GetCollection()
	collection.Where(mongo.M{
		"appName": appName,
	})
	data := make(map[string]*DevicePO)
	return data
}

//@TODO
func saveUser(user *DevicePO) error {
	device, _ := NewDevice(context.Background())
	device.LoadByMac(user.Mac)

	if device.HasId() {
		device.GetPlainObject().WIFI = user.WIFI
		device.GetPlainObject().Relay = user.Relay
		device.GetPlainObject().IP = user.IP
		device.GetPlainObject().HeartbeatAt = user.HeartbeatAt
		device.GetPlainObject().RelayPin = user.RelayPin
		if device.GetPlainObject().HasCustomRelayPin && device.GetPlainObject().RelayPin != device.GetPlainObject().CustomRelayPin {
			NewApp(user.Name,NewAppNameOption(user.AppName)).SendMessageToUser(user.Mac,NewCmd("setRelayPIN",user.HasCustomRelayPin))
			//{"cmd":"setRelayPIN","executedAt":19939838,data:0}
		}
	} else {
		device.SetData(user)
	}
	device.Save()
	return nil
}
