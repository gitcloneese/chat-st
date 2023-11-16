package state

import (
	"context"
	pb "xy3-proto/coordinator"
	ec "xy3-proto/errcode"
	"xy3-proto/pkg/log"
)

func (s *State) HeartBeat(ctx context.Context, req *pb.HeartBeatReq) (resp *pb.CommonResp, err error) {
	if req.AppID == "" || req.Name == "" {
		log.Error("[HeartBeat] empty param service %s ID %s fail", req.Name, req.AppID)
		return nil, ec.IllegalParams
	}
	resp = &pb.CommonResp{}
	if req.ServiceShutdownFlag != 0 {
		err = s.DeleteServiceState2(context.TODO(), req.Name, req.AppID)
		if err != nil {
			log.Error("[HeartBeat] DeleteServiceState error %s", err.Error())
			return
		}
		log.Warn("HeartBeat service:%s app id:%s version:%d shutdown", req.Name, req.AppID, req.Version)
		return
	}

	err = s.UpdateServiceHeartBeat2(context.TODO(), req)
	if err != nil {
		log.Error("[HeartBeat] UpdateServiceHeartBeat2 service %s ID %s fail, msg %s", req.Name, req.AppID, err.Error())
		return
	}
	return
}
