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

func sendIRData(data string) mqtt.Token {
	cmd := "irs|" + data
	fmt.Println("cmd:",cmd)
	topic := "camera360-global"
	return client.Publish(topic,0,false,cmd)
}


func ServeMQTT() {
	opts := mqtt.NewClientOptions()
	fmt.Println("mqServer",mqServer);
	opts.AddBroker(mqServer)
	opts.SetClientID("aaaaaaa")
	opts.OnConnect = func(client mqtt.Client) {
		client.Subscribe("camera360-ir-received",2, func(client mqtt.Client, message mqtt.Message) {
			fmt.Println("ir received",string(message.Payload()))
			query,err := url.ParseQuery(string(message.Payload()))
			if err != nil {
				fmt.Println("parse query data error:",err)
			}
			data := query.Get("data")
			if data != "" {
				go func() {
					time.Sleep(time.Second)
					//sendIRData(data)  会死循环
				}()
			}
			fmt.Println("ir data:",data)
		})
		client.Subscribe("camera360-hart-beat",2, func(client mqtt.Client, message mqtt.Message) {
			fmt.Println("message",fmt.Sprintf("%s",message.Payload()))
			now := time.Now().Unix()
			query,err := url.ParseQuery(string(message.Payload()))
			if err != nil {
				fmt.Println("parse query data error:",err)
			}
			fmt.Println("query.Encode()",query.Encode())
			mac := query.Get("mac")
			if user,ok := users[mac];ok {
				user.Relay = query.Get("relay")
				user.HeartbeatAt = now
			} else {
				users[mac] = &User{
					Mac:mac,
					WIFI:query.Get("wifi"),
					IP:query.Get("ip"),
					UserName:query.Get("clientId"),
					Relay:query.Get("relay"),
					ConnectedAt:now,
					HeartbeatAt:now,
				}
			}

		})
	}
	client = mqtt.NewClient(opts)
	token := client.Connect()
	token.Wait()
}
