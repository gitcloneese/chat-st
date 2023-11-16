package arena

import (
	"context"
	"errors"
	"fmt"
	"x-server/core/pkg/util"

	"x-server/core/dao/model"
	"xy3-proto/pkg/log"
	pbrank "xy3-proto/rank"

	"github.com/go-redis/redis/v8"
	v8 "github.com/go-redis/redis/v8"
)

var (
	_pvpRankCount int64 = 1000
)

func serverPrefix(serverID int64) string {
	return fmt.Sprintf("server_%d", serverID)
}

func pvpLevelPrefix(pvpLevel int64) string {
	return fmt.Sprintf("pvp_level_%d", pvpLevel)
}

func GetPvpRankMaxCount(group int64) int64 {
	return _pvpRankCount
}

func GetPvpRankingsKey(serverID int64, pvpLevel int64) string {
	return fmt.Sprintf(model.RedisKey_PvpFighting_Rankings, serverPrefix(serverID), pvpLevelPrefix(pvpLevel))
}

// GetPlayerPvpRank 获取唯一id对应的名次.
func (a *Arena) GetPlayerPvpRank(serverID int64, pvpLevel int64, playerID int64) int64 {
	rankingsKey := GetPvpRankingsKey(serverID, pvpLevel)
	playerIDStr := util.Int64ToStr(playerID)
	result, err := a.client.ZScore(context.TODO(), rankingsKey, playerIDStr).Result()
	switch {
	case errors.Is(err, v8.Nil):
		return -1
	case err != nil:
		return -1
	default:
		return int64(result)
	}
}

// AppendPvpRank 追加
func (a *Arena) AppendPvpRank(serverID int64, pvpLevel int64, list []*pbrank.RankData) (appendedRank int64, err error) {
	rankingsKey := GetPvpRankingsKey(serverID, pvpLevel)
	size, err := a.client.ZCard(context.TODO(), rankingsKey).Result()
	appendedRank = -1
	if size >= GetPvpRankMaxCount(pvpLevel) {
		log.Error("AppendPvpRank size over max count size:%d max count:%d", size, GetPvpRankMaxCount(pvpLevel))
		return appendedRank, err
	}

	ret, err := a.client.ZRange(context.TODO(), rankingsKey, 0, -1).Result()
	if err != nil {
		log.Error("AppendPvpRank ZRange err:%v", err)
		return appendedRank, err
	}

	ids := make(map[int64]int64)
	for k, v := range ret {
		ids[util.StrToInt64(v)] = int64(0 + k)
	}

	// 查询双方排行
	pipe := a.client.Pipeline()
	// Start from index 1 instead of 0
	rank := size + 1
	for _, v := range list {
		if _, ok := ids[v.UniqueID]; ok {
			log.Warn("AppendPvpRank ids not find:%v", v.UniqueID)
			continue
		}

		pipe.ZAdd(context.TODO(), rankingsKey, &redis.Z{
			Score:  float64(rank),
			Member: v.UniqueID,
		})
		rank++
	}
	_, err = pipe.Exec(context.TODO())
	if err != nil {
		log.Error("AppendPvpRank Exec err:%v", err)
		return appendedRank, err
	}
	log.Info("AppendPvpRank rank:%d", rank)
	// Have to -1 the current value of rank is the next empty slot, after incrementing above
	appendedRank = rank - 1
	return appendedRank, err
}

