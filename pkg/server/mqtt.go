package server

import (
	"camera360.com/tv/pkg"
	"camera360.com/tv/pkg/tools"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"net/url"
	"time"
)


func newServerId() string {
	return "server-" + tools.RandStringBytes(15)
}

func ServeMQTT(app *pkg.App, onConnectedCallback func(mqClient mqtt.Client) error) bool {
	statChan := make(chan bool, 1)
	clientId := newServerId()
	opts := mqtt.NewClientOptions()
	log.Println("mqServer:", app.MQTTServer)
	temp, err := url.Parse(app.MQTTServer)
	if err != nil {
		log.Fatal("invalid mq server: ", app.MQTTServer)
	}
	opts.AddBroker(temp.Hostname() + ":" + temp.Port())
	opts.ConnectTimeout = time.Second * 5
	opts.SetClientID(clientId)
	opts.SetPingTimeout(1 * time.Second)
	opts.Username = temp.User.Username()
	opts.Password, _ = temp.User.Password()
	opts.OnConnect = func(client mqtt.Client) {
		log.Println("on mqtt connected")
		if onConnectedCallback != nil {
			onConnectedCallback(client)
		}
		app.SetMQTTClient(client)
		t := client.Subscribe("/"+app.Name+"/heart-beat", 2, func(client mqtt.Client, message mqtt.Message) {
			app.MQTTOnMessageReceived(client,message)
			go func() {
				time.Sleep(time.Second * 1)
				hub.RefreshUsers(app.Name)
			}()
		})
		if t.Wait() && t.Error() == nil {
			statChan <- true
		}
	}
	opts.OnConnectionLost = func(i mqtt.Client, e error) {
		log.Println("connect lost: ", e.Error())
	}
	opts.OnReconnecting = func(i mqtt.Client, options *mqtt.ClientOptions) {
		log.Println("reconnect")
	}
	var token mqtt.Token
	client := mqtt.NewClient(opts)
	token = client.Connect()
	token.Wait()
	if token.Error() != nil {
		clientId = newServerId()
		log.Fatal("mqtt connect: ", token.Error(), client.IsConnected())
	}
	log.Println("mqtt connected")
	return <-statChan
}

