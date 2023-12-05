package main

import (
	"bytes"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"time"
	"x-server/example/chat-st/tools"
)

func ping(c *websocket.Conn) {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(3))
	if err := c.WriteMessage(websocket.BinaryMessage, buf.Bytes()); err != nil {
		log.Infof("websocket send ping err:%v", err)
	}
}

// 测试本地
// -addr=http://127.0.0.1:8200 -playerNum=100 -chatCount=100 -local=1 -loginAdd=http://127.0.0.1:8000  -t=1
// 测试远端
// -t=0 -addr=http://8.219.160.79:81 -playerNum=30 -chatCount=50 -local=0
// 正式环境
// -t=0 -addr=http://xy3api.firerock.sg -playerNum=30 -chatCount=50 -local=0
// uat环境
// -t=0 -addr=http://8.219.59.226:81 -playerNum=30 -chatCount=50 -local=0

const (
	TALl            = iota // 流程全跑一遍
	TSendMessage           // 发送消息
	TReceiveMessage        // 接收消息
)

func main() {
	tools.RunPreparePlayers()
	// 开始聊天测试
	time.Sleep(1 * time.Second)
	switch tools.T {
	case TALl: // 0
		tools.RunChat() // 发送消息, 接收消息
	case TSendMessage: // 1 // 发送多少条消息， 平均耗时， 成功率， 失败率
		tools.RunSendMessage()
	case TReceiveMessage: // 2 测试能建立多少ws长连接
		tools.RunTestReceiveMessage()
	}
	//监听os.Signal

}
