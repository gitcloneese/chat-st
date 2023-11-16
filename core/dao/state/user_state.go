package state

import (
	"context"
	"errors"
	"strconv"
	"x-server/core/dao/model"
	"x-server/core/pkg/util"
	"xy3-proto/pkg/log"

	redis "github.com/go-redis/redis/v8"
)

func (s *State) GetPlayerState(ctx context.Context, playerId int64) (state int64, err error) {
	result, err := s.client.Get(context.TODO(), model.GetPlayerStateKey(playerId)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return -1, err
	}
	if result == "" {
		return
	}
	state, err = strconv.ParseInt(result, 10, 64)
	return
}

func (s *State) AddPlayerState(ctx context.Context, serverName, appId string, playerId int64) error {
	// player 添加到 服务器玩家信息表
	pipe := s.client.Pipeline()
	pipe.SAdd(context.TODO(), model.GetServiceUserKey(serverName, appId), playerId)
	pipe.Set(context.TODO(), model.GetPlayerStateKey(playerId), util.GetTimeStamp(), 0)
	_, err := pipe.Exec(context.TODO())
	return err
}

func (s *State) GetPlayerOneServiceAddress(ctx context.Context, playerID int64, serviceName string) (addrMap map[string]string, err error) {
	addrMap = make(map[string]string)

	result, err := s.client.HGet(context.TODO(), model.GetUserServerListKey(playerID), serviceName).Result()
	if err != nil {
		return nil, err
	}
	addrMap[serviceName] = result
	return
}

func (s *State) GetMutliPlayerOneServiceAddress(ctx context.Context, playerIDS []int64, serviceName string) (addrMap map[int64]string, err error) {
	addrMap = make(map[int64]string)
	pipe := s.client.Pipeline()
	for _, v := range playerIDS {
		pipe.HGet(context.TODO(), model.GetUserServerListKey(v), serviceName)
	}

	result, err := pipe.Exec(context.TODO())
	if err != nil {
		log.Error("GetMutliPlayerOneServiceAddress service state fail, msg %s", err.Error())
	}
	for k, cmder := range result {
		var m string
		m, err = cmder.(*redis.StringCmd).Result()
		if err != nil {
			log.Error("GetMutliPlayerOneServiceAddress Result err:%s", err.Error())
			continue
			// return nil, err
		}
		addrMap[playerIDS[k]] = m
	}
	return
}

func (s *State) GetPlayerAllServiceAddress(ctx context.Context, playerID int64) (addrMap map[string]string, err error) {
	return s.client.HGetAll(context.TODO(), model.GetUserServerListKey(playerID)).Result()
}

func (s *State) UpdatePlayerState(ctx context.Context, playerID int64, addressMap map[string]string, state int64) (err error) {
	pipe := s.client.Pipeline()
	var userBindings []interface{}
	for k, v := range addressMap {
		// player 添加到 服务器玩家信息表
		pipe.SAdd(context.TODO(), model.GetServiceUserKey(k, v), playerID)

		userBindings = append(userBindings, k, v)
	}

	// 创建玩家绑定服务器信息
	if len(userBindings) > 0 && len(userBindings)%2 == 0 {
		pipe.HMSet(context.TODO(), model.GetUserServerListKey(playerID), userBindings...)
	} else {
		log.Warn("UpdatePlayerState player %d no user binding service ", playerID)
	}

	// 更新玩家状态
	pipe.Set(context.TODO(), model.GetPlayerStateKey(playerID), state, 0)
	_, err = pipe.Exec(context.TODO())
	if err != nil {
		log.Error("UpdatePlayerState update player %d and service state fail, msg %s", playerID, err.Error())
		return
	}
	return
}

func (s *State) DeletePlayerAndServiceState(ctx context.Context, playerID int64, addressMap map[string]string) (err error) {
	pipe := s.client.Pipeline()
	for k, v := range addressMap {
		// player 从 服务器玩家信息表 删除
		pipe.SRem(context.TODO(), model.GetServiceUserKey(k, v), playerID)
		// 更新服务器人数
		pipe.HIncrBy(context.TODO(), model.GetServerHostsKey2(k, v), StsPlayerNum, -1)
	}

	// 删除玩家绑定服务器信息表
	pipe.Del(context.TODO(), model.GetUserServerListKey(playerID))

	// 删除玩家状态
	pipe.Del(context.TODO(), model.GetPlayerStateKey(playerID))

	_, err = pipe.Exec(context.TODO())
	if err != nil {
		log.Error("delete player %d and service state fail, msg %s", playerID, err.Error())
		return
	}
	return
}
