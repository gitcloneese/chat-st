package logic

import (
	"context"
	"x-server/core/model"
	"x-server/logger/internal/dao"

	pb "xy3-proto/logger"
	"xy3-proto/pkg/log"
)

func (p *Logic) Login(ctx context.Context, req *pb.LogMsg) error {
	obj, err := model.UnmarshalToTLogLogin(req)
	if err != nil {
		log.Error("Login Unmarshal req:%v error: %v", req, err)
		return err
	}
	if data := dao.GetMysqlDB().Create(obj); data.RowsAffected != 1 {
		log.Error("Login Create err! req:%v data:%v insert error ", req, data)
	}
	return nil
}

func (p *Logic) Register(ctx context.Context, req *pb.LogMsg) error {
	obj, err := model.UnmarshalToTLogRegister(req)
	if err != nil {
		log.Error("Register Unmarshal req:%v error: %v", req, err)
		return err
	}
	if data := dao.GetMysqlDB().Create(obj); data.RowsAffected != 1 {
		log.Error("Register Create err! req:%v data%v: insert error ", req, data)
	}
	return nil
}
