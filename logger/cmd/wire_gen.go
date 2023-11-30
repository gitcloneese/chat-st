// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"x-server/logger/internal/dao"
	"x-server/logger/internal/logic"
	"x-server/logger/internal/server"
	"x-server/logger/internal/service"
	"xy3-proto"
)

// Injectors from wire.go:

func initApp() (*kratos.App, func(), error) {
	db, err := dao.NewDB()
	if err != nil {
		return nil, nil, err
	}
	daoDao, cleanup, err := dao.NewDao(db)
	if err != nil {
		return nil, nil, err
	}
	receiver, cleanup2, err := dao.NewKafkaReceiver()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	logicLogic, err := logic.New(daoDao)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	serviceService, cleanup3, err := service.New(daoDao, receiver, logicLogic)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	httpServer, err := server.NewHTTPServer(serviceService)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	grpcServer, err := server.NewGRPCServer(serviceService)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	registrar := proto.NewRegister()
	app, cleanup4 := newApp(httpServer, grpcServer, registrar)
	return app, func() {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}