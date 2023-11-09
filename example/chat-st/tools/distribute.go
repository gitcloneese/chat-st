package tools

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	newChat "xy3-proto/new-chat"
)

func distribute(data []byte) {
	msg := &newChat.Message{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		log.Errorf("distribute proto Unmarshal err:%v", err)
		return
	}
	cmdLogic(msg.Ops, msg.Data)
}
