package logic

import (
	"context"

	"x-server/core/model"
	"x-server/logger/internal/dao"

	pb "xy3-proto/logger"
	"xy3-proto/pkg/log"
)

func (p *Logic) Online(ctx context.Context, req *pb.LogMsg) error {
	obj, err := model.UnmarshalToTLogOnline(req)
	log.Info("Online Unmarshal req:%v error: %v", req, obj)
	if err != nil {
		log.Error("Online Unmarshal req:%v error: %v", req, err)
		return err
	}
	data := dao.GetMysqlDB().Create(obj)
	log.Debug("Online Create req:%v, obj:%v, data:%v", req, obj, data)

	return nil
}
