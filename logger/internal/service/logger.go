package service

import (
	"context"
	"encoding/json"
	"x-server/core/mq/kafka"
	"x-server/logger/internal/dao"
	"x-server/logger/internal/logic"
	pblogger "xy3-proto/logger"
	"xy3-proto/pkg/log"

	"google.golang.org/protobuf/types/known/emptypb"
)

// Service service.
type Service struct {
	pblogger.UnimplementedLoggerServer
	dao   dao.Dao
	logic *logic.Logic
}

// New new a service and return.
func New(d dao.Dao, receiver kafka.Receiver, logic *logic.Logic) (s *Service, cf func(), err error) {
	s = &Service{
		dao:   d,
		logic: logic,
	}
	cf = s.Close
	if err := receiver.Receive(context.Background(), s.Handler); err != nil {
		panic(err)
	}
	return
}

func (s *Service) Handler(_ context.Context, msg kafka.Event) error {
	var logMsg pblogger.LogMsgs
	err := json.Unmarshal(msg.Value(), &logMsg)
	if err != nil {
		log.Error("Message Handler MnMarshal log msg:%v err:%v", msg, err)
	}
	return s.logic.Handler(&logMsg)
}

// Close close the resource.
func (s *Service) Close() {
}

func (s *Service) Ping(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
