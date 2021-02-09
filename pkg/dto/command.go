package dto

import (
	"encoding/json"
	"fmt"
	"time"
)

type Command struct {
	Cmd        string 		`json:"cmd"`
	Data       interface{}  `json:"data"`
	ExecutedAt int64  		`json:"execAt"`
}

func (s *Command) IsTurnOff() bool {
	return s.Cmd == "off"
}

func (s *Command) ToString() string {
	return fmt.Sprintf("%s,%v,%v", s.Cmd, s.ExecutedAt, s.Data)
}

func (s *Command) ToJSON() string {
	v,_ := json.Marshal(s)
	return string(v)
}

func NewCmd(cmd string, data interface{}) *Command {
	c := &Command{
		Cmd:        cmd,
		Data:       data,
		ExecutedAt: time.Now().Unix(),
	}
	return c
}
func NewIrSendCommand(data string) *Command {
	return NewCmd("irs", data)
}
func NewOffCommand() *Command {
	return NewCmd("off", "")
}

func NewOnCommand() *Command {
	return NewCmd("on", "")
}
