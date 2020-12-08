package main

import (
	"camera360.com/tv/pkg/tv"
	"github.com/urfave/cli/v2" // imports as package "cli"
	"os"
	"sync"
)

var app = cli.NewApp()

func initCtx(ctx *cli.Context) {
	mq := ctx.String("mq")
	db := ctx.String("sqlite")
	tv.SetDbFileName(db)
	tv.SetMQServer(mq)
}

func main() {
	app.Commands = []*cli.Command{
		&cli.Command{
			Name:"init",
			Usage:"init db",
			Action: func(ctx *cli.Context) error {
				initCtx(ctx)
				tv.CreateTables()
				return nil
			},
		},
		&cli.Command{
			Name:"serve",
			Usage:"start http server and subscribe to mqtt server",
			Action: func(ctx *cli.Context) error {
				initCtx(ctx)
				var wg sync.WaitGroup
				wg.Add(2)
				tv.ServeMQTT()
				tv.ServeHttp()
				wg.Wait()
				return nil
			},
		},
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:"mq",
			Usage:"mqtt server url, such as 1.0.0.0:1883",
			Value:"tcp://192.168.18.60:1883",
		},
		&cli.StringFlag{
			Name:"sqlite",
			Usage:"sqlite file name",
			Value:"data.db",
		},
	}
	app.Run(os.Args)
}
