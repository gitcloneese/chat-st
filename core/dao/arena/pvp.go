package arena

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"x-server/core/pkg/util"

	"x-server/core/dao/model"
	pbarena "xy3-proto/arena"
	"xy3-proto/pkg/log"

	v8 "github.com/go-redis/redis/v8"
)

func GetPvpFightingRecordKey(serverID int64, roleID int64) string {
	return fmt.Sprintf(model.RedisKey_PvpFighting_Record, serverPrefix(serverID), roleID)
}

func GetPvpFightingLockKey(serverID int64, pvpLevel int32) string {
	return fmt.Sprintf(model.RedisKey_PvpFighting_Lock, serverPrefix(serverID), pvpLevelPrefix(int64(pvpLevel)))
}

func (a *Arena) IsPvpInit(serverID int64) (b bool, err error) {
	b, err = a.client.SetNX(context.TODO(), fmt.Sprintf(model.RedisKey_PvpFighting_Init, serverPrefix(serverID)), "1", 0).Result()
	if err != nil {
		log.Error("IsPvpInit ...err:%v", err)
		return
	}

	a.SavePvpTime(serverID, time.Now().Unix())
	return
}

func (a *Arena) GetPvpTime(serverID int64) (time int64, err error) {
	result, err := a.client.HGet(context.TODO(), fmt.Sprintf(model.RedisKey_PvpFighting_Time, serverPrefix(serverID)), "time").Result()
	if err != nil {
		if errors.Is(nil, v8.Nil) {
			return 0, nil
		}
		return
	}
	if result != "" {
		return util.StrToInt64(result), nil
	}
	return 0, err
}

func (a *Arena) SavePvpTime(serverID int64, time int64) {
	a.client.HSet(context.TODO(), fmt.Sprintf(model.RedisKey_PvpFighting_Time, serverPrefix(serverID)), "time", time)
}

func (a *Arena) IsNeedEveryRefresh(serverID int64) (b bool, err error) {
	b, err = a.client.SetNX(context.TODO(), fmt.Sprintf(model.RefreshEveryDay, serverPrefix(serverID)), "1", 10*time.Minute).Result()
	if err != nil {
		log.Error("IsNeedEveryRefresh ... err:%v", err)
		return
	}
	return
}

func (a *Arena) GetPvpFighting(serverID int64, roleID int64) (pvpFighting *model.PvpFightingReport, err error) {
	key := GetPvpFightingRecordKey(serverID, roleID)
	result, err := a.client.HGetAll(context.TODO(), key).Result()
	if err != nil {
		return
	}
	pvpFighting = &model.PvpFightingReport{}
	fightReport := []byte(result["FightReport"])
	if len(fightReport) > 0 {
		report := make([]*pbarena.FightReport, 0)
		err = json.Unmarshal(fightReport, &report)
		if err != nil {
			return
		}
		pvpFighting.FightReport = report
	}

	return
}

func (a *Arena) SavePvpFighting(serverID int64, roleID int64, pvpFighting *model.PvpFightingReport) (err error) {
	key := GetPvpFightingRecordKey(serverID, roleID)
	data := make(map[string]interface{}, 0)
	reportBytes, err := json.Marshal(pvpFighting.FightReport)
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

// 斗法挑战锁定双方排行
func (a *Arena) PvpFightingLocked(serverID int64, pvpLevel int32, selfRank int32, targetRank int32) (islock bool, err error) {
	key := GetPvpFightingLockKey(serverID, pvpLevel)
	// 未进入排行只用锁对面
	if selfRank == -1 {
		result, err := a.client.SIsMember(context.TODO(), key, targetRank).Result()
		if err != nil || result {
			//nolint:nilerr
			return true, nil
		}
		r, err := a.client.SAdd(context.TODO(), key, targetRank).Result()
		if err != nil {
			//nolint:nilerr
			return true, nil
		}
		if r == 0 {
			return true, nil
		}
	} else {
		result1, err := a.client.SIsMember(context.TODO(), key, selfRank).Result()
		if err != nil {
			//nolint:nilerr
			return true, nil
		}
		result2, err := a.client.SIsMember(context.TODO(), key, targetRank).Result()
		if err != nil {
			//nolint:nilerr
			return true, nil
		}
		if result1 || result2 {
			return true, nil
		}
		r, err := a.client.SAdd(context.TODO(), key, selfRank, targetRank).Result()
		if err != nil {
			//nolint:nilerr
			return true, nil
		}
		if r == 0 {
			return true, nil
		}
	}
	return false, nil
}

// 斗法挑战完毕解锁双方排行
func (a *Arena) PvpFightingUnlock(serverID int64, pvpLevel int32, selfRank int32, targetRank int32) (err error) {
	key := GetPvpFightingLockKey(serverID, pvpLevel)
	_, err = a.client.SRem(context.TODO(), key, selfRank, targetRank).Result()
	return
}
