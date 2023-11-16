package dao

import (
	"context"
	"fmt"
	"strconv"
	pbscene "xy3-proto/scene"

	"x-server/core/dao/model"
	coremodel "x-server/core/model"
	battle "xy3-proto/battle"
	"xy3-proto/pkg/log"

	"google.golang.org/protobuf/proto"
)

func (d *dao) CrossGetLineupOnline(userid int64, group int32) (rsp *pbscene.CrossGetLineupRsp) {
	rsp = &pbscene.CrossGetLineupRsp{}

	//获取玩家所在服务器ID
	result, err := d.client.HGet(context.TODO(), fmt.Sprintf(model.RedisKey_Player, userid), "ServerID").Result()
	if err != nil {
		log.Error("CrossGetLineupOnline redis-player err! userid:%v group:%v err:%v", userid, group, err)
		return rsp
	}
	serverid, _ := strconv.ParseInt(result, 10, 64)

	//获取玩家所在服的scene分线
	appid, err := d.Scene.PlayerLine(userid)
	if err != nil {
		log.Error("CrossGetLineupOnline redis-line err! userid:%v group:%v err:%v", userid, group, err)
		return rsp
	}

	req := &pbscene.CrossSceneReq{
		UserID: userid,
		Ops:    int32(pbscene.GameCommand_CrossGetLineup),
	}
	req.Data, err = proto.Marshal(&pbscene.CrossGetLineupReq{GroupID: pbscene.GroupType(group)})
	if err != nil {
		log.Error("CrossGetLineupOnline Marshal err! userid:%v group:%v err:%v", userid, group, err)
		return rsp
	}

	msg, err := pbscene.CrossRPC(context.TODO(), fmt.Sprintf(coremodel.NamespaceV, serverid), appid, req)
	if err != nil {
		log.Error("CrossGetLineupOnline rpc err! userid:%v group:%v err:%v", userid, group, err)
		return rsp
	}

	err = proto.Unmarshal(msg.Data, rsp)
	if err != nil {
		log.Error("CrossGetLineupOnline Unmarshal err! userid:%v group:%v err:%v", userid, group, err)
		return rsp
	}

	log.Info("CrossGetLineupOnline ok! userid:%v group:%v", userid, group)
	return rsp
}

func (d *dao) GetLineupFromOnlineOrRedis(userid int64, group int32) (lineups []*battle.LineupInfo, powers []int64, camps []*battle.CampParam) {
	rsp := d.CrossGetLineupOnline(userid, group)

	if rsp == nil || len(rsp.Lineups) == 0 { //没有查询到在线信息，需要使用redis缓存
		lineups, powers = d.Scene.GetCacheLineup(userid, group)
		camps = d.Scene.GetCacheCampParam(userid, group)

		log.Info("GetLineupFromOnlineOrRedis redis! userid:%v group:%v", userid, group)
	} else { //查询到玩家在线阵容，直接使用
		lineups = rsp.Lineups
		camps = rsp.Camps

		powers = make([]int64, 0, len(lineups))
		for _, lineup := range lineups {
			num := int64(0)
			for _, item := range lineup.LineupItems {
				num += item.Hero.Power
			}
			powers = append(powers, num)
		}

		log.Info("GetLineupFromOnlineOrRedis online! userid:%v group:%v", userid, group)
	}
	return
}

// CrossGetMainLineupOnline
// 获取正在使用的主阵容信息
func (d *dao) CrossGetMainLineupOnline(userid int64) (rsp *pbscene.CrossGetMainLineupRsp) {
	rsp = &pbscene.CrossGetMainLineupRsp{}

	//获取玩家所在服务器ID
	result, err := d.client.HGet(context.TODO(), fmt.Sprintf(model.RedisKey_Player, userid), "ServerID").Result()
	if err != nil {
		log.Error("CrossGetMainLineupOnline redis-player err! userid:%v err:%v", userid, err)
		return rsp
	}
	serverid, _ := strconv.ParseInt(result, 10, 64)

	//获取玩家所在服的scene分线
	appid, err := d.Scene.PlayerLine(userid)
	if err != nil {
		log.Error("CrossGetMainLineupOnline redis-line err! userid:%v err:%v", userid, err)
		return rsp
	}

	req := &pbscene.CrossSceneReq{
		UserID: userid,
		Ops:    int32(pbscene.GameCommand_CrossGetMainLineup),
	}
	req.Data, err = proto.Marshal(&pbscene.CrossGetMainLineupReq{})
	if err != nil {
		log.Error("CrossGetMainLineupOnline Marshal err! userid:%v err:%v", userid, err)
		return rsp
	}

	msg, err := pbscene.CrossRPC(context.TODO(), fmt.Sprintf(coremodel.NamespaceV, serverid), appid, req)
	if err != nil {
		log.Error("CrossGetMainLineupOnline rpc err! userid:%v err:%v", userid, err)
		return rsp
	}

	err = proto.Unmarshal(msg.Data, rsp)
	if err != nil {
		log.Error("CrossGetMainLineupOnline Unmarshal err! userid:%v err:%v", userid, err)
		return rsp
	}

	log.Info("CrossGetMainLineupOnline ok! userid:%v", userid)
	return rsp
}

// GetMainLineupFromOnlineOrRedis
// 获取正在使用的主阵容信息
func (d *dao) GetMainLineupFromOnlineOrRedis(userid int64) (lineups []*battle.LineupInfo, powers []int64, camps []*battle.CampParam) {
	rsp := d.CrossGetMainLineupOnline(userid)
	if rsp == nil || len(rsp.Lineups) == 0 { //没有查询到在线信息，需要使用redis缓存
		group := int32(pbscene.GroupType_LineupMain)
		lineups, powers = d.Scene.GetCacheLineup(userid, group)
		camps = d.Scene.GetCacheCampParam(userid, group)
		log.Info("GetMainLineupFromOnlineOrRedis redis! userid:%v group:%v", userid, group)
	} else { //查询到玩家在线阵容，直接使用
		lineups = rsp.Lineups
		camps = rsp.Camps

		powers = make([]int64, 0, len(lineups))
		for _, lineup := range lineups {
			num := int64(0)
			for _, item := range lineup.LineupItems {
				num += item.Hero.Power
			}
			powers = append(powers, num)
		}

		log.Info("GetMainLineupFromOnlineOrRedis online! userid:%v", userid)
	}
	return
}
