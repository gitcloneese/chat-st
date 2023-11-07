package tools

import (
	"google.golang.org/protobuf/proto"
	"log"
	newChat "xy3-proto/new-chat"
)

func distribute(data []byte) {
	msg := &newChat.Message{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		log.Println("distribute proto Unmarshal err:", err)
		return
	}
	cmdLogic(msg.Ops, msg.Data)
}
