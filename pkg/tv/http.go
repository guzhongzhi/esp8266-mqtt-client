package tv

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func ServeHttp() {
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		js,err := json.Marshal(users)
		if err != nil {
			http.Error(w,err.Error(),http.StatusBadGateway)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	})

	http.HandleFunc("/ir", func(writer http.ResponseWriter, request *http.Request) {
		code := request.URL.Query().Get("code")
		sendIRData(code)
	})

	http.HandleFunc("/send", func(writer http.ResponseWriter, request *http.Request) {
		topic := request.URL.Query().Get("topic")
		cmd := request.URL.Query().Get("cmd")
		if topic == "" {
			topic = "camera360-global"
		}
		token := client.Publish(topic,0,false,cmd)
		token.Wait()
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("OK"))
	})

	p,_ :=filepath.Abs(filepath.Dir(filepath.Dir(os.Args[0])))
	fmt.Println("ppp",p)
	static := http.Dir(p + "/static/")
	fmt.Println("static",static)
	http.Handle("/static/",http.FileServer(static))

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		tmp := `<html><head>
    <meta charset="utf-8" />
    <title>Daemon</title> 
 <meta name="robots" content="index,follow" />
              <script src='https://libs.baidu.com/jquery/2.0.0/jquery.min.js'></script>
              <script src='https://cdnjs.cloudflare.com/ajax/libs/knockout/3.5.0/knockout-min.js'></script>
              <script src='http://esp8266.gulusoft.com/config.js'></script>
              <link rel='stylesheet' type='text/css' href='http://esp8266.gulusoft.com/main.css'>
 <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1" />
</head><body>

              <div class='title'><h1>控制中心<span id='loading'>Loading</span></h1></div>

              <div id='content'>

<div>
<a href="/send?cmd=off">电源关</a>
<a href="/send?cmd=on">电源开</a>
</div>
</div>
</body></html>
`
		writer.Write([]byte(tmp))
	})

	err := http.ListenAndServe("0.0.0.0:9900", nil)
	if err != nil {
		fmt.Printf("http.ListenAndServe()函数执行错误,错误为:%v\n", err)
		return
	}
}
