package arena

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"x-server/core/pkg/util"

	"x-server/core/dao/model"
	pbarena "xy3-proto/arena"

	"xy3-proto/pkg/log"

	v8 "github.com/go-redis/redis/v8"
)

// ztpvpPlayerDailyTaskKey 玩家每日任务任务进度
func ztpvpDailyTaskKey(serverID int64, season int32, seasonDay int32, playerID int64) string {
	return fmt.Sprintf(model.RedisKey_ZTPvp_DailyTask, serverID, season, seasonDay, playerID)
}

// ztpvpPlayerDailyTaskRewardKey 玩家每日任务领取记录
func ztpvpDailyTaskRewardKey(serverID int64, season int32, seasonDay int32, playerID int64) string {
	return fmt.Sprintf(model.RedisKey_ZTPvp_DailyTask_Reward, serverID, season, seasonDay, playerID)
}

// ztpvpLevelTaskRewardKey 玩家段位奖励领取记录
func ztpvpLevelTaskRewardKey(serverID int64, season int32, playerID int64) string {
	return fmt.Sprintf(model.RedisKey_ZTPvp_LevelTask_Reward, serverID, season, playerID)
}

// ztpvpLevelTaskRewardStatKey 本服任务奖励领取计数
func ztpvpLevelTaskRewardCountKey(serverID int64, season int32) string {
	return fmt.Sprintf(model.RedisKey_ZTPvp_LevelTask_Count, serverID, season)
}

