package logic

import (
	"context"

	"x-server/core/model"
	"x-server/logger/internal/dao"

	pb "xy3-proto/logger"
	"xy3-proto/pkg/log"
)

// PvpRanking 每日一次斗法日志
func (p *Logic) PvpRanking(_ context.Context, req *pb.LogMsg) error {
	obj, err := model.UnmarshalToMLogRanking(req)
	if err != nil {
		log.Error("PvpRanking Unmarshal req:%v error: %v", req, err)
		return err
	}
	data := dao.GetMysqlDB().Create(obj)
	log.Debug("PvpRanking Create req:%v, obj:%v, data:%v", req, obj, data)

	return nil
}
