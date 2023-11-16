package config

import (
	"x-server/core/apollo"
	"x-server/kafkaSend/config/apolloConfig"
	"xy3-proto/pkg/conf/paladin"
	"xy3-proto/pkg/log"
)

var (
	priority = make([]string, 0, 10)
	mm       = make(map[string]paladin.Setter)
)

func init() {
	log.Info("config init")
}

func Init() error {
	ch := make(chan struct{})
	// 文件将废弃, 不在挂载文件
	go fileInit(ch)
	<-ch
	if apollo.Switch() {
		apolloConfig.Init()
	}
	return nil
}

// using apollo, this should not panic
func fileInit(ch chan struct{}) {
	defer func() {
		if err1 := recover(); err1 != nil {
			log.Error("file init err:%v", err1)
		}
		ch <- struct{}{}
	}()
	if err0 := paladin.Init(); err0 != nil {
		log.Error("fileInit paladin err:%v", err0)
		return
	}
	paladin.Watch(priority, mm)
}
