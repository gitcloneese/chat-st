package arena

import (
	"context"
	"errors"
	"fmt"
	"math"
	"x-server/core/pkg/util"
	util2 "x-server/core/pkg/util"

	"x-server/core/dao/model"
	pbrank "xy3-proto/rank"

	"xy3-proto/pkg/log"

	v8 "github.com/go-redis/redis/v8"
)

// 将配置迁移到core之后配置化
var _ttpvpZoneRankCount = 2000
var _ttpvpRankCount = 2000

// AddTTPvpRobotToRank 将机器人批量插入排行榜
func (a *Arena) AddTTPvpRobotToRank(zoneID int64, seasonID int32, robots map[int64]*model.TTPvpPlayerInfo, robotScore map[int64]int64) (err error) {
	var list []*pbrank.RankData
	for _, v := range robots {
		list = append(list, &pbrank.RankData{
			UniqueID: v.ID,
			Score:    robotScore[v.ID],
			ServerID: v.ServerID,
		})
	}
	err = a.AppendTTPVPRank(zoneID, seasonID, list)
	if err != nil {
		return err
	}
	return
}

// 将玩家的积分记录进排行
func (a *Arena) updateTTRankScore(rankKey string, PlayerID int64, val float64) (newScore float64, err error) {
	member := fmt.Sprintf("%d", PlayerID)
	newScore, err = a.client.ZIncrBy(context.TODO(), rankKey, val, member).Result()
	if err != nil {
		log.Error("updateTTRankScore err:%v", err)
		return 0, err
	}
	if newScore < 0 {
		newScore = math.Abs(newScore)
		newScore, err = a.client.ZIncrBy(context.TODO(), rankKey, newScore, member).Result()
		if err != nil {
			log.Error("updateTTRankScore err:%v", err)
			return 0, err
		}
	}
	return newScore, nil
}

// AppendTTPVPRank 增加诸天pvp排行
func (a *Arena) AppendTTPVPRank(zoneID int64, seasonID int32, list []*pbrank.RankData) (err error) {
	// 查询双方排行
	pipe := a.client.Pipeline()
	for _, v := range list {
		member := &v8.Z{
			Score:  float64(v.Score),
			Member: v.UniqueID,
		}
		// 追加到 战区排行
		rankKey := fmt.Sprintf(model.RedisKey_TTPvp_Zone_Rank, zoneID, seasonID)
		pipe.ZIncrNX(context.TODO(), rankKey, member)

		// 追加到 玩家所在本服排行
		rankKey = fmt.Sprintf(model.RedisKey_TTPvp_Rank, v.ServerID, seasonID)
		pipe.ZIncrNX(context.TODO(), rankKey, member)
	}
	_, err = pipe.Exec(context.TODO())
	if err != nil {
		log.Error("AppendTTPVPRank Exec Error: %v", err)
	}
	defer pipe.Close()
	return
}

// UpdateTTPvpRank 更新诸天斗法玩家分数
func (a *Arena) UpdateTTPvpScore(zoneID int64, seasonID int32, pid1, svrid1 int64, score1 float64, pid2, svrid2 int64, score2 float64) (pid1NewScore, pid2NewScore float64, err error) {
	// 更新战区排行榜
	rankKey := fmt.Sprintf(model.RedisKey_TTPvp_Zone_Rank, zoneID, seasonID)
	_, err = a.updateTTRankScore(rankKey, pid1, score1)
	if err != nil {
		log.Error("UpdateTTPvpScore updateTTRankScore Error: %v", err)
	}
	_, err = a.updateTTRankScore(rankKey, pid2, score2)
	if err != nil {
		log.Error("UpdateTTPvpScore updateTTRankScore Error: %v", err)
	}
	// 更新玩家1所在本服排行榜
	rankKey = fmt.Sprintf(model.RedisKey_TTPvp_Rank, svrid1, seasonID)
	pid1NewScore, err = a.updateTTRankScore(rankKey, pid1, score1)
	if err != nil {
		log.Error("UpdateTTPvpScore %v", err)
		return
	}

	// 更新玩家2所在本服排行榜
	rankKey = fmt.Sprintf(model.RedisKey_TTPvp_Rank, svrid2, seasonID)
	pid2NewScore, err = a.updateTTRankScore(rankKey, pid2, score2)
	if err != nil {
		log.Error("UpdateTTPvpScore %v", err)
		return
	}
	return pid1NewScore, pid2NewScore, err
}

// GetTTPvpRankInfo 查询唯一id对应的排行信息
func (a *Arena) GetTTPvpRankInfo(zoneID int64, serverID int64, seasonID int32, uniqueID int64) (info *pbrank.RankData, err error) {
	info = &pbrank.RankData{UniqueID: uniqueID}

	var (
		rankKey string
	)
	if serverID > 0 {
		rankKey = fmt.Sprintf(model.RedisKey_TTPvp_Rank, serverID, seasonID)
	} else {
		rankKey = fmt.Sprintf(model.RedisKey_TTPvp_Zone_Rank, zoneID, seasonID)
	}
	member := fmt.Sprintf("%d", uniqueID)

	rank, err := a.client.ZRevRank(context.TODO(), rankKey, member).Result()
	if errors.Is(err, v8.Nil) { // key 不存在
		info.Rank = -1
		err = nil
		return info, err
	}
	if err != nil { // 操作失败
		info.Rank = -1
		return info, err
	}
	info.Rank = rank

	// 获取分数
	score, _ := a.client.ZScore(context.TODO(), rankKey, member).Result()
	info.Score = int64(score)
	return info, err
}

