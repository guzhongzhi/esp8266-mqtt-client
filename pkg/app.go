package pkg

import (
	"camera360.com/tv/pkg/app"
	"camera360.com/tv/pkg/dto"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"go.etcd.io/bbolt"
	"strings"
	"sync"
	"time"
)

const (
	userBucketName = "users"
)

var apps map[string]*App

func Apps() map[string]*App {
	return apps
}

func NewApp(appName, mq, relayBootStatus string) *App {
	if apps == nil {
		apps = make(map[string]*App)
	}
	if v, ok := apps[appName]; ok {
		return v
	}
	a := &App{
		Name:            appName,
		RelayBootStatus: relayBootStatus,
		MQTTServer:      mq,
		Users:           make(map[string]*app.User),
	}
	apps[appName] = a
	err := a.init()
	if err != nil {
		panic(err.Error())
	}
	return a
}

func GetApp(appName string) *App {
	if v, ok := apps[appName]; ok {
		return v
	}
	return nil
}

type App struct {
	Name            string               `json:"name"`
	RelayBootStatus string               `json:"relay_boot_status"`
	MQTTServer      string               `json:"mqtt_server"`
	Users           map[string]*app.User `json:"users"`
	BaseDir         string
	locker          sync.Mutex `json:"-"`
	mqClient        mqtt.Client
	db              *bbolt.DB
}

func (s *App) Close() error {
	if s.db != nil {
		s.db.Close()
	}

	return nil
}

func (s *App) init() error {
	db, err := bbolt.Open("data/"+s.Name+".db", 0700, nil)
	if err != nil {
		return err
	}
	s.db = db

	err = db.Update(func(tx *bbolt.Tx) error {
		var err error
		var bucket *bbolt.Bucket
		bucket, err = tx.CreateBucketIfNotExists([]byte(userBucketName))
		if err != nil {
			return err
		}

		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			user := &app.User{}
			json.Unmarshal(v, user)
			s.Users[user.Mac] = user
		}
		return err
	})
	return err
}

func (s *App) SetMQTTClient(c mqtt.Client) {
	s.mqClient = c
}

func (s *App) SendUserCommand(mac string, command *dto.Command) {
	fmt.Println("send command", command.ToJSON())
	s.mqClient.Publish("/"+s.Name+"/user/"+mac, 1, false, command.ToJSON())
}

func (s *App) addUser(req *dto.BeatRequest) error {
	if _, ok := s.Users[req.Mac]; !ok {
		s.locker.Lock()
		defer s.locker.Unlock()
		if _, ok := s.Users[req.Mac]; !ok {
			s.Users[req.Mac] = &app.User{
				Mac:       req.Mac,
				IP:        req.IP,
				WIFI:      req.WIFI,
				Gateway:   req.Gateway,
				ClientId:  req.ClientId,
				RelayPin:  req.RelayPIN,
				Relay:     req.Relay,
				StatePin:  req.StatePIN,
				IRPin:     req.IrPIN,
				AppName:   req.App,
				IsNewBoot: req.IsNewBoot,
				Version:   req.Version,
			}
		}
	}
	s.Users[req.Mac].RefreshedAt = int(time.Now().Unix())
	s.Users[req.Mac].IsNewBoot = req.IsNewBoot
	s.Users[req.Mac].WIFI = req.WIFI

	s.Users[req.Mac].RelayPin = req.RelayPIN
	s.Users[req.Mac].Relay = req.Relay
	s.Users[req.Mac].IRPin = req.IrPIN
	s.Users[req.Mac].IRReceivePin = req.IrReceivePIN
	s.Users[req.Mac].Version = req.Version

	buf, err := json.Marshal(s.Users[req.Mac])
	if err != nil {
		return err
	}
	s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(userBucketName))
		if b == nil {
			return fmt.Errorf("invalid bucket '%s'", userBucketName)
		}
		return b.Put([]byte(req.Mac), buf)
	})
	return nil
}

func (s *App) SaveUser(mac string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(userBucketName))
		if b == nil {
			return fmt.Errorf("invalid bucket '%s'", userBucketName)
		}

		if _, ok := s.Users[mac]; !ok {
			return fmt.Errorf("invalid user mac '%s'", mac)
		}
		buf, err := json.Marshal(s.Users[mac])
		if err != nil {
			return err
		}
		return b.Put([]byte(mac), buf)
	})
}

