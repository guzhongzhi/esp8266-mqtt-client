package tv

import (
	"camera360.com/tv/pkg/remotecontrol"
	"camera360.com/tv/pkg/tools"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"strings"
	"sync"
	"time"
)

var apps = make(map[string]*app)
var appLocker sync.Mutex

type App interface {
	GetName() string
	GetPublicTopic() string
	GetUsers() []*DevicePO
	SendMessage(*Command) mqtt.Token
	SendMessageToTopic(topic string, command *Command) mqtt.Token
	GetUserByMac(mac string) *DevicePO
	SendMessageToUser(mac string, command *Command) (mqtt.Token, error)
}

func Apps() map[string]*app {
	return apps
}

func NewApp(clientId string, opts ...AppOption) (*app) {
	options := NewAppOptions(opts...)
	if options.client == nil {
		options.client = client
		log.Println("init application failure,there is no mqtt client,use default mqtt client")
	}
	log.Println("options.Name ", options.Name)
	if options.Name == "" {
		appName := strings.Split(clientId, "-")[0]
		log.Println("new app:", appName)
		options.Name = appName
	}
	if v, ok := apps[options.Name]; ok {
		log.Println("app existing", options.Name)
		return v
	}
	appLocker.Lock()
	defer appLocker.Unlock()
	if v, ok := apps[options.Name]; ok {
		log.Println("app existing logic 2", options.Name)
		return v
	}

	newApp := &app{
		options: options,
		Users:   make(map[string]*DevicePO),
	}
	newApp.init()
	apps[options.Name] = newApp
	return newApp
}

type app struct {
	Users   map[string]*DevicePO
	locker  sync.Mutex
	options *AppOptions
}

func (s *app) Options() *AppOptions {
	return s.options
}

//发送消息到整个app
func (s *app) GetPublicTopic() string {
	return "/" + s.options.Name + "/public-topic"
}

func (s *app) GetUserByMac(mac string) *DevicePO {
	if u, ok := s.Users[mac]; ok {
		return u
	}
	return nil
}

//客户端心跳上报
func (s *app) GetUserHeartBeatTopic() string {
	return "/" + s.options.Name + "/heart-beat"
}

//客户端接收的红外消息上报
func (s *app) GetIRReceivedTopic() string {
	return "/" + s.options.Name + "/ir-received"
}

func (s *app) SendMessage(message *Command) mqtt.Token {
	return s.SendMessageToTopic(s.GetPublicTopic(), message)
}

func (s *app) SendMessageToTopic(topic string, message *Command) mqtt.Token {
	j, _ := json.Marshal(message)
	data := string(j);
	log.Println("publish message to:", topic, data)
	return s.options.client.Publish(topic, s.options.Qos, false, data)
}

//客户端接收消息的topic
func (s *app) GetUserTopic(u *DevicePO) string {
	return "/" + s.options.Name + "/user/" + u.GetTopic()
}

//发送消息给指定客户端
func (s *app) SendMessageToUser(mac string, message *Command) (mqtt.Token, error) {
	user, ok := s.Users[mac]
	if !ok {
		return nil, errors.New("invalid user mac address")
	}
	topic := s.GetUserTopic(user)
	token := s.SendMessageToTopic(topic, message)
	token.Wait()
	fmt.Println("user.ExecutedAt", user.ExecutedAt, token.Error())
	return token, token.Error()
}

func (s *app) OnIRReceived(client mqtt.Client, message mqtt.Message) {
	body := string(message.Payload())
	request := &HeartBeatRequest{}
	fmt.Println(body, request.Data)
	if request.Data == "" {
		log.Println("ir data is empty")
		return
	}
	v := `{label:"%s",value:"%s"},`
	fmt.Println(fmt.Sprintf(v, tools.RandStringBytes(10), request.Data))

	btn, _ := remotecontrol.NewButton(context.Background())
	po := btn.GetPO()
	po.AppName = s.options.Name
	po.Name = tools.RandStringBytes(10)
	po.IrCode = request.Data
	po.CreatedAt = time.Now().Unix()
	po.UpdatedAt = time.Now().Unix()
	btn.Save()
}

