// Package logic
// 战斗结构
package logic

import (
	"context"
	"x-server/core/model"
	"x-server/logger/internal/dao"

	pb "xy3-proto/logger"
	"xy3-proto/pkg/log"
)

func (p *Logic) TutorialProgress(_ context.Context, req *pb.LogMsg) error {
	obj, err := model.UnmarshalToTLogTutorialProgress(req)
	if err != nil {
		log.Error("TutorialProgress Unmarshal req:%v error: %v", req, err)
		return err
	}
	if data := dao.GetMysqlDB().Create(
		obj,
	); data.RowsAffected != 1 {
		log.Error("TutorialProgress Create err! req:%v err: insert error ", req, data)
	}
	return nil
}
