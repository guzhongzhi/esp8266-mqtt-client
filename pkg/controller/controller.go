package controller

import (
	"code.aliyun.com/MIG-server/micro-base/orm/mongo"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type Controller struct {
	response    http.ResponseWriter
	request     *http.Request
	actionIndex int
}

func (c *Controller) ActionIndex(v int) {
	c.actionIndex = v
}

func (c *Controller) Dispatch(ctl interface{}) {
	r := c.request
	path := strings.Split(r.URL.Path, "/")
	if c.actionIndex == 0 {
		c.actionIndex = 3
	}
	methodName := path[c.actionIndex]

	temp := strings.Split(methodName, "-")
	if len(temp) > 1 {
		methodName = ""
		for _, v := range temp {
			methodName += strings.ToUpper(string(v[0])) + string(v[1:])
		}
	} else {
		methodName = strings.ToUpper(string(methodName[0])) + string(methodName[1:])
	}
	log.Println(fmt.Sprintf("methodName: %v,%T", methodName, ctl))
	st := reflect.ValueOf(ctl)
	v := st.MethodByName(methodName)
	res := v.Call([]reflect.Value{})
	if len(res) == 0 {
		return
	}
	v1 := res[0].Interface()
	switch v1.(type) {
	case error:
		c.WriteStatusData(nil, http.StatusInternalServerError, fmt.Sprintf("%v", v1))
	}
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
	var js []byte
	var err error
	switch data.(type) {
	case []byte:
		js = data.([]byte)
	case string:
		js = []byte(data.(string))
	default:
		js, err = json.Marshal(data)
		if err != nil {
			log.Println("WriteJSON error:", err.Error())
			return err
		}
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

func (c *Controller) Int64(name string) int64 {
	v := c.Int(name)
	return int64(v)
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
