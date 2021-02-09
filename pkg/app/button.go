package app

type Button struct {
	Id uint64 `json:"id"`
	Name string  `json:"name"`
	Code string `json:"code"`
	NEC string `json:"nec"`
	Data string `json:"data"`
	GroupId uint64 `json:"group"` //遥控板
	CreatedAt int `json:"created_at"`
	UpdatedAt int `json:"updated_at"`
}