func (s *app) OnHeartBeat(client mqtt.Client, request *HeartBeatRequest) {
	now := time.Now().Unix()
	if request.Mac == "" {
		return
	}
	if user, ok := s.Users[request.Mac]; ok {
		fmt.Println("user existing: ", request.Mac)
		user.Relay = request.Relay
		user.HeartbeatAt = now
		user.IP = request.IP
		user.WIFI = request.WIFI
		user.IsNewBoot = request.IsNewBoot
		user.RelayPin = request.RelayPIN
		if request.ExecutedAt > 0 {
			user.ExecutedAt = request.ExecutedAt
		}
		fmt.Println("user.Relay", user.Relay)
		s.saveUser(user)
	} else {
		fmt.Println("no user: ", request.Mac)
		user := &DevicePO{
			AppName:     s.options.Name,
			ModeId:      []string{},
			Mac:         request.Mac,
			WIFI:        request.WIFI,
			IP:          request.IP,
			Name:        request.ClientId,
			Relay:       request.Relay,
			RelayPin:    request.RelayPIN,
			ExecutedAt:  request.ExecutedAt,
			ConnectedAt: now,
			HeartbeatAt: now,
			IsNewBoot:   request.IsNewBoot,
		}
		s.saveUser(user)
		s.AddUser(user)
	}
	s.sendUsersToWS()
}

func (s *app) saveUser(user *DevicePO) error {
	device, _ := NewDevice(context.Background())
	device.LoadByMac(user.Mac)

	//设备初始化函数
	onBootFunc := func(device *DevicePO) error {
		time.Sleep(time.Second)
		//开机时初始化继电器状态
		if s.options.DeviceRelayStatusOnBoot == "on" {
			//开机时打开继电器
			if device.RelayTriggeredByLowLevel && device.Relay == "on" {
				s.SendMessageToUser(device.Mac, NewOffCommand())
			} else if !device.RelayTriggeredByLowLevel && device.Relay == "off" {
				s.SendMessageToUser(device.Mac, NewOnCommand())
			}
		} else {
			//开机时关闭继电器
			if device.RelayTriggeredByLowLevel && device.Relay == "off" {
				s.SendMessageToUser(device.Mac, NewOnCommand())
			} else if !device.RelayTriggeredByLowLevel && device.Relay == "on" {
				s.SendMessageToUser(device.Mac, NewOffCommand())
			}
		}

		//开机其它自定义命令
		if device.OnBootCommand == "" {
			return nil
		}
		tmp := strings.Split(strings.TrimSpace(device.OnBootCommand), ",")
		for _, cmd := range tmp {
			s.SendMessageToUser(device.Mac, NewCmd(cmd, ""))
			time.Sleep(time.Second * 3)
		}

		return nil
	}

	if device.HasId() {
		device.GetPlainObject().WIFI = user.WIFI
		device.GetPlainObject().Relay = user.Relay
		device.GetPlainObject().RelayPin = user.RelayPin
		device.GetPlainObject().IP = user.IP
		device.GetPlainObject().HeartbeatAt = user.HeartbeatAt
		device.GetPlainObject().RelayPin = user.RelayPin
		user.Id = device.GetPlainObject().Id

		//开机的时候初始化继电器引脚，
		if device.GetPlainObject().HasCustomRelayPin && device.GetPlainObject().CustomRelayPin != user.RelayPin {
			s.SendMessageToUser(device.GetPlainObject().Mac, NewCmd("setRelayPIN", device.GetPlainObject().CustomRelayPin))
		}
		if user.IsNewBoot {
			//开机的时候重新初始化状态
			onBootFunc(device.GetPlainObject())
		}
	} else {
		device.SetData(user)
	}
	device.Save()
	return nil
}

func (s *app) sendUsersToWS() error {
	device, _ := NewDevice(context.Background())
	for _, user := range s.Users {
		device.Load(user.Id)
		if device.HasId() {
			user.Name = device.GetPlainObject().Name
		}
	}
	msg := WebSocketClientMessage{
		Operation: "users",
		Data:      s.Users,
	}
	js, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	for c, _ := range hub.clients {
		if c.appName == s.options.Name {
			c.send <- js
		}
	}
	return nil
}

func (s *app) init() {
	var err error
	var users map[string]*DevicePO
	client := s.options.client
	log.Println("subscribe to public ir received:", s.GetIRReceivedTopic())
	log.Println("app boardcast topic:", s.GetPublicTopic())
	fmt.Println("s.GetIRReceivedTopic()", s.GetIRReceivedTopic())
	client.Subscribe(s.GetIRReceivedTopic(), s.options.Qos, s.OnIRReceived)
	users, err = loadUsers(s.options.Name)
	if err != nil {
		log.Println("err", err)
	} else {
		s.Users = users
	}
}

func (s *app) AddUser(user *DevicePO) App {
	s.locker.Lock()
	defer s.locker.Unlock()
	s.Users[user.Mac] = user
	log.Println("user topic:", s.GetUserTopic(user))
	return s
}

func (s *app) GetUsers() []*DevicePO {
	users := make([]*DevicePO, 0)
	for _, v := range s.Users {
		users = append(users, v)
	}
	return users
}

func (s *app) GetName() string {
	return s.options.Name
}
