package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
	"x-server/example/chat-st/tools"

	"github.com/gorilla/websocket"
)

func ping(c *websocket.Conn) {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(3))
	if err := c.WriteMessage(websocket.BinaryMessage, buf.Bytes()); err != nil {
		log.Printf("websocket send ping err:%v", err)
	}
}

//1. 先造玩家

// -userId  11  -addr  8.219.59.226:81  -wsPath  xy3-cross/new-chat/Connect  -token  "Bearer eyJ1c2VyaWQiOjEwMDAwMTU2fQ==.03258446b8eb1b3cd6507dfc2737128b"
// go run main.go -userId 11 -zoneId xy3-1 -serverId 1 -addr 127.0.0.1:8002
func main() {
	flag.Parse()

	u := url.URL{
		Scheme:   "ws",
		Host:     tools.Addr,
		Path:     tools.WsPath,
		RawQuery: fmt.Sprintf("userid=%v", 1000),
	}
	log.Printf("connect to chat %v", u.String())

	reqheader := http.Header{}
	reqheader.Add("Authorization", "xxxx")

	c, _, err := websocket.DefaultDialer.Dial(u.String(), reqheader)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	go func() {
		t := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-t.C:
				ping(c)
			}
		}
	}()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read websocket err: ", err)
			break
		}
		flag := message[0]
		if flag == 0 {
			distribute(message[1:])
		}
	}
}
