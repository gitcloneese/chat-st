package main

import (
	"x-server/example/common/tools"
)

const (
	RunAll        = iota // 运行所有测试用例
	RunManualList        // 运行获取任务手册
	Friend               // 好友相关
)

// -platformAddr=http://8.219.160.79:82 -accountAddr=http://8.219.160.79:81 -accountNum=100 --debug=true -c=100
// -platformAddr=http://8.219.160.79:82 -accountAddr=http://8.219.160.79:81 -accountNum=10
// -c=50 -t=1 -platformAddr=https://xy3api.firerock.sg -addr=https://xy3api.firerock.sg -accountNum=3000 -chatCount=100 --debug=false -platformId=4
// -c=200 -t=1 -platformAddr=https://xy3api.firerock.sg -addr=https://xy3api.firerock.sg -accountNum=3000 -chatCount=100 --debug=false -testOne=true -n=10000 -platformId=4 -serverId=1
// -serverId=1 -platformId=4 -platformAddr=https://xy3api.firerock.sg -accountAddr=https://xy3api.firerock.sg -accountNum=10000 --debug=true -c=300 -testOne=false -n=1000 -t=2 -accountId=panll035
// -data='{"SeachParam":""}' -serverId=1 -platformId=4 -platformAddr=https://xy3api.firerock.sg -accountAddr=https://xy3api.firerock.sg -accountNum=50000 --debug=true -c=500 -testOne=true -n=1000 -t=2 -accountId=panll035
// -data='{"SeachParam":""}' -serverId=2 -platformAddr=http://8.219.160.79:82 -accountAddr=http://8.219.160.79:81 -accountNum=500000 --debug=true -c=500 -testOne=false -n=1000 -t=2 -accountId=panll035
func main() {
	if !tools.GetDBPlayer() {
		// 游客登录
		tools.RunPlatformGuestLoginReq()
		// 访问account对应的玩家列表
		tools.RunAccountRoleListReq()
		// 获取游戏登录token
		tools.RunGameLoginReq()
	}
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
	case Friend: // 好友
		tools.RunFriendRequestReq()
		tools.RunFriendRequestListReq()
		tools.RunFriendSearchReq()
	}
}
