package main

import (
	"x-server/example/common/tools"
)

const (
	RunAll        = iota //  运行所有测试用例
	RunManualList        //  运行获取任务手册
)

// -platformAddr=http://8.219.160.79:82 -accountAddr=http://8.219.160.79:81 -accountNum=100 --debug=true -c=100
// -platformAddr=http://8.219.160.79:82 -accountAddr=http://8.219.160.79:81 -accountNum=10
// -c=50 -t=1 -platformAddr=https://xy3api.firerock.sg -addr=https://xy3api.firerock.sg -accountNum=3000 -chatCount=100 --debug=false -platformId=4
// -c=200 -t=1 -platformAddr=https://xy3api.firerock.sg -addr=https://xy3api.firerock.sg -accountNum=3000 -chatCount=100 --debug=false -testOne=true -n=10000 -platformId=4
func main() {
	// 游客登录
	tools.RunPlatformGuestLoginReq()
	// 访问account对应的玩家列表
	tools.RunAccountRoleListReq()
	// 获取游戏登录token
	tools.RunGameLoginReq()
	switch tools.T {
	case RunAll:
		// 任务手册列表
		tools.RunManualListReq()
		// 任务手册总奖励领取状态
		tools.RunManualGrandTotalListReq()
	case RunManualList:
		// 任务手册
		// need 登录
		tools.RunManualListReq()
		tools.RunManualGrandTotalListReq()
	}
}
