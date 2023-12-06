package main

import (
	"x-server/example/common/tools"
)

const (
	RunAll             = iota //  运行所有测试用例
	RunAccountRoleList        //  运行AccountRoleList
	//RunGetGameLoginToken        //  运行GetGameLoginToken
)

// -platformAddr=http://8.219.160.79:82 -accountAddr=http://8.219.160.79:81 -accountNum=100 --debug=true -c=100
// -platformAddr=http://8.219.160.79:82 -accountAddr=http://8.219.160.79:81 -accountNum=10
func main() {
	// 游客登录
	tools.RunPlatformGuestLoginReq()
	switch tools.T {
	case RunAll:
		// 访问account对应的玩家列表
		tools.RunAccountRoleListReq()
		// 获取游戏登录token
		tools.RunGameLoginReq()
	case RunAccountRoleList:
		tools.RunAccountRoleListReq()
	}
}
