package tools

import (
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/proto"
	newChat "xy3-proto/new-chat"
	"xy3-proto/pkg/log"
)

func cmdLogic(ops newChat.Operation, data []byte) {
	v := operationMsg(ops)
	if v == nil {
		log.Error("ops %v not find parse message", ops)
		return
	}
	err := proto.Unmarshal(data, v)
	if err != nil {
		log.Error("ops %v message parse err %v", ops, err)
		return
	}

	buf, err := json.MarshalIndent(v, "\t", "    ")
	if err != nil {
		log.Error("json marshal indent err %v", err)
		return
	}

	log.Info("ops %v message  %v", ops, string(buf))
	fmt.Println(string(buf))
}
