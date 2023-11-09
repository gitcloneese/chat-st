package tools

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	newChat "xy3-proto/new-chat"
)

func cmdLogic(ops newChat.Operation, data []byte) {
	v := operationMsg(ops)
	if v == nil {
		log.Errorf("ops %v not find parse message", ops)
		return
	}
	err := proto.Unmarshal(data, v)
	if err != nil {
		log.Errorf("ops %v message parse err %v", ops, err)
		return
	}

	buf, err := json.MarshalIndent(v, "\t", "    ")
	if err != nil {
		log.Errorf("json marshal indent err %v", err)
		return
	}

	log.Infof("ops %v message  %v", ops, string(buf))
}
