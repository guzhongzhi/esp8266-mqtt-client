package tv

import "github.com/eclipse/paho.mqtt.golang"

type AppOption func(opts *AppOptions)

func NewAppClientIdOption(id string) AppOption {
	return func(opts *AppOptions) {
		opts.Id = id
	}
}

func NewMQTTClientOption(client mqtt.Client) AppOption {
	return func(opts *AppOptions) {
		opts.client = client
	}
}

type AppOptions struct {
	client mqtt.Client
	Id     string
	Qos    byte
}

func NewAppOptions(opts ...AppOption) *AppOptions {
	o := &AppOptions{
		Qos: 2,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}
