package tv

import (
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"math/rand"
	"net/url"
	"time"
)

var mqServer = ""
var client mqtt.Client

func SetMQServer(v string) {
	mqServer = v
}

func RegistryApp(mqttClient mqtt.Client, message mqtt.Message) {
	query, err := url.ParseQuery(string(message.Payload()))
	if err != nil {
		log.Println("init application failure:", err.Error())
		return
	}
	clientId := query.Get("clientId")
	app := NewApp(clientId, NewMQTTClientOption(mqttClient))
	app.OnHeartBeat(mqttClient, message)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func ServeMQTT() {
	opts := mqtt.NewClientOptions()
	log.Println("mqServer:", mqServer)
	temp, err := url.Parse(mqServer)
	if err != nil {
		log.Fatal("invalid mq server: ", mqServer)
	}
	opts.AddBroker(temp.Hostname() + ":" + temp.Port())
	opts.ConnectTimeout = time.Second * 5
	opts.SetClientID("server-" + RandStringBytes(5))
	opts.SetPingTimeout(1 * time.Second)
	opts.Username = temp.User.Username()
	opts.Password, _ = temp.User.Password()
	opts.OnConnect = func(client mqtt.Client) {
		log.Println("heart-beat")
		client.Subscribe("/camera360/heart-beat", 2, RegistryApp)
	}
	opts.OnConnectionLost = func(i mqtt.Client, e error) {
		log.Println("connect lost: ", e.Error())
	}
	opts.OnReconnecting = func(i mqtt.Client, options *mqtt.ClientOptions) {
		log.Println("reconnect")
	}
	var token mqtt.Token
	client = mqtt.NewClient(opts)
	token = client.Connect()
	token.Wait()
	if token.Error() != nil {
		log.Fatal("mqtt connect: ", token.Error(), client.IsConnected())
	}
	log.Println("mqtt connected")
}
