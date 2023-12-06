package main

import (
	"time"
	"x-server/example/common/tools"
)

// 正式环境
// -t=0 -addr=http://xy3api.firerock.sg
// uat环境
// -t=0 -addr=http://8.219.59.226:81
// 测试环境
//-c=50 -t=2 -platformAddr=http://8.219.160.79:82 -addr=http://8.219.160.79:81 -accountNum=100 -chatCount=100 --debug=true

const (
	TALl            = iota // 流程全跑一遍
	TSendMessage           // 发送消息
	TReceiveMessage        // 接收消息
)

// 默认是在第二服务器
func main() {
	// 游客登录
	tools.RunPlatformGuestLoginReq()
	// 访问account对应的玩家列表
	tools.RunAccountRoleListReq()
	// 获取游戏登录token
	tools.RunGameLoginReq()
	// 设置世界聊天频道
	tools.RunSetZoneServer()
	// 开始聊天测试
	switch tools.T {
	case TALl: // 0
		tools.RunChat() // 发送消息, 接收消息
	case TSendMessage: // 1 // 发送多少条消息， 平均耗时， 成功率， 失败率
		tools.RunSendMessage()
	case TReceiveMessage: // 2 测试能建立多少ws长连接
		tools.RunTestReceiveMessage()
	}
	time.Sleep(time.Second * 10000)
}
