package logic

import (
	"context"

	"x-server/core/model"
	"x-server/logger/internal/dao"

	pb "xy3-proto/logger"
	"xy3-proto/pkg/log"
)

func (p *Logic) Mail(ctx context.Context, req *pb.LogMsg) error {
	log.Info("logger Mail: %v", req)

	if err := dao.GetMysqlDB().Create(model.UnmarshalToTLogMail(req)); err != nil {
		log.Error("Mail err! err:%v req:%v", err, req)
	}

	return nil
}
