module camera360.com/tv

go 1.13

require (
	github.com/eclipse/paho.mqtt.golang v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/robfig/cron/v3 v3.0.0
	github.com/urfave/cli/v2 v2.3.0
	go.etcd.io/bbolt v1.3.5
	go.mongodb.org/mongo-driver v1.4.2
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
