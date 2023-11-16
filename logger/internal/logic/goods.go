package logic

import (
	"context"
	"time"

	"x-server/core/model"
	"x-server/logger/internal/dao"

	pb "xy3-proto/logger"
	"xy3-proto/pkg/log"
)

func (p *Logic) Resource(ctx context.Context, req *pb.LogMsg) error {
	//log.Info("logger Resource: %v", req)

	obj, err := model.UnmarshalToTLogResource(req)
	if err != nil {
		log.Error("Resource Unmarshal req:%v error: %v", req, err)
		return err
	}
	if data := dao.GetMysqlDB().Create(&model.MLogTransfer{
		PlayerId:   obj.UserID,
		PlayerName: obj.UserName,
		OS:         obj.LogEmbed.OS,
		WarZone:    obj.WarZone,
		Server:     obj.Server,
		ItemID:     obj.ID,
		ItemNum:    obj.NowNum - obj.PreNum,
		ItemTypeId: obj.Type,
		ResultNum:  obj.NowNum,
		Source:     obj.Source,
		Time:       time.Unix(obj.LogEmbed.CreatedAt, 0),
		Level:      obj.Level,
	}); data.RowsAffected != 1 {
		log.Error("Resource Create err! req:%v err: insert error ", req, data)
	}
	return nil
}