// 获取单个玩家基本信息
func (a *Arena) GetZTPvpPlayerInfo(seasonID int32, id int64) (info *model.ZTPvpPlayerInfo, err error) {
	key := fmt.Sprintf(model.RedisKey_ZTPvp_Player, seasonID, id)
	result, err := a.client.HGetAll(context.TODO(), key).Result()

	if errors.Is(err, v8.Nil) {
		//nolint:nilnil
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	if len(result) <= 0 {
		return info, err
	}

	info = &model.ZTPvpPlayerInfo{
		ID:             id,
		PvpLevel:       util.StrToInt32(result["PvpLevel"]),
		ChallengeCount: util.StrToInt32(result["ChallengeCount"]),
		WinCount:       util.StrToInt32(result["WinCount"]),
		BestScore:      util.StrToInt64(result["BestScore"]),
		BestPvpLevel:   util.StrToInt32(result["BestPvpLevel"]),
		ServerID:       util.StrToInt64(result["ServerID"]),
		WinningStreak:  util.StrToInt32(result["WinningStreak"]),
	}

	whiteList := make(map[int64]int32)
	err = json.Unmarshal([]byte(result["WhiteList"]), &whiteList)
	if err != nil {
		log.Error("GetZTPvpPlayerInfo unmarshal WhiteList err:[%v]", err)
		return info, err
	}
	info.WhiteList = whiteList
	return info, err
}

// 获取多个玩家基本信息
func (a *Arena) GetMutliZTPvpPlayerInfo(seasonId int32, ids []int64) (infos map[int64]*model.ZTPvpPlayerInfo, err error) {
	pipe := a.client.Pipeline()
	for _, id := range ids {
		pipe.HGetAll(context.TODO(), fmt.Sprintf(model.RedisKey_ZTPvp_Player, seasonId, id))
	}

	result, err := pipe.Exec(context.TODO())
	if err != nil {
		return nil, err
	}

	infos = make(map[int64]*model.ZTPvpPlayerInfo)
	for k, cmder := range result {
		m, err := cmder.(*v8.StringStringMapCmd).Result()
		if errors.Is(err, v8.Nil) {
			continue
		} else if err != nil {
			return nil, err
		}
		if len(m) <= 0 {
			continue
		}

		info := &model.ZTPvpPlayerInfo{
			ID:             ids[k],
			PvpLevel:       util.StrToInt32(m["PvpLevel"]),
			ChallengeCount: util.StrToInt32(m["ChallengeCount"]),
			WinCount:       util.StrToInt32(m["WinCount"]),
			BestPvpLevel:   util.StrToInt32(m["BestPvpLevel"]),
			BestScore:      util.StrToInt64(m["BestScore"]),
			ServerID:       util.StrToInt64(m["ServerID"]),
			WinningStreak:  util.StrToInt32(m["WinningStreak"]),
		}
		whiteList := make(map[int64]int32)
		err = json.Unmarshal([]byte(m["WhiteList"]), &whiteList)
		if err != nil {
			log.Error("GetZTPvpPlayerInfo unmarshal WhiteList err:[%v]", err)
			continue
		}
		info.WhiteList = whiteList
		infos[ids[k]] = info
	}

	return infos, err
}

// 获取多个玩家基本信息
func (a *Arena) UpdateZTPvpPlayerInfo(seasonID int32, info *model.ZTPvpPlayerInfo) (err error) {
	key := fmt.Sprintf(model.RedisKey_ZTPvp_Player, seasonID, info.ID)
	data := make(map[string]interface{})
	data["ID"] = info.ID
	data["PvpLevel"] = info.PvpLevel
	data["ChallengeCount"] = info.ChallengeCount
	data["WinCount"] = info.WinCount
	data["BestScore"] = info.BestScore
	data["BestPvpLevel"] = info.BestPvpLevel
	data["ServerID"] = info.ServerID
	data["WinningStreak"] = info.WinningStreak
	bytes, err := json.Marshal(info.WhiteList)
	if err != nil {
		log.Error("UpdateZTPvpPlayerInfo marshal WhiteList err:[%v]", err)
		return
	}
	data["WhiteList"] = bytes

	_, err = a.client.HSet(context.TODO(), key, data).Result()
	if err != nil {
		return
	}

	return
}

// 更新多个玩家基本信息
func (a *Arena) UpdateZTPvpMutliPlayerInfo(seasonID int32, infos map[int64]*model.ZTPvpPlayerInfo) (err error) {
	pipe := a.client.Pipeline()
	for _, info := range infos {
		key := fmt.Sprintf(model.RedisKey_ZTPvp_Player, seasonID, info.ID)
		data := make(map[string]interface{})
		data["ID"] = info.ID
		data["PvpLevel"] = info.PvpLevel
		data["ChallengeCount"] = info.ChallengeCount
		data["WinCount"] = info.WinCount
		data["BestScore"] = info.BestScore
		data["BestPvpLevel"] = info.BestPvpLevel
		data["ServerID"] = info.ServerID
		data["WinningStreak"] = info.WinningStreak
		bytes, err := json.Marshal(info.WhiteList)
		if err != nil {
			log.Error("UpdateZTPvpPlayerInfo marshal WhiteList err:[%v]", err)
			continue
		}
		data["WhiteList"] = bytes

		pipe.HSet(context.TODO(), key, data)
	}
	_, err = pipe.Exec(context.TODO())
	return
}

// 获取玩家录像列表
func (a *Arena) GetZTPvpRecord(serverID int64, id int64) (info *model.ZTPvpReport, err error) {
	result, err := a.client.HGetAll(context.TODO(), fmt.Sprintf(model.RedisKey_ZTPvp_Record, serverID, id)).Result()
	if err != nil {
		return
	}
	info = &model.ZTPvpReport{RoleID: id, FightReport: make([]*pbarena.FightReport, 0)}
	fightReport := []byte(result["FightReport"])
	if len(fightReport) > 0 {
		err = json.Unmarshal(fightReport, &info.FightReport)
		if err != nil {
			return
		}
	}
	return
}

// 保存玩家录像信息
func (a *Arena) SaveZTPvpRecord(serverID int64, id int64, info *model.ZTPvpReport) (err error) {
	key := fmt.Sprintf(model.RedisKey_ZTPvp_Record, serverID, id)
	data := make(map[string]interface{})
	data["RoleID"] = id
	reportBytes, err := json.Marshal(info.FightReport)
	if err != nil {
		return
	}
	data["FightReport"] = reportBytes
	_, err = a.client.HSet(context.TODO(), key, data).Result()
	if err != nil {
		return
	}
	return
}

// 获取巅峰对决战报
func (a *Arena) GetZTPvpBossRecord(serverID int64) (replayIDs []string, err error) {
	key := fmt.Sprintf(model.RedisKey_ZTPvp_Boss_Record, serverID)
	replayIDs, err = a.client.LRange(context.TODO(), key, 0, -1).Result()
	if err != nil {
		return
	}
	return
}

// 保存巅峰对决战报
func (a *Arena) SaveZTPvpBossRecord(serverID int64, replayID string) (err error) {
	key := fmt.Sprintf(model.RedisKey_ZTPvp_Boss_Record, serverID)
	_, err = a.client.LPush(context.TODO(), key, replayID).Result()
	if err != nil {
		return
	}

	// 巅峰对决有记录长度限制
	a.client.LTrim(context.TODO(), key, 0, 20)

	return
}

//////////////////////////////////////////////////////////////////////////////////////////////

// GetZTPvpDailyTaskInfo 查询玩家每日任务进度数据
func (a *Arena) GetZTPvpDailyTaskInfo(serverID int64, season int32, seasonDay int32, playerID int64) (info *model.ZTPvpPlayerDailyTaskInfo, err error) {
	key := ztpvpDailyTaskKey(serverID, season, seasonDay, playerID)
	val, err := a.client.HGetAll(context.TODO(), key).Result()
	if err != nil {
		return nil, err
	}
	if errors.Is(err, v8.Nil) {
		info = &model.ZTPvpPlayerDailyTaskInfo{}
		err = nil
		return
	}
	if err != nil {
		return
	}
	if len(val) <= 0 {
		info = &model.ZTPvpPlayerDailyTaskInfo{
			ChallengeCount: 0,
			WinCount:       0,
		}

		return
	}

	info = &model.ZTPvpPlayerDailyTaskInfo{
		ChallengeCount: util.StrToInt32(val["challenge"]),
		WinCount:       util.StrToInt32(val["win"]),
	}

	return
}

// AddZTPvpDailyWinCount 累加玩家每日挑战次数
func (a *Arena) AddZTPvpDailyChallengeCount(serverID int64, season, seasonDay int32, playerID int64) (ok bool, err error) {
	key := ztpvpDailyTaskKey(serverID, season, seasonDay, playerID)
	field := "challenge"
	val, err := a.client.HIncrBy(context.TODO(), key, field, 1).Result()
	if err != nil {
		return
	}
	if val <= 0 {
		return false, nil
	}
	return true, nil
}

// AddZTPvpDailyWinCount 累加玩家每日胜利次数
func (a *Arena) AddZTPvpDailyWinCount(serverID int64, season, seasonDay int32, playerID int64) (ok bool, err error) {
	key := ztpvpDailyTaskKey(serverID, season, seasonDay, playerID)
	field := "win"
	val, err := a.client.HIncrBy(context.TODO(), key, field, 1).Result()
	if err != nil {
		return
	}
	if val <= 0 {
		return false, nil
	}
	return true, nil
}

// GetZTPvpDailyTaskRewardRecord 查询玩家每日任务已领取奖励记录
func (a *Arena) GetZTPvpDailyTaskRewardRecord(serverID int64, season, seasonDay int32, playerID int64) (list []int32, err error) {
	key := ztpvpDailyTaskRewardKey(serverID, season, seasonDay, playerID)
	val, err := a.client.SMembers(context.TODO(), key).Result()
	if err != nil {
		return
	}
	list = make([]int32, len(val))
	for k, v := range val {
		list[k] = util.StrToInt32(v)
	}
	return
}

// AddZTPvpDailyTaskRewardRecord  更新每日任务领奖记录
func (a *Arena) AddZTPvpDailyTaskRewardRecord(serverID int64, season, seasonDay int32, playerID int64, taskID int32) (ok bool, err error) {
	key := ztpvpDailyTaskRewardKey(serverID, season, seasonDay, playerID)
	val, err := a.client.SAdd(context.TODO(), key, taskID).Result()
	if err != nil {
		return
	}
	if val <= 0 {
		return false, nil
	}
	return true, nil
}

// ExistZTPvpDailyTaskRewardRecord 查询玩家是否已领取指定每日任务奖励
func (a *Arena) ExistZTPvpDailyTaskRewardRecord(serverID int64, season, seasonDay int32, playerID int64, taskID int32) (ok bool, err error) {
	key := ztpvpDailyTaskRewardKey(serverID, season, seasonDay, playerID)
	return a.client.SIsMember(context.TODO(), key, taskID).Result()
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// GetZTPvpLevelTaskRewardRecord 查询玩家段位任务已领取奖励记录
func (a *Arena) GetZTPvpLevelTaskRewardRecord(serverID int64, season int32, playerID int64) (list []int32, err error) {
	key := ztpvpLevelTaskRewardKey(serverID, season, playerID)
	val, err := a.client.SMembers(context.TODO(), key).Result()
	if err != nil {
		return
	}
	list = make([]int32, len(val))
	for k, v := range val {
		list[k] = util.StrToInt32(v)
	}
	return
}

// GetMulityZTPvpLevelTaskRewardRecord 获取一批玩家的领奖记录
func (a *Arena) GetMulityZTPvpLevelTaskRewardRecord(serverID int64, season int32, playerIDs []int64) (recordMap map[int64][]int32, err error) {
	pipe := a.client.Pipeline()
	for _, playerID := range playerIDs {
		key := ztpvpLevelTaskRewardKey(serverID, season, playerID)
		pipe.SMembers(context.TODO(), key)
	}
	result, err := pipe.Exec(context.TODO())
	if err != nil {
		return
	}
	recordMap = make(map[int64][]int32)
	for k, res := range result {
		ids, err := res.(*v8.StringSliceCmd).Result()
		if err != nil {
			continue
		}
		for _, id := range ids {
			recordMap[playerIDs[k]] = append(recordMap[playerIDs[k]], util.StrToInt32(id))
		}
	}
	return
}

// AddZTPvpDailyTaskRewardRecord  更新段位已领奖记录
func (a *Arena) AddZTPvpLevelTaskRewardRecord(serverID int64, season int32, playerID int64, taskID int32) (ok bool, err error) {
	key := ztpvpLevelTaskRewardKey(serverID, season, playerID)
	field := fmt.Sprintf("%d", taskID)
	val, err := a.client.SAdd(context.TODO(), key, field).Result()
	if err != nil {
		return
	}
	if val <= 0 {
		return false, nil
	}
	return true, nil
}

// ExistZTPvpLevelTaskRewardRecord 查询玩家是否已领取指定段位任务奖励
func (a *Arena) ExistZTPvpLevelTaskRewardRecord(serverID int64, season int32, playerID int64, taskID int32) (ok bool, err error) {
	key := ztpvpLevelTaskRewardKey(serverID, season, playerID)
	return a.client.SIsMember(context.TODO(), key, taskID).Result()
}

// GetZTPvpLevelTaskRewardStat 查询段位任务累计领取奖励次数
func (a *Arena) GetZTPvpLevelTaskRewardStat(serverID int64, season int32) (infos map[int32]int64, err error) {
	key := ztpvpLevelTaskRewardCountKey(serverID, season)
	val, err := a.client.HGetAll(context.TODO(), key).Result()
	if err != nil {
		return nil, err
	}
	if errors.Is(err, v8.Nil) {
		err = nil
		return
	}
	if err != nil {
		return
	}
	infos = make(map[int32]int64)
	for k, v := range val {
		infos[util.StrToInt32(k)] = util.StrToInt64(v)
	}
	return
}

// AddZTPvpLevelTaskRewardStat 累加本服段位任务已领奖数量统计
func (a *Arena) AddZTPvpLevelTaskRewardStat(serverID int64, season int32, taskID int32) (ok bool, err error) {
	// 校验taskid是否受到限制

	key := ztpvpLevelTaskRewardCountKey(serverID, season)
	field := fmt.Sprintf("%d", taskID)
	val, err := a.client.HIncrBy(context.TODO(), key, field, 1).Result()
	if err != nil {
		return
	}
	if val <= 0 {
		return false, nil
	}
	return true, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// 天梯诸天获取段位限制
func (a *Arena) GetZTPvpLevelLimit(serverID int64, seasonID int32) map[int32][]int64 {
	key := fmt.Sprintf(model.RedisKey_ZTPvp_Level_Limit, serverID, seasonID)
	result, err := a.client.HGetAll(context.TODO(), key).Result()
	if err != nil {
		return nil
	}
	limitMap := map[int32][]int64{}
	for k, v := range result {
		var uids []int64
		err = json.Unmarshal([]byte(v), &uids)
		if err != nil {
			continue
		}
		limitMap[util.StrToInt32(k)] = uids
	}
	return limitMap
}

// 保存诸天斗法段位限制
func (a *Arena) SaveZTPvpLevelLimit(serverID int64, seasonID int32, limitMap map[int32][]int64) (err error) {
	key := fmt.Sprintf(model.RedisKey_ZTPvp_Level_Limit, serverID, seasonID)
	data := map[string]interface{}{}
	for limitID, uids := range limitMap {
		uidsBytes, err := json.Marshal(uids)
		if err != nil {
			log.Error("SaveZTPvpLevelLimit err:[%v]", err)
			return err
		}
		data[util.Int32ToStr(limitID)] = uidsBytes
	}
	_, err = a.client.HSet(context.TODO(), key, data).Result()
	return
}
