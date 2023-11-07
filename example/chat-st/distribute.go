package main

import (
	"google.golang.org/protobuf/proto"
	"log"
	chat "xy3-proto/new-chat"
)

func distribute(data []byte) {
	msg := &chat.Message{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		log.Println("distribute proto Unmarshal err:", err)
		return
	}
	cmdsLogic(msg.Ops, msg.Data)
}
