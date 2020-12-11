package main

import (
	"camera360.com/tv/pkg/runtime"
	"camera360.com/tv/pkg/tv"
	"github.com/urfave/cli/v2" // imports as package "cli"
	"os"
	"path/filepath"
	"sync"
)

var app = cli.NewApp()

func initCtx(ctx *cli.Context) {
	mq := ctx.String("mq")
	db := ctx.String("sqlite")
	base := ctx.String("base")
	runtime.PATH = base
	tv.SetDbFileName(db)
	tv.SetMQServer(mq)
}

func main() {
	app.Commands = []*cli.Command{
		&cli.Command{
			Name:  "init",
			Usage: "init db",
			Action: func(ctx *cli.Context) error {
				initCtx(ctx)
				tv.CreateTables()
				return nil
			},
		},
		&cli.Command{
			Name:  "serve",
			Usage: "start http server and subscribe to mqtt server",
			Action: func(ctx *cli.Context) error {
				initCtx(ctx)
				listen := ctx.String("listen")
				var wg sync.WaitGroup
				wg.Add(2)
				go tv.NewHub().Run()
				go func() {
					tv.ServeMQTT()
				}()
				go func() {
					tv.ServeHttp(listen)
				}()
				wg.Wait()
				return nil
			},
		},
	}

	dir := filepath.Dir(filepath.Dir(os.Args[0]))
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "mq",
			Usage: "mqtt server url, such as 1.0.0.0:1883",
			Value: "tcp://mqtt:mqtt@118.31.246.195:1883",
		},
		&cli.StringFlag{
			Name:  "sqlite",
			Usage: "sqlite file name",
			Value: "data.db",
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
	app.Run(os.Args)
}
