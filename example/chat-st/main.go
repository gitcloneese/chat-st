package main

import (
	"bytes"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func ping(c *websocket.Conn) {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(3))
	if err := c.WriteMessage(websocket.BinaryMessage, buf.Bytes()); err != nil {
		log.Printf("websocket send ping err:%v", err)
	}
}

//1. 先造玩家

// -addr = http: //127.0.0.1:8200 -playerNum=100 -local=1 -loginAdd=http://127.0.0.1:8000
func main() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	log.Println("exit")
}
