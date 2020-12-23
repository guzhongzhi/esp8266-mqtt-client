package main

import (
	"camera360.com/tv/pkg/server"
	"camera360.com/tv/pkg/tv"
	"code.aliyun.com/MIG-server/micro-base/config"
	"code.aliyun.com/MIG-server/micro-base/microclient"
	"code.aliyun.com/MIG-server/micro-base/runtime"
	"code.aliyun.com/MIG-server/micro-base/utils"
	"fmt"
	"github.com/urfave/cli/v2" // imports as package "cli"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var app = cli.NewApp()

func main() {
	app.Commands = []*cli.Command{
		&cli.Command{
			Name:  "serve",
			Usage: "start http server and subscribe to mqtt server",
			Action: func(ctx *cli.Context) error {
				mq := ctx.String("mq")
				tv.SetMQServer(mq)
				appName := ctx.String("appName")
				listen := ctx.String("listen")
				var wg sync.WaitGroup
				wg.Add(2)
				tv.ServeMQTT(appName)
				go tv.NewHub().Run()
				go func() {
					server.ServeHttp(listen)
				}()
				wg.Wait()
				return nil
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "appName",
					Usage: "appName",
					Value: "camera360",
				},
			},
		},
	}

	dir := filepath.Dir(filepath.Dir(os.Args[0]))
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "env",
			Usage: "current　environment, values are: qa,dev,prod",
			Value: "dev",
		},
		&cli.StringFlag{
			Name:  "mq",
			Usage: "mqtt server url, such as 1.0.0.0:1883",
			Value: "tcp://mqtt:mqtt@118.31.246.195:1883",
		},
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "debug mode",
			Value: false,
		},
		&cli.StringFlag{
			Name:  "listen",
			Usage: "http listen port address and port",
			Value: "0.0.0.0:9900",
		},

		&cli.StringFlag{
			Name:  "base",
			Usage: "base directory",
			Value: dir,
		},
	}
	app.Before = func(ctx *cli.Context) error {
		runtime.SetDebug(ctx.Bool("debug"))
		env := ctx.String("env")
		clientId := "com.camera360.srv.tvads"
		//configUrl := "http://localhost:8100"
		configPath := utils.GetBinPath("../configs")
		loaderOptions := config.NewOptions(config.ConfigReloadDurationOption(time.Second*600),
			config.NewCallBackOption(func(loader *config.Loader) {
				for key, value := range loader.GetRemoteConfigData() {
					if runtime.IsDebug() {
						fmt.Println(key, value)
					}
				}
			}))
		_, err :=
			config.InitLoader(env, configPath,
				loaderOptions,
				//microclient.HttpCallUrlOption(configUrl),
				microclient.ClientCallTypeOption(microclient.ClientCallTypeHttp),
				microclient.NewClientIdOption(clientId))
		if err != nil {
			panic(err)
		}
		return nil
	}
	app.Run(os.Args)
}
