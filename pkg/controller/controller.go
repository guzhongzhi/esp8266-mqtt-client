package controller

import (
	"code.aliyun.com/MIG-server/micro-base/orm/mongo"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Controller struct {
	response http.ResponseWriter
	request  *http.Request
}

func (c *Controller) WriteStatusData(data interface{}, status int, message string) error {
	js, _ := json.Marshal(mongo.M{
		"data":    data,
		"status":  status,
		"message": message,
	})
	c.response.Write(js)
	return nil
}

func (c *Controller) WriteJSON(data interface{}) error {
	js, err := json.Marshal(data)
	if err != nil {
		log.Println("WriteJSON error:", err.Error())
		return err
	}
	c.response.Write(js)
	return nil
}

func (c *Controller) SetResponse(w http.ResponseWriter) {
	c.response = w
}

func (c *Controller) SetRequest(r *http.Request) {
	c.request = r
}

func (c *Controller) Request() *http.Request {
	return c.request
}

func (c *Controller) Response() http.ResponseWriter {
	return c.response
}

func (c *Controller) value(name string) string {
	value := c.request.URL.Query().Get(name)
	if value == "" {
		value = c.request.FormValue(name)
	}
	return value
}

func (c *Controller) String(name string) string {
	return c.value(name)
}

func (c *Controller) Boolean(name string) bool {
	v := strings.ToLower(c.value(name))
	switch v {
	case "1", "true", "真", "yes":
		return true
	}
	return false
}

func (c *Controller) Int(name string) int {
	v := c.value(name)
	if v == "" {
		return 0
	}
	k, err := strconv.Atoi(v)
	if err != nil {
		log.Println("get int value failure: ", v)
	}
	return k
}
