package app


type User struct {
	Mac string `json:"mac"`
	IP string `json:"ip"`
	WIFI string `json:"wifi"`
	Gateway string `json:"gateway"`
	ClientId string `json:"client_id"`
	Relay string  `json:"relay"`//继电器状态 
	RelayPin int `json:"relay_pin"` //继电器引脚
	StatePin int  `json:"state_pin"`//状态引脚
	IRPin int  `json:"ir_pin"`//红外引脚
	AppName  string `json:"app_name"`
	IsNewBoot bool `json:"is_new_boot"`
	Version string `json:"version"`
	RefreshedAt int `json:"refreshed_at"`

	OnBootCommand			 string 			`json:"on_boot_command" bson:"onBootCommand"` //开机执行命令
	CustomRelayPin 			 int 				`json:"custom_relay_pin" bson:"customRelayPin"`
	HasCustomRelayPin		 bool 				`json:"has_custom_relay_pin" bson:"hasCustomRelayPin"`
	RelayTriggeredByLowLevel bool               `json:"relay_triggered_by_low_level" bson:"relayTriggeredByLowLevel"`

	Groups []int32 `json:"groups"`
}
