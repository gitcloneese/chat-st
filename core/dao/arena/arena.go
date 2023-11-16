package arena

import (
	"context"
	"errors"
	"fmt"
	"x-server/core/pkg/util"

	v8 "github.com/go-redis/redis/v8"

	"x-server/core/dao/model"
	"xy3-proto/pkg/log"
)

type Arena struct {
	client *v8.Client
}

func New(r *v8.Client) *Arena {
	return &Arena{client: r}
}

// 记录斗法每日排行
func (a *Arena) SaveDailyRank(serverID int64, pvpType int32, day string, dailyRankMap map[int64]int32) (err error) {
	pipe := a.client.Pipeline()
	for uid, dailyRank := range dailyRankMap {
		key := fmt.Sprintf(model.RedisKey_Arena_Daily_Rank, serverID, pvpType, day)
		_, err = pipe.HSet(context.TODO(), key, util.Int64ToStr(uid), dailyRank).Result()
		if err != nil {
			log.Error("SaveDailyRank HSet Error: %v", err)
		}
	}
	_, err = pipe.Exec(context.TODO())
	if err != nil {
		log.Error("SaveDailyPvpFightingRank err:[%v]")
		return
	}
	return
}

// 获取斗法每日排行
func GetDailyRank(redis *v8.Client, serverID int64, pvpType int32, day string) (dailyRankMap map[int64]int32, err error) {
	key := fmt.Sprintf(model.RedisKey_Arena_Daily_Rank, serverID, pvpType, day)
	result, err := redis.HGetAll(context.TODO(), key).Result()
	if err != nil {
		log.Error("GetDailyPvpFightingRank day:[%v] err:[%v]", day, err)
		//nolint:all
		return nil, nil
	}
	dailyRankMap = make(map[int64]int32)
	for k, v := range result {
		dailyRankMap[util.StrToInt64(k)] = util.StrToInt32(v)
	}
	return
}

// 获取玩家斗法每日排行
func (a *Arena) GetDailyRankByUid(serverID int64, pvpType int32, day string, uid int64) (rank int32, err error) {
	key := fmt.Sprintf(model.RedisKey_Arena_Daily_Rank, serverID, pvpType, day)
	result, err := a.client.HGet(context.TODO(), key, util.Int64ToStr(uid)).Result()
	if err != nil {
		log.Error("GetDailyRankByUid day:[%v] err:[%v]", day, err)
		return -1, nil
	}
	rank = util.StrToInt32(result)
	return
}

// 删除斗法每日排行
func (a *Arena) DelDailyRank(serverID int64, pvpType int32, day string) (err error) {
	key := fmt.Sprintf(model.RedisKey_Arena_Daily_Rank, serverID, pvpType, day)
	_, err = a.client.Del(context.TODO(), key).Result()
	if err != nil {
		log.Error("DelDailyRank err:[%v]", err)
		return
	}
	return
}

// 记录斗法成就排行
func (a *Arena) SaveAchievementRank(serverID int64, pvpType int32, achievementRankMap map[int64]int32) (err error) {
	pipe := a.client.Pipeline()
	for uid, rank := range achievementRankMap {
		key := fmt.Sprintf(model.RedisKey_Arena_Achievement_Rank, serverID, pvpType)
		_, err = pipe.HSet(context.TODO(), key, util.Int64ToStr(uid), rank).Result()
		if err != nil {
			log.Error("SaveAchievementRank HSet Error: %v", err)
		}
	}
	_, err = pipe.Exec(context.TODO())
	if err != nil {
		log.Error("SaveDailyPvpFightingRank err:[%v]")
		return
	}
	return
}

// 获取斗法成就排行
func (a *Arena) GetAchievementRank(serverID int64, pvpType int32) (achievementRankMap map[int64]int32, err error) {
	key := fmt.Sprintf(model.RedisKey_Arena_Achievement_Rank, serverID, pvpType)
	result, err := a.client.HGetAll(context.TODO(), key).Result()
	if err != nil {
		log.Error("GetDailyPvpFightingRank err:[%v]", err)
		return nil, err
	}
	achievementRankMap = make(map[int64]int32)
	for k, v := range result {
		uid := util.StrToInt64(k)
		rank := util.StrToInt32(v)
		achievementRankMap[uid] = rank
	}
	return
}

// 获取斗法单人成就排行
func (a *Arena) GetAchievementRankByUid(serverID int64, pvpType int32, uid int64) (rank int32, err error) {
	key := fmt.Sprintf(model.RedisKey_Arena_Achievement_Rank, serverID, pvpType)
	result, err := a.client.HGet(context.TODO(), key, util.Int64ToStr(uid)).Result()
	if err != nil {
		if errors.Is(nil, v8.Nil) {
			return -1, nil
		}
		log.Error("GetDailyPvpFightingRank err:[%v]", err)
		return -1, err
	}
	rank = util.StrToInt32(result)
	return
}
