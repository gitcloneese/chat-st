package main

import (
	"time"
	"x-server/example/login-st/tools"
)

const (
	RunAll             = iota //  运行所有测试用例
	RunPlatform               //  运行PlatformLogin
	RunAccountRoleList        //  运行AccountRoleList
	//RunGetGameLoginToken        //  运行GetGameLoginToken
)

// -platformAddr=http://8.219.160.79:82 -accountAddr=http://8.219.160.79:81 -accountNum=100
func main() {
	switch tools.T {
	case RunAll:
		// 游戏登录platform
		tools.PreparePlatformAccount()
		time.Sleep(time.Second)
		// 访问account对应的玩家列表
		tools.AccountRoleList()
		time.Sleep(time.Second)
		// 获取游戏登录token
		tools.GetLoginToken()
	case RunPlatform:
		tools.PreparePlatformAccount()
	case RunAccountRoleList:
		tools.PreparePlatformAccount()
		time.Sleep(time.Second * 5)
		tools.AccountRoleList()
	}
}
