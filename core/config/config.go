package config

import (
	herocfg "x-server/core/config/hero"
	"xy3-proto/pkg/conf/paladin"
)

// var Provider = wire.NewSet(InitConfig)

// Init 配置初始化
func Init(mm map[string]paladin.Setter) {
	herocfg.Init(mm)
}
