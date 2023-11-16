package apolloConfig

import (
	"x-server/core/apollo"
	"xy3-proto/pkg/conf/paladin"
	"xy3-proto/pkg/log"
)

var (
	apolloConf paladin.Client
)

func init() {
	apollo.AddNs(
		apollo.DbNS,
		apollo.KafkaNS,
	)
}

// Init
// apollo配置初始化
func Init() {
	//其他配置检查依赖物品
	priorities := make([]string, 0, 3)
	configMap := make(map[string]paladin.Setter)
	fillNamespace(&priorities, configMap)
	var err error
	apolloConf, err = apollo.Init(priorities, configMap)
	if err != nil {
		panic(err)
	}
	value := apolloConf.Get("http.txt")
	log.Info("apolloConf get http.txt:%v", value)
}

func fillNamespace(_ *[]string, _ map[string]paladin.Setter) {
	//其他配置检查依赖物品
}
