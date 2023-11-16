package main

import (
	"flag"
	"time"
	"x-server/core/pkg/util"
	"x-server/kafkaSend/config"
	pb "xy3-proto/logger"
)

func main() {

	flag.Parse()
	if err := config.Init(); err != nil {
		panic(err)
	}

	msg := &pb.LogMsgs{
		Messages: []*pb.LogMsg{
			{
				Category: pb.ELogCategory_ELC_Login,
				Os:       1,
				Time:     time.Now().Unix(),
				Json:     `{"test": 1}`,
			},
		},
	}
	err := util.Record(msg)
	if err != nil {
		panic(err)
	}
}
