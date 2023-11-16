package scene

import (
	"context"
	"encoding/json"
	"fmt"

	"x-server/core/dao/model"
	battle "xy3-proto/battle"
	"xy3-proto/pkg/log"
	pbscene "xy3-proto/scene"

	v8 "github.com/go-redis/redis/v8"
)

func (s *Scene) GetCacheLineup(userid int64, group int32) (list []*battle.LineupInfo, powers []int64) {
	key := fmt.Sprintf(model.RedisKey_Lineup, group, userid)
	ll, err := s.client.LLen(context.TODO(), key).Result()
	if err != nil {
		return list, powers
	}
	result, err := s.client.LRange(context.TODO(), key, 0, ll-1).Result()
	if err != nil {
		return list, powers
	}
	// 从redis中按list方式取数据是从后到前,为了保证顺序应当从后开始遍历
	list = make([]*battle.LineupInfo, 0)
	powers = make([]int64, 0)
	for _, r := range result {
		lineupBytes := []byte(r)
		lineupInfo := &battle.LineupInfo{}
		err = json.Unmarshal(lineupBytes, lineupInfo)
		if err != nil {
			continue
		}
		list = append(list, lineupInfo)
		power := int64(0)
		for _, lineupItem := range lineupInfo.LineupItems {
			power += lineupItem.Hero.Power
		}
		power += lineupInfo.LeaderPower
		powers = append(powers, power)
	}
	return list, powers
}

func (s *Scene) GetCacheCampParam(userid int64, group int32) (list []*battle.CampParam) {
	list = []*battle.CampParam{}

	key := fmt.Sprintf(model.RedisKey_Battle, group, userid)
	result, err := s.client.LRange(context.TODO(), key, 0, -1).Result()
	if err != nil {
		log.Error("GetCacheCampParam! userid:%v group:%v err:%v", userid, group, err)
		return list
	}
	for _, v := range result {
		buffer := []byte(v)
		camppm := battle.CampParam{Fighters: []*battle.Fighter{}, LevelUpInfos: []*battle.LevelUpInfo{}}
		if err := json.Unmarshal(buffer, &camppm); err != nil {
			log.Error("GetCacheCampParam json err! userid:%v group:%v err:%v", userid, group, err)
			continue
		}
		list = append(list, &camppm)
	}

	return list
}

// CacheRoleBattle 缓存角色战斗数据和展示数据
func (s *Scene) CacheRoleLineup(id int64, m pbscene.GroupType, lineup []*battle.LineupInfo, battle []*battle.CampParam) (err error) {
	key := fmt.Sprintf(model.RedisKey_Lineup, int32(m), id)
	pipe := s.client.Pipeline()
	pipe.Del(context.TODO(), key)
	for _, l := range lineup {
		bytes, err := json.Marshal(l)
		if err != nil {
			continue
		}
		pipe.RPush(context.TODO(), key, bytes)
	}

	key = fmt.Sprintf(model.RedisKey_Battle, int32(m), id)
	pipe.Del(context.TODO(), key)
	for _, b := range battle {
		bytes, err := json.Marshal(b)
		if err != nil {
			continue
		}
		pipe.RPush(context.TODO(), key, bytes)
	}
	_, err = pipe.Exec(context.TODO())
	if err != nil {
		return err
	}
	return err
}

func GetPlayerLineups(redis *v8.Client, ids []int64, mi int32) (lineupMap map[int64][]*battle.LineupInfo, err error) {
	lineupMap = make(map[int64][]*battle.LineupInfo)
	pipeline := redis.Pipeline()
	for _, id := range ids {
		key := fmt.Sprintf(model.RedisKey_Lineup, mi, id)
		pipeline.LRange(context.TODO(), key, 0, -1)
	}
	result, err := pipeline.Exec(context.TODO())
	if err != nil {
		log.Error("GetPlayerLineups Exec err:%v", err)
		return lineupMap, err
	}

	for index, cmder := range result {
		var m []string
		m, err = cmder.(*v8.StringSliceCmd).Result()
		if err != nil {
			log.Error("GetPlayerLineups Result err:%v", err)
			continue
		}
		lineups := make([]*battle.LineupInfo, 0)
		for _, lineup := range m {
			lineupInfo := &battle.LineupInfo{}
			err = json.Unmarshal([]byte(lineup), lineupInfo)
			if err != nil {
				log.Error("GetPlayerLineups Unmarshal err:%v", err)
				continue
			}
			lineups = append(lineups, lineupInfo)
		}
		lineupMap[ids[index]] = lineups
	}
	return lineupMap, err
}
