package main

import (
	"bytes"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"x-server/example/chat-st/tools"
)

func ping(c *websocket.Conn) {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(3))
	if err := c.WriteMessage(websocket.BinaryMessage, buf.Bytes()); err != nil {
		log.Printf("websocket send ping err:%v", err)
	}
}

// 测试本地
// -addr = http://127.0.0.1:8200 -playerNum=100 -chatCount=100 -local=1 -loginAdd=http://127.0.0.1:8000  -t=1
// 测试远端
// -t=0 -addr=http://8.219.160.79:81 -playerNum=30 -chatCount=50 -local=0
func main() {
	now := time.Now()
	tools.PreparePlayers()
	// 开始聊天测试
	tools.PrepareChat()
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	latency := time.Since(now).Seconds()
	log.Printf("exit 总耗时:%v\n", latency)
}
