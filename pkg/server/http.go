package server

import (
	"camera360.com/tv/pkg/ads"
	"camera360.com/tv/pkg/remotecontrol"
	"camera360.com/tv/pkg/tv"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func ServeHttp(listen string) {
	log.Println("http listen: ", listen)
	r := mux.NewRouter()

	//静态文件
	dir := filepath.Dir(filepath.Dir(os.Args[0])) + "/static/"
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("OK"))
	})

	//应用列表
	r.HandleFunc("/apps", func(writer http.ResponseWriter, request *http.Request) {
		j, err := json.Marshal(tv.Apps())
		fmt.Println("apps", err)
		writer.Write(j)
	})

	//websocket
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		tv.WebSocketHandler(tv.NewHub(), w, r)
	})

	r.HandleFunc("/api/index", func(writer http.ResponseWriter, request *http.Request) {
		ads.NewApi().ServeHTTP(writer, request)
	})

	//应用接口
	subRouter := r.PathPrefix("/app/{appId}/").Subrouter().MatcherFunc(func(r *http.Request, match *mux.RouteMatch) bool {
		result, _ := regexp.MatchString("/app/.*?", r.URL.Path)
		return result
	})
	subRouter.HandlerFunc(func(w http.ResponseWriter, r2 *http.Request) {
		//遥控板管理
		if strings.Index(r2.URL.Path, "/mode/") != -1 {
			remotecontrol.NewControl().ServeHTTP(w, r2)
		} else if strings.Index(r2.URL.Path, "/ads/") != -1 {
			ads.NewController().ServeHTTP(w, r2)
		} else {
			tv.NewApi().ServeHTTP(w, r2)
		}
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
