package main

import (
	"x-server/example/login-st/tools"
)

// -platformAddr=http://8.219.160.79:82 -accountAddr=http://8.219.160.79:81 -accountNum=100
func main() {
	// 游戏登录platform
	tools.PreparePlatformAccount()
	// 访问account对应的玩家列表
	tools.AccountRoleList()
	// 获取游戏登录token
	tools.GetLoginToken()
}
