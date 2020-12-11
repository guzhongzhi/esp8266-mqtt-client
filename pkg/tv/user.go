package tv

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

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

func loadUsers(dir string) map[string]*User {
	data := make(map[string]*User)
	pathRoot := dir + string(os.PathSeparator)
	fileNames, err := filepath.Glob(pathRoot + "*")
	if err != nil {
		return data
	}
	fmt.Println(pathRoot+"*", fileNames)
	for _, fileName := range fileNames {
		content, err := ioutil.ReadFile(pathRoot + fileName)
		if err != nil {
			fmt.Println("err", err)
			continue
		}
		u := &User{}
		err = json.Unmarshal(content, u)
		if err != nil {
			fmt.Println("err", err)
			continue
		}
		u.Relay = "off"
		data[u.Mac] = u
	}
	fmt.Println(data)
	return data
}

func saveUser(dir string, user *User) error {
	content, err := json.Marshal(user)
	if err != nil {
		return err
	}
	s := strings.Replace(user.Mac, ":", "", -1)
	fileName := dir + string(os.PathSeparator) + s + ".json"
	return ioutil.WriteFile(fileName, content, os.ModePerm)
}
