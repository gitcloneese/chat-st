//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
	"x-server/logger/internal/dao"
	"x-server/logger/internal/logic"
	"x-server/logger/internal/server"
	"x-server/logger/internal/service"
)

func initApp() (*kratos.App, func(), error) {
	panic(wire.Build(dao.ProviderSet, server.ProviderSet, service.ProviderSet, logic.ProviderSet, newApp))
}
