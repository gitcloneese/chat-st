package tools

import (
	"google.golang.org/protobuf/proto"
	newChat "xy3-proto/new-chat"
	"xy3-proto/pkg/log"
)

func distribute(data []byte) {
	msg := &newChat.Message{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		log.Error("distribute proto Unmarshal err:%v", err)
		return
	}
	cmdLogic(msg.Ops, msg.Data)
}
