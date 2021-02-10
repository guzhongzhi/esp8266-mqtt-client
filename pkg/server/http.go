package server

import (
	"camera360.com/tv/pkg"
	app2 "camera360.com/tv/pkg/app"
	"camera360.com/tv/pkg/dto"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func appByRequest(writer http.ResponseWriter, request *http.Request) *pkg.App {
	r := mux.Vars(request)
	if _, ok := r["appName"]; !ok {
		writer.Write([]byte("invalid appName"))
		return nil
	}
	appName := r["appName"]
	app := pkg.GetApp(appName)
	if app == nil {
		writer.Write([]byte("invalid appName"))
		return nil
	}
	return app
}
func ServeHttp(listen string) {
	log.Println("http listen: ", listen)
	r := mux.NewRouter()

	//静态文件
	dir := filepath.Dir(filepath.Dir(os.Args[0])) + "/static/"
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("OK"))
	})
	//websocket
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		WebSocketHandler(NewHub(), w, r)
	})

	//应用列表
	r.HandleFunc("/apps", func(writer http.ResponseWriter, request *http.Request) {
		apps := pkg.Apps()
		j, _ := json.Marshal(apps)
		writer.Write(j)
	})

	r.HandleFunc("/app/{appName}/users", func(writer http.ResponseWriter, request *http.Request) {
		app := appByRequest(writer, request)
		if app == nil {
			return
		}
		b, _ := json.Marshal(app.Users)
		writer.Write(b)
	})

	r.HandleFunc("/app/{appName}/user/send-message", func(writer http.ResponseWriter, request *http.Request) {
		app := appByRequest(writer, request)
		if app == nil {
			return
		}
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			fmt.Print("err:", err.Error())
		}
		fmt.Println("body", string(body))
		var req struct {
			Mac  string      `json:"mac"`
			Cmd  string      `json:"cmd"`
			Data interface{} `json:"data"`
		}
		json.Unmarshal(body, &req)
		if req.Data == nil {
			req.Data = ""
		}
		c := dto.NewCmd(req.Cmd, req.Data)
		if req.Mac != "" {
			app.SendUserCommand(req.Mac, c)
		} else {

		}
		writer.Write([]byte("{}"))
	})
	r.HandleFunc("/app/{appName}/buttons", func(writer http.ResponseWriter, request *http.Request) {
		app := appByRequest(writer, request)
		if app == nil {
			return
		}
		b, _ := json.Marshal(app.GetButtons())
		writer.Write(b)
	})

	r.HandleFunc("/app/{appName}/button/groups", func(writer http.ResponseWriter, request *http.Request) {
		app := appByRequest(writer, request)
		if app == nil {
			return
		}
		b, _ := json.Marshal(app.GetButtonGroups())
		writer.Write(b)
	})

	r.HandleFunc("/app/{appName}/button/group-save", func(writer http.ResponseWriter, request *http.Request) {
		app := appByRequest(writer, request)
		if app == nil {
			return
		}

		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			fmt.Print("err:", err.Error())
		}
		b := &app2.ButtonGroup{}
		json.Unmarshal(body, b)
		if b.Id == 0 {
			writer.Write([]byte("Error"))
		}
		app.SaveButtonGroup(b)
		writer.Write([]byte("Group Save OK"))
	})

	r.HandleFunc("/app/{appName}/button/group-delete", func(writer http.ResponseWriter, request *http.Request) {
		app := appByRequest(writer, request)
		if app == nil {
			return
		}

		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			fmt.Print("err:", err.Error())
		}
		b := &app2.ButtonGroup{}
		json.Unmarshal(body, b)
		if b.Id == 0 {
			writer.Write([]byte("Error"))
		}
		app.DeleteButtonGroup(uint64(b.Id))
		writer.Write([]byte("OK"))
	})

	r.HandleFunc("/app/{appName}/button/group-buttons", func(writer http.ResponseWriter, request *http.Request) {
		app := appByRequest(writer, request)
		if app == nil {
			return
		}

		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			fmt.Print("err:", err.Error())
		}
		b := &app2.ButtonGroup{}
		json.Unmarshal(body, b)
		if b.Id == 0 {
			writer.Write([]byte("Error"))
		}
		d := app.GetGroupButtons(uint64(b.Id))
		js, _ := json.Marshal(d)
		writer.Write(js)
	})

	r.HandleFunc("/app/{appName}/button/save", func(writer http.ResponseWriter, request *http.Request) {
		app := appByRequest(writer, request)
		if app == nil {
			return
		}

		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			fmt.Print("err:", err.Error())
		}
		b := &app2.Button{}
		json.Unmarshal(body, b)
		if b.Id == 0 {
			writer.Write([]byte("Error"))
		}
		app.SaveButton(b)
		writer.Write([]byte("OK"))
	})

	r.HandleFunc("/app/{appName}/button/delete", func(writer http.ResponseWriter, request *http.Request) {
		app := appByRequest(writer, request)
		if app == nil {
			return
		}

		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			fmt.Print("err:", err.Error())
		}
		b := &app2.Button{}
		json.Unmarshal(body, b)
		if b.Id == 0 {
			writer.Write([]byte("Error"))
			return
		}
		app.DeleteButton(b.Id)
		writer.Write([]byte("OK"))
	})

	r.HandleFunc("/app/{appName}/user/save", func(writer http.ResponseWriter, request *http.Request) {
		app := appByRequest(writer, request)
		if app == nil {
			return
		}
		b, _ := json.Marshal(app.Users)
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			fmt.Print("err:", err.Error())
		}
		var req struct {
			Mac      string  `json:"mac"`
			Groups   []int32 `json:"groups"`
			ClientId string  `json:"client_id"`
			CustomRelayPin int `json:"custom_relay_pin"`
			HasCustomRelayPin bool `json:"has_custom_relay_pin"`
			RelayPin int `json:"relay_pin"`
		}
		json.Unmarshal(body, &req)


		if v, ok := app.Users[req.Mac]; ok {
			oldHasCustomPin := v.HasCustomRelayPin
			v.ClientId = req.ClientId
			v.Groups = req.Groups
			v.CustomRelayPin = req.CustomRelayPin
			v.HasCustomRelayPin = req.HasCustomRelayPin
			v.RelayPin = req.RelayPin

			if oldHasCustomPin != req.HasCustomRelayPin || (v.HasCustomRelayPin  && v.CustomRelayPin != req.CustomRelayPin) {
				if req.HasCustomRelayPin {
					app.SendUserCommand(v.Mac,dto.NewCmd("srp", req.CustomRelayPin))
				} else {
					app.SendUserCommand(v.Mac,dto.NewCmd("srp", req.RelayPin))
				}
			}
		}


		app.SaveUser(req.Mac)
		writer.Write(b)
	})

	r.HandleFunc("/app/{appName}/user/delete", func(writer http.ResponseWriter, request *http.Request) {
		app := appByRequest(writer, request)
		if app == nil {
			return
		}
		b, _ := json.Marshal(app.Users)
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			fmt.Print("err:", err.Error())
		}
		var req struct {
			Mac  string `json:"mac"`
			Name string `json:"name"`
		}
		json.Unmarshal(body, &req)
		app.DeleteUser(req.Mac)
		writer.Write(b)
	})

	r.HandleFunc("/app/{appName}/dashboard", func(writer http.ResponseWriter, request *http.Request) {
		app := appByRequest(writer, request)
		if app == nil {
			return
		}
		fileName := app.BaseDir + "/static/index.html"
		body, err := ioutil.ReadFile(fileName)
		if err != nil {
			writer.Write([]byte("invalid appName: " + fileName))
			return
		}
		html := strings.Replace(string(body), "__APP_NAME__", app.Name, -1)
		writer.Write([]byte(html))
	})

	srv := &http.Server{
		Handler:      r,
		Addr:         listen,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("listen http server failed:", err.Error())
	}
}
