package main

import (
	"flag"
	"os"

	"github.com/dreaminglwj/sage/internal/resource"
	kratos "github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	"github.com/dreaminglwj/sage/internal/conf"
	"github.com/dreaminglwj/sage/internal/plugin/log"
)

var (
	flagConf string
	id, _    = os.Hostname()
)

func init() {
	flag.StringVar(&flagConf, "conf", "./configs", "config path, eg: -conf config.yaml")
}

func newApp(c *conf.Config, logger *log.Logger, gs *grpc.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(c.App.Name),
		kratos.Version(c.App.Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(gs),
	)
}

func main() {
	flag.Parse()
	defer resource.Release()
	c, err := conf.NewConfig(flagConf)
	if err != nil {
		panic(err)
	}

	app, err := wireApp(c)
	if err != nil {
		panic(err)
	}

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
