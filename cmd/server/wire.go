//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	kratos "github.com/go-kratos/kratos/v2"
	"github.com/google/wire"

	"github.com/dreaminglwj/sage/internal/conf"
	"github.com/dreaminglwj/sage/internal/plugin/log"
	"github.com/dreaminglwj/sage/internal/server"
	"github.com/dreaminglwj/sage/internal/service"
	"github.com/dreaminglwj/sage/internal/storage"
	"github.com/dreaminglwj/sage/internal/storage/repository"
)

// wireApp init kratos application.
func wireApp(config *conf.Config) (*kratos.App, error) {
	panic(wire.Build(
		log.NewLogger,
		log.NewHelper,

		server.NewGRPCServer,
		storage.NewStorage,
		storage.NewEngine,

		repository.ProviderSet,
		service.NewSchemaService,
		newApp,
	))
}
