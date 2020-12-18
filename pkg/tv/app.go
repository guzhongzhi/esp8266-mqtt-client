package tv

import (
	"camera360.com/tv/pkg/remotecontrol"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"net/url"
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
	SendMessage(message interface{}) mqtt.Token
	SendMessageToTopic(topic string, message interface{}) mqtt.Token
	GetUserByMac(mac string) *DevicePO
	SendMessageToUser(mac string, message interface{}) (mqtt.Token, error)
}

func Apps() map[string]*app {
	return apps
}

func NewApp(clientId string, opts ...AppOption) *app {
	appName := strings.Split(clientId, "-")[0]
	if v, ok := apps[appName]; ok {
		log.Println("app existing", appName)
		return v
	}
	log.Println("new app:", appName)
	appLocker.Lock()
	defer appLocker.Unlock()
	if v, ok := apps[appName]; ok {
		return v
	}
	opts = append(opts, NewAppNameOption(appName))
	options := NewAppOptions(opts...)
	if options.client == nil {
		log.Println("init application failure,there is no mqtt client")
		options.client = client
	}

	newApp := &app{
		options: options,
		Users:   make(map[string]*DevicePO),
	}
	newApp.init()
	apps[appName] = newApp
	return newApp
}

type app struct {
	Users   map[string]*DevicePO
	locker  sync.Mutex
	options *AppOptions
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

func (s *app) SendMessage(message interface{}) mqtt.Token {
	return s.SendMessageToTopic(s.GetPublicTopic(), message)
}

func (s *app) SendMessageToTopic(topic string, message interface{}) mqtt.Token {
	log.Println("publish message to:", topic, message)
	return s.options.client.Publish(topic, s.options.Qos, false, message)
}

//客户端接收消息的topic
func (s *app) GetUserTopic(u *DevicePO) string {
	return "/" + s.options.Name + "/user/" + u.GetTopic()
}

//发送消息给指定客户端
func (s *app) SendMessageToUser(mac string, message interface{}) (mqtt.Token, error) {
	user, ok := s.Users[mac]
	if !ok {
		return nil, errors.New("invalid user mac address")
	}
	topic := s.GetUserTopic(user)
	log.Println("publish message to user:", topic)
	return s.SendMessageToTopic(topic, message), nil
}

func (s *app) OnUserTopicDataReceived(client mqtt.Client, message mqtt.Message) {
	topic := message.Topic()
	fmt.Println("user topic: ", topic)
	temp := strings.Split(topic, "/")
	appId := temp[1]
	mac := temp[3]
	fmt.Println("send message to ", fmt.Sprintf("%v,%v,%v", appId, mac, string(message.Payload())))
}

func (s *app) OnIRReceived(client mqtt.Client, message mqtt.Message) {
	//fmt.Println("ir received", string(message.Payload()))
	query, err := url.ParseQuery(string(message.Payload()))
	if err != nil {
		fmt.Println("parse query data error:", err)
	}
	data := query.Get("data")
	v := `{label:"%s",value:"%s"},`
	fmt.Println(fmt.Sprintf(v, RandStringBytes(10), data))

	btn, _ := remotecontrol.NewButton(context.Background())
	po := btn.GetPO()
	po.AppName = s.options.Name
	po.Name = RandStringBytes(10)
	po.IrCode = data
	btn.Save()
}

func (s *app) OnHeartBeat(client mqtt.Client, message mqtt.Message) {
	fmt.Println("on heart beat message", fmt.Sprintf("%s", message.Payload()))
	now := time.Now().Unix()
	query, err := url.ParseQuery(string(message.Payload()))
	if err != nil {
		fmt.Println("parse query data error:", err)
	}
	mac := query.Get("mac")
	if mac == "" {
		return
	}
	if user, ok := s.Users[mac]; ok {
		fmt.Println("user existing: ", mac)
		user.Relay = query.Get("relay")
		user.HeartbeatAt = now
		user.IP = query.Get("ip")
		user.WIFI = query.Get("wifi")
		fmt.Println("user.Relay", user.Relay)
		saveUser(user)
	} else {
		fmt.Println("no user: ", mac)
		user := &DevicePO{
			AppName:     s.options.Name,
			ModeId:      []string{},
			Mac:         mac,
			WIFI:        query.Get("wifi"),
			IP:          query.Get("ip"),
			Name:        query.Get("clientId"),
			Relay:       query.Get("relay"),
			ConnectedAt: now,
			HeartbeatAt: now,
		}
		s.AddUser(user)
	}
	s.sendUsersToWS()
}

func (s *app) sendUsersToWS() error {
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
	client := s.options.client
	log.Println("subscribe to public ir received:", s.GetIRReceivedTopic())
	log.Println("app boardcast topic:", s.GetPublicTopic())
	fmt.Println("s.GetIRReceivedTopic()", s.GetIRReceivedTopic())
	client.Subscribe(s.GetIRReceivedTopic(), s.options.Qos, s.OnIRReceived)
	//client.Subscribe("ir-received", s.options.Qos, s.OnIRReceived)
	//log.Println("subscribe to public heart beat topic:", s.GetUserHeartBeatTopic())
	//client.Subscribe(s.GetUserHeartBeatTopic(), s.options.Qos, s.OnHeartBeat)
	s.Users = loadUsers(s.options.Name)
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
