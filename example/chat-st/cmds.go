package main

import (
	"log"
	"x-server/example/chat-st/tools"

	"encoding/json"

	"google.golang.org/protobuf/proto"
	chat "xy3-proto/new-chat"
)

func cmdsLogic(ops chat.Operation, data []byte) {
	log.Printf("Recv Ops %v Data %v", ops, data)

	v, has := tools.CmdM[ops]
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

	log.Printf("ops %v message \n\t%v", chat.Operation(ops), string(buf))
}
