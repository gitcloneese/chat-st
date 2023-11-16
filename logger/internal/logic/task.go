// Package logic
// 任务数据
package logic

import (
	"context"
	"x-server/core/model"
	"x-server/logger/internal/dao"

	pb "xy3-proto/logger"
	"xy3-proto/pkg/log"
)

func (p *Logic) Task(_ context.Context, req *pb.LogMsg) error {
	obj, err := model.UnmarshalToTLogTask(req)
	if err != nil {
		log.Error("Task Unmarshal req:%v error: %v", req, err)
		return err
	}
	if data := dao.GetMysqlDB().Create(
		obj,
	); data.RowsAffected != 1 {
		log.Error("Task Create err! req:%v err: insert error ", req, data)
	}
	return nil
}
