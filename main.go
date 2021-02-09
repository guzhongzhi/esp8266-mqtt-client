package main

import (
	"camera360.com/tv/pkg"
	"camera360.com/tv/pkg/server"
	"github.com/urfave/cli/v2" // imports as package "cli"
	"os"
	"path/filepath"
	"sync"
)

var app = cli.NewApp()

func main() {
	app.Commands = []*cli.Command{
		&cli.Command{
			Name:  "serve",
			Usage: "start http server and subscribe to mqtt server",
			Action: func(ctx *cli.Context) error {
				mq := ctx.String("mq")
				relayBootStatus := ctx.String("relayBootStatus")
				appName := ctx.String("appName")

				app := pkg.NewApp(appName,mq,relayBootStatus)
				app.BaseDir = ctx.String("base")
				//listen := ctx.String("listen")
				var wg sync.WaitGroup
				wg.Add(1)
				go server.ServeHttp(ctx.String("listen"))
				go server.ServeMQTT(app, nil)
				go server.NewHub().Run()
				wg.Wait()
				return nil
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "appName",
					Usage: "appName",
					Value: "guz",
				},
				&cli.StringFlag{
					Name:  "relayBootStatus",
					Usage: "device relay status after boot: on or off",
					Value: "off",
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
			Value: "",
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
		return nil
	}
	app.Run(os.Args)
}
