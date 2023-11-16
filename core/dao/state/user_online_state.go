package state

import (
	"context"
	pb "xy3-proto/coordinator"
	ec "xy3-proto/errcode"
	"xy3-proto/pkg/log"
)

// PlayerOnline 调用 用于玩家建立长连接后把playerState改成1
func (s *State) PlayerOnline(ctx context.Context, req *pb.PlayerOnlineReq) (resp *pb.CommonResp, err error) {
	if req.PlayerID <= 0 {
		log.Error("PlayerOnline player id:%s wrong param", req.PlayerID)
		return nil, ec.IllegalParams
	}
	resp = &pb.CommonResp{}

	// 更新玩家状态为在线
	err = s.AddPlayerState(ctx, req.ServerName, req.AppId, req.PlayerID)
	if err != nil {
		log.Error("PlayerOnline player id:%d ZAdd err:%v", req.PlayerID, err)
		return
	}
	return
}

func (s *State) PlayerOffline(ctx context.Context, req *pb.PlayerOfflineReq) (resp *pb.CommonResp, err error) {
	if req.PlayerID <= 0 {
		log.Error("PlayerOffline player id:%s wrong param", req.PlayerID)
		return nil, ec.IllegalParams
	}
	resp = &pb.CommonResp{}

	// 玩家下线
	err = s.PlayerOfflineHandler(context.TODO(), req.PlayerID)
	if err != nil {
		log.Error("PlayerOffline PlayerOfflineHandler player id:%d err:%v", req.PlayerID, err)
		return
	}
	return
}

func (s *State) PlayerOfflineHandler(ctx context.Context, playerID int64) (err error) {
	// 检查玩家状态
	address, err := s.GetPlayerAllServiceAddress(context.TODO(), playerID)
	if err != nil {
		log.Error("[PlayerOfflineHandler] get player %d service address fail, msg %s", playerID, err.Error())
		return
	}

	err = s.DeletePlayerAndServiceState(context.TODO(), playerID, address)
	if err != nil {
		log.Error("[PlayerOfflineHandler] delete player %d service state fail, msg %s", playerID, err.Error())
		return
	}
	return
}
