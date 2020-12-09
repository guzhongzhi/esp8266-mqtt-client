package tv

type User struct {
	Id          int    `json:"id"`
	UserName    string `json:"username"`
	IP          string `json:"ip"`
	WIFI        string `json:"wifi"`
	Relay       string `json:"relay"`
	Mac         string `json:"mac"`
	ConnectedAt int64  `json:"connected_at"`
	HeartbeatAt int64  `json:"heartbeat_at"`
}

func (s *User) GetTopic() string {
	return s.Mac
}
