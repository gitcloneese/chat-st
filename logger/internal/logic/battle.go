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

func (p *Logic) Battle(_ context.Context, req *pb.LogMsg) error {
	obj, err := model.UnmarshalToTLogBattle(req)
	if err != nil {
		log.Error("Battle Unmarshal req:%v error: %v", req, err)
		return err
	}

	if data := dao.GetMysqlDB().Create(
		&(obj.Battle),
	); data.RowsAffected != 1 {
		log.Error("Battle Create err! insert error data: %v", data)
		return nil
	}

	obj.BattleParam.LogBattleID = obj.Battle.ID
	if data := dao.GetMysqlDB().Create(
		&(obj.BattleParam),
	); data.RowsAffected != 1 {
		log.Error("BattleParam Create err! insert error data: %v", data)
	}

	for index := range obj.BattleResults {
		obj.BattleResults[index].LogBattleID = obj.Battle.ID
	}
	if data := dao.GetMysqlDB().CreateInBatches(
		&(obj.BattleResults),
		len(obj.BattleResults)); data.RowsAffected == 0 {
		log.Error("BattleResults Create err! insert error data: %v", data)
	} else {
		log.Debug("[BATTLE LOG] Results: %v | Rows Affected:%v", obj.BattleResults, data.RowsAffected)
	}
	return nil
}