// GetTTPvpZoneRank 获取诸天斗法战区排行榜
func (a *Arena) GetTTPvpRankList(zoneID int64, serverID int64, seasonID int32, start, stop int64) (list []*pbrank.RankData, err error) {
	var (
		rankKey string
		maxLen  int64
	)
	if serverID > 0 {
		rankKey = fmt.Sprintf(model.RedisKey_TTPvp_Rank, serverID, seasonID)
		maxLen = int64(_ttpvpRankCount)
	} else {
		rankKey = fmt.Sprintf(model.RedisKey_TTPvp_Zone_Rank, zoneID, seasonID)
		maxLen = int64(_ttpvpZoneRankCount)
	}

	list = make([]*pbrank.RankData, 0)

	res := a.client.ZRevRangeWithScores(context.TODO(), rankKey, start, stop)
	ret, _ := res.Result() //
	ids := make([]int64, len(ret))
	for k, z := range ret {
		rankData := &pbrank.RankData{
			UniqueID: util2.StrToInt64(z.Member.(string)),
			Score:    int64(z.Score),
			Rank:     int64(k) + start,
		}
		list = append(list, rankData)

		ids[k] = rankData.UniqueID
	}

	// 针对长度进行裁剪
	count, _ := a.client.ZCard(context.TODO(), rankKey).Result()
	if count > maxLen {
		a.client.ZRemRangeByRank(context.TODO(), rankKey, 0, count-maxLen-1)
	}

	return list, err
}

// GetTTPvpScoreList 获取诸天斗法积分排行榜
func (a *Arena) GetTTPvpScoreList(zoneID int64, serverID int64, seasonID int32, startScore, stopScore int64) (list []*pbrank.RankData, err error) {
	var (
		rankKey string
		maxLen  int64
	)
	if serverID > 0 {
		rankKey = fmt.Sprintf(model.RedisKey_TTPvp_Rank, serverID, seasonID)
		maxLen = int64(_ttpvpRankCount)
	} else {
		rankKey = fmt.Sprintf(model.RedisKey_TTPvp_Zone_Rank, zoneID, seasonID)
		maxLen = int64(_ttpvpZoneRankCount)
	}

	list = make([]*pbrank.RankData, 0)

	opt := &v8.ZRangeBy{
		Min: util2.Int64ToStr(startScore),
		Max: util2.Int64ToStr(stopScore),
	}

	res := a.client.ZRevRangeByScoreWithScores(context.TODO(), rankKey, opt)
	ret, _ := res.Result() //
	ids := make([]int64, len(ret))
	for k, z := range ret {
		rankData := &pbrank.RankData{
			UniqueID: util2.StrToInt64(z.Member.(string)),
			Score:    int64(z.Score),
		}
		list = append(list, rankData)

		ids[k] = rankData.UniqueID
	}

	// 针对长度进行裁剪
	count, _ := a.client.ZCard(context.TODO(), rankKey).Result()
	if count > maxLen {
		a.client.ZRemRangeByRank(context.TODO(), rankKey, 0, count-maxLen-1)
	}

	return list, err
}

// GetTTPvpLastRankList 获取诸天斗法战区历史排行榜
func GetTTPvpLastRankList(redis *v8.Client, zoneID int64, serverID int64, seasonID int32, start, stop int64) (list []*pbrank.RankData, err error) {
	var (
		rankKey string
	)
	if serverID > 0 {
		rankKey = fmt.Sprintf(model.RedisKey_TTPvp_Rank, serverID, seasonID)
	} else {
		rankKey = fmt.Sprintf(model.RedisKey_TTPvp_Zone_Rank, zoneID, seasonID)
	}

	list = make([]*pbrank.RankData, 0)
	res := redis.ZRevRangeWithScores(context.TODO(), rankKey, start, stop)
	ret, _ := res.Result() //
	ids := make([]int64, len(ret))
	for k, z := range ret {
		rankData := &pbrank.RankData{
			UniqueID: util2.StrToInt64(z.Member.(string)),
			Score:    int64(z.Score),
			Rank:     int64(k) + start,
		}
		list = append(list, rankData)
		ids[k] = rankData.UniqueID
	}

	return list, err
}

// 保存 TTPVP 之前的排名
func (a *Arena) SaveTTPvpPreviousScore(zoneID int64, serverID int64, seasonID int32, uniqueID int64, previousScore int64) (err error) {
	rankKey := fmt.Sprintf(model.RedisKey_TTPvp_Zone_PreviousScore, zoneID, seasonID)
	if serverID > 0 {
		rankKey = fmt.Sprintf(model.RedisKey_TTPvp_PreviousScore, serverID, seasonID)
	}

	member := fmt.Sprintf("%d", uniqueID)
	_, err = a.client.HSet(context.TODO(), rankKey, member, previousScore).Result()
	if err != nil {
		return err
	}

	return nil
}

// 获取 TTPVP 的上一次排名
func (a *Arena) GetTTPvpPreviousScore(zoneID int64, serverID int64, seasonID int32, uniqueID int64) (previousScore int64, err error) {
	rankKey := fmt.Sprintf(model.RedisKey_TTPvp_Zone_PreviousScore, zoneID, seasonID)
	if serverID > 0 {
		rankKey = fmt.Sprintf(model.RedisKey_TTPvp_PreviousScore, serverID, seasonID)
	}

	member := fmt.Sprintf("%d", uniqueID)
	score, err := a.client.HGet(context.TODO(), rankKey, member).Result()
	if errors.Is(err, v8.Nil) { // key 不存在
		return -1, nil
	} else if err != nil {
		return -1, err
	}
	previousScore = util.StrToInt64(score)

	return previousScore, nil
}
