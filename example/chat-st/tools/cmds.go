package tools

import (
	"encoding/json"
	"log"

	"google.golang.org/protobuf/proto"
	newChat "xy3-proto/new-chat"
)

func cmdLogic(ops newChat.Operation, data []byte) {
	log.Printf("Recv Ops %v Data %v", ops, data)

	v, has := CmdM[ops]
	if !has || v == nil {
		log.Printf("ops %v not find parse message", ops)
		return
	}
	err := proto.Unmarshal(data, v)
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
