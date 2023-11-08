package tools

import (
	"encoding/json"
	"google.golang.org/protobuf/proto"
	"log"
	newChat "xy3-proto/new-chat"
)

func cmdLogic(ops newChat.Operation, data []byte) {
	v, has := CmdM[ops]
	if !has || v == nil {
		log.Printf("ops %v not find parse message", ops)
		return
	}

	vv := v.ProtoReflect().New().Interface()

	err := proto.Unmarshal(data, vv)
	if err != nil {
		log.Printf("ops %v message parse err %v", ops, err)
		return
	}

	buf, err := json.MarshalIndent(v, "\t", "    ")
	if err != nil {
		log.Println("json marshal indent err ", err)
		return
	}

	log.Printf("ops %v message \n\t%v", ops, string(buf))
}
