package tv

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type AppRequest struct {
	AppId string
	Mac   string
}

func (s *AppRequest) GetApp() App {
	fmt.Println("s.AppId", s.AppId)
	return NewApp(s.AppId)
}

func (s *AppRequest) GetUser() *User {
	app := s.GetApp()
	return app.GetUserByMac(s.Mac)
}

func NewAppRequest(r *http.Request) *AppRequest {
	vars := mux.Vars(r)
	appId, ok := vars["appId"]
	if !ok {
		appId = "camera360"
	}
	mac, ok := vars["mac"]
	if ok {
		mac = ""
	}
	newRequest := &AppRequest{
		AppId: appId,
		Mac:   mac,
	}
	return newRequest
}

func ServeHttp(listen string) {
	log.Println("http listen: ", listen)
	r := mux.NewRouter()
	dir := filepath.Dir(filepath.Dir(os.Args[0])) + "/static/"
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		tmp := `<html><head>
    <meta charset="utf-8" />
    <title>Daemon</title> 
 <meta name="robots" content="index,follow" />
              <script src='https://libs.baidu.com/jquery/2.0.0/jquery.min.js'></script>
              <script src='/static/knockout.js'></script>
              <script src='/static/config.js'></script>
<!-- 最新版本的 Bootstrap 核心 CSS 文件 -->
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@3.3.7/dist/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">

<!-- 可选的 Bootstrap 主题文件（一般不用引入） -->
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@3.3.7/dist/css/bootstrap-theme.min.css" integrity="sha384-rHyoN1iRsVXV4nD0JutlnGaslCJuC7uwjduW9SVrLvRYooPp2bWYgmgJQIXwl/Sp" crossorigin="anonymous">

<!-- 最新的 Bootstrap 核心 JavaScript 文件 -->
<script src="https://cdn.jsdelivr.net/npm/bootstrap@3.3.7/dist/js/bootstrap.min.js" integrity="sha384-Tc5IQib027qvyjSMfHjOMaLkfuWVxZxUPnCJA7l2mCWNIpG9mGCD8wGNIcPD7Txa" crossorigin="anonymous"></script>
              <link rel='stylesheet' type='text/css' href='/static/main.css'>
              <script src='/static/app.js'></script>
 <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1" />
</head><body>

              <div class='title'><h1>控制中心<span id='loading'>Loading</span></h1></div>
              <div id='content'></div>
</body></html>
`
		writer.Write([]byte(tmp))
	})

	r.HandleFunc("/apps", func(writer http.ResponseWriter, request *http.Request) {
		j, err := json.Marshal(Apps())
		fmt.Println("apps", err)
		writer.Write(j)
	})

	r.HandleFunc("/{appId}/users", func(w http.ResponseWriter, r *http.Request) {
		request := NewAppRequest(r)
		users := request.GetApp().GetUsers()
		js, err := json.Marshal(users)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	})

	r.HandleFunc("/{appId}/ir", func(writer http.ResponseWriter, request *http.Request) {
		code := request.URL.Query().Get("code")
		code = "irs|" + code
		NewAppRequest(request).GetApp().SendMessage(code)
	})

	r.HandleFunc("/{appId}/{mac}/ir", func(writer http.ResponseWriter, request *http.Request) {
		code := request.URL.Query().Get("code")
		code = "irs|" + code
		r := NewAppRequest(request)
		app := r.GetApp()
		user := r.GetUser()
		if user != nil {
			app.SendMessageToUser(user.Mac, code)
		} else {
			fmt.Println("send message to application user failure, user do not existing")
		}
	})
	r.HandleFunc("/{appId}/{mac}/message", func(writer http.ResponseWriter, request *http.Request) {
		appRequest := NewAppRequest(request)
		app := appRequest.GetApp()
		cmd := request.URL.Query().Get("cmd")
		user := appRequest.GetUser()
		if user != nil {
			app.SendMessageToUser(user.Mac, cmd)
		} else {
			fmt.Println("invalid user mac:", appRequest.Mac)
		}

		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("OK"))
	})
	r.HandleFunc("/{appId}/message", func(writer http.ResponseWriter, request *http.Request) {
		appRequest := NewAppRequest(request)
		app := appRequest.GetApp()
		topic := request.URL.Query().Get("topic")
		cmd := request.URL.Query().Get("cmd")
		if topic != "" {
			app.SendMessageToTopic(topic, cmd)
		} else {
			app.SendMessage(cmd)
		}

		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("OK"))
	})

	srv := &http.Server{
		Handler:      r,
		Addr:         listen,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("liste http server failed:", err.Error())
	}
}
