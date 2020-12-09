package tv

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
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
		fmt.Println("init application failure:", err.Error())
		return
	}
	clientId := query.Get("clientId")
	app := NewApp(clientId, NewMQTTClientOption(mqttClient))
	app.OnHeartBeat(mqttClient,message)
}

func ServeMQTT() {
	opts := mqtt.NewClientOptions()
	fmt.Println("mqServer:", mqServer)
	opts.AddBroker(mqServer)
	opts.ConnectTimeout = time.Second * 5
	opts.SetClientID("server")
	opts.OnConnect = func(client mqtt.Client) {
		client.Subscribe("/camera360/heart-beat", 2, RegistryApp)
	}
	var token mqtt.Token
	client = mqtt.NewClient(opts)
	for !client.IsConnected() {
		fmt.Println("mqtt start to connect ")
		token = client.Connect()
		time.Sleep(time.Second * 5)
	}
	if token != nil {
		fmt.Println("mqtt connect: ", token.Error(), client.IsConnected())
	}
}