func (s *App) DeleteUser(mac string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(userBucketName))
		if b == nil {
			return fmt.Errorf("invalid bucket '%s'", userBucketName)
		}
		return b.Delete([]byte(mac))
	})
}

func (s *App) MQTTOnMessageReceived(mqttClient mqtt.Client, message mqtt.Message) {
	clientId := ""
	fmt.Println(string(message.Payload()), clientId)
	request := &dto.BeatRequest{}
	body := strings.TrimSpace(string(message.Payload()))
	if body[0] == '{' {
		json.Unmarshal([]byte(body), request)
		clientId = request.ClientId
	} else {
		return
	}
	if request.App != s.Name {
		return
	}
	switch request.Command {
	case "beat":
		s.addUser(request)
	case "irr":
		//红外接收
		data := request.Data
		s.irReceived(strings.TrimSpace(data))
	}
}

func (s *App) GetButtonGroups() []*app.ButtonGroup {
	data := make([]*app.ButtonGroup, 0)
	s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("button-groups"))
		if err != nil {
			return err
		}
		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			btn := &app.ButtonGroup{}
			json.Unmarshal(v, btn)
			data = append(data, btn)
		}

		return nil
	})
	return data
}

func (s *App) GetGroupButtons(groupId uint64) []*app.Button {
	data := make([]*app.Button, 0)
	s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("buttons"))
		if err != nil {
			return err
		}
		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			btn := &app.Button{}
			json.Unmarshal(v, btn)
			if btn.GroupId == groupId {
				data = append(data, btn)
			}
		}

		return nil
	})
	return data
}

func (s *App) GetButtons() []*app.Button {
	data := make([]*app.Button, 0)
	s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("buttons"))
		if err != nil {
			return err
		}
		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			btn := &app.Button{}
			json.Unmarshal(v, btn)
			data = append(data, btn)
		}

		return nil
	})
	return data
}

func (s *App) DeleteButtonGroup(id uint64) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("button-groups"))
		if err != nil {
			return err
		}
		bucket.Delete(itob(uint64(id)))
		return nil
	})
}
func (s *App) SaveButtonGroup(b *app.ButtonGroup) error {
	s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("button-groups"))
		if err != nil {
			return err
		}
		if b.Id == 0 {
			id, err := bucket.NextSequence()
			if err != nil {
				return err
			}
			b.Id = id
			b.CreatedAt = int32(time.Now().Unix())
		}
		b.UpdatedAt = int32(time.Now().Unix())
		if b.CreatedAt == 0 {
			b.CreatedAt = int32(time.Now().Unix())
		}
		js, err := json.Marshal(b)
		if err != nil {
			return err
		}
		d := bucket.Get(itob(b.Id))
		if d == nil || len(d) == 0 {
			old := app.ButtonGroup{}
			json.Unmarshal(d, &old)
			b.CreatedAt = old.CreatedAt
		} else {
			return bucket.Put(itob(b.Id), js)
		}
		return bucket.Put(itob(b.Id), js)
	})
	return nil
}

func (s *App) SaveButton(b *app.Button) error {
	s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("buttons"))
		if err != nil {
			return err
		}
		js, err := json.Marshal(b)
		if err != nil {
			return err
		}
		if b.CreatedAt == 0 {
			b.CreatedAt = int(time.Now().Unix())
		}
		return bucket.Put(itob(b.Id), js)
	})
	return nil
}

func (s *App) DeleteButton(id uint64) error {
	s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("buttons"))
		if err != nil {
			return err
		}
		return bucket.Delete(itob(uint64(id)))
	})
	return nil
}

func (s *App) irReceived(data string) {
	temp := strings.Split(data, ";")
	if len(temp) <2 {
		return
	}
	data = strings.Trim(temp[0], "{")
	data = strings.Trim(data, "}")
	data = strings.TrimSpace(data)
	nec := strings.TrimSpace(temp[1])
	nec = strings.ReplaceAll(nec, "//UNKNOWN", "")
	nec = strings.ReplaceAll(nec, "//NEC", "")
	nec = strings.ReplaceAll(nec, "//SAMSUNG", "")

	s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("buttons"))
		if err != nil {
			return err
		}
		id, err := bucket.NextSequence()
		if err != nil {
			return err
		}
		n := int(time.Now().Unix())
		btn := &app.Button{
			NEC:       nec,
			Data:      data,
			Id:        id,
			UpdatedAt: n,
			CreatedAt: n,
		}

		body, err := json.Marshal(btn)
		if err != nil {
			return err
		}
		return bucket.Put(itob(id), body)
	})
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