// ReplaceRank 替换双方排行
func (a *Arena) ReplacePvpRank(serverID int64, pvpLevel int64, uniqueID int64, targetUniqueID int64, isTargetRobot bool) (err error) {
	// 查询双方排行
	attackerRank := a.GetPlayerPvpRank(serverID, pvpLevel, uniqueID)
	targetRank := a.GetPlayerPvpRank(serverID, pvpLevel, targetUniqueID)
	pipe := a.client.Pipeline()
	rankingsKey := GetPvpRankingsKey(serverID, pvpLevel)

	pipe.ZAdd(context.TODO(), rankingsKey, &redis.Z{
		Score:  float64(targetRank),
		Member: uniqueID,
	})

	// 设置玩家2排行
	if attackerRank > 0 {
		// If it is a player with ranking, swap ranking
		pipe.ZAdd(context.TODO(), rankingsKey, &redis.Z{
			Score:  float64(attackerRank),
			Member: targetUniqueID,
		})
	} else if isTargetRobot {
		// If it is a robot to back of rank
		size, err := a.client.ZCard(context.TODO(), rankingsKey).Result()
		if err != nil {
			log.Error("ReplacePvpRank err:[%v]", err)
		}
		lastRankingPosition := size + 1
		pipe.ZAdd(context.TODO(), rankingsKey, &redis.Z{
			Score:  float64(lastRankingPosition),
			Member: targetUniqueID,
		})
	} else {
		// If it is an unranked player, defeated player is now unranked
		pipe.ZRem(context.TODO(), rankingsKey, targetUniqueID)
	}

	_, err = pipe.Exec(context.TODO())
	if err != nil {
		log.Error("ReplacePvpRank err:[%v]", err)
	}
	return err
}

// 玩家排行晋升
func (a *Arena) UpPvpLevel(serverID, pvpLevel, playerID int64) (err error) {
	// 删除旧分组的排行
	rankingsKey := GetPvpRankingsKey(serverID, pvpLevel)
	err = a.client.ZRem(context.TODO(), rankingsKey, playerID).Err()
	if err != nil {
		log.Error("UpPvpLevel ZRem err:%v", err)
		return
	}

	// Shift everyone's rank upwards by 1
	members, err := a.client.ZRange(context.TODO(), rankingsKey, 0, -1).Result()
	if err != nil {
		log.Error("UpPvpLevel ZRange err:%v", err)
		return
	}
	for _, member := range members {
		_, err = a.client.ZIncrBy(context.TODO(), rankingsKey, -1, member).Result()
		if err != nil {
			log.Error("UpPvpLevel ZIncrBy err:%v", err)
			return
		}
	}
	return
}

// GetPvpRankInfo 查询唯一id对应的排行信息
func (a *Arena) GetPvpRankInfo(serverID int64, pvpLevel int64, uniqueID int64) (info *pbrank.RankData, err error) {
	info = &pbrank.RankData{UniqueID: uniqueID}
	rankingsKey := GetPvpRankingsKey(serverID, pvpLevel)
	uniqueIDStr := util.Int64ToStr(uniqueID)
	result, err := a.client.ZScore(context.TODO(), rankingsKey, uniqueIDStr).Result()
	switch {
	case errors.Is(err, v8.Nil):
		info.Rank = -1
		err = nil
		return
	case err != nil:
		log.Error("GetPvpRankInfo err:[%v]", err)
		info.Rank = -1
		return
	default:
		info.Rank = int64(result)
	}

	return
}

// FullRankList 填充排行列表数据
func (a *Arena) GetPvpRankList(serverID int64, group int64, start, stop int64) (list []*pbrank.RankData, err error) {
	rankingsKey := GetPvpRankingsKey(serverID, group)
	// Because we start from index 1
	startOffset := start - 1
	stopOffset := stop - 1
	ret, err := a.client.ZRange(context.TODO(), rankingsKey, startOffset, stopOffset).Result()
	if err != nil {
		log.Error("GetPvpRankList LRange err:[%v]", err)
		return list, err
	}

	ids := make([]int64, len(ret))
	for k, v := range ret {
		rankData := &pbrank.RankData{
			UniqueID: util.StrToInt64(v),
			Score:    int64(0),
			Rank:     int64(k) + start,
		}
		list = append(list, rankData)
		ids[k] = rankData.UniqueID
	}

	// 针对长度进行裁剪
	count, err := a.client.ZCard(context.TODO(), rankingsKey).Result()
	if err != nil {
		log.Error("GetPvpRankList ZCard err:[%v]", err)
		return list, err
	}
	if count > GetPvpRankMaxCount(group) {
		// Range of elements to be removed is from int64(GetPvpRankMaxCount(group)) to infinity (-1)
		_, err := a.client.ZRemRangeByRank(context.TODO(), rankingsKey, GetPvpRankMaxCount(group), -1).Result()
		if err != nil {
			log.Error("GetPvpRankList ZRemRangeByRank err:[%v]", err)
			return list, err
		}
	}
	return list, err
}
