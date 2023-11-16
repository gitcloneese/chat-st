package scene

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"x-server/core/pkg/util"

	"x-server/core/dao/model"
	"xy3-proto/pkg/log"
	pbscene "xy3-proto/scene"

	v8 "github.com/go-redis/redis/v8"
)

func (s *Scene) GetSystemScore(userid int64, t pbscene.CompareType) int64 {
	key := fmt.Sprintf(model.RedisKey_SysScore, userid, int64(t))
	str, err := s.client.Get(context.TODO(), key).Result()
	if errors.Is(err, v8.Nil) {
		return 0
	}
	if err != nil {
		log.Error("GetSystemScore redis-err! userid:%v t:%v err:%v", userid, t, err)
		return 0
	}
	val, _ := strconv.ParseInt(str, 10, 64)
	return val
}

func (s *Scene) SetSystemScore(userid int64, t pbscene.CompareType, val int64) {
	key := fmt.Sprintf(model.RedisKey_SysScore, userid, int64(t))

	if _, err := s.client.Set(context.TODO(), key, fmt.Sprintf("%v", val), 0).Result(); err != nil {
		log.Error("SetSystemScore redis-err! userid:%v t:%v val:%v err:%v", userid, t, val, err)
	}
}

func (s *Scene) GetCacheEW(userid int64, ewid int32) (info *model.CacheExclusiveWeapon) {
	key := fmt.Sprintf(model.RedisKey_ExclusiveWeapon, userid, ewid)
	result, err := s.client.HGetAll(context.TODO(), key).Result()
	if err != nil {
		return nil
	}
	if len(result) == 0 {
		return nil
	}
	info = &model.CacheExclusiveWeapon{
		ID:      util.StrToInt32(result["ID"]),
		Star:    util.StrToInt32(result["Star"]),
		HoleNum: util.StrToInt32(result["HoleNum"]),
		SuitID:  util.StrToInt32(result["SuitID"]),
	}
	info.RunesFromString(result["Runes"])
	info.VBarFromString(result["VBar"])
	return
}

func (s *Scene) UpdateCacheEW(userid int64, info *model.CacheExclusiveWeapon) bool {
	key := fmt.Sprintf(model.RedisKey_ExclusiveWeapon, userid, info.ID)

	data := map[string]interface{}{
		"ID":      info.ID,
		"Star":    info.Star,
		"HoleNum": info.HoleNum,
		"SuitID":  info.SuitID,
		"Runes":   info.RunesToString(),
		"VBar":    info.VBarToString(),
	}

	if _, err := s.client.HSet(context.TODO(), key, data).Result(); err != nil {
		log.Error("updateEW err! userid:%v ewid:%v err:%v", userid, info.ID, err)
		return false
	}
	return true
}

func (s *Scene) GetDHClock(userid int64) int64 {
	key := fmt.Sprintf(model.RedisKey_TouhouClock, userid)
	str, err := s.client.Get(context.TODO(), key).Result()
	if errors.Is(err, v8.Nil) {
		return 0
	}

	if err != nil {
		log.Error("GetDHClock redis-err! userid:%v err:%v", userid, err)
		return 0
	}
	val, _ := strconv.ParseInt(str, 10, 64)
	return val
}

func (s *Scene) SetDHClock(userid int64, val int64) {
	key := fmt.Sprintf(model.RedisKey_TouhouClock, userid)

	if _, err := s.client.Set(context.TODO(), key, fmt.Sprintf("%v", val), 0).Result(); err != nil {
		log.Error("SetDHClock redis-err! userid:%v val:%v err:%v", userid, val, err)
	}
}

func (s *Scene) GetCacheLingbao(userid int64, id int32) (info *model.CacheLingbao) {
	key := fmt.Sprintf(model.RedisKey_LingBao, userid, id)

	result, err := s.client.HGetAll(context.TODO(), key).Result()
	if err != nil {
		log.Error("GetCacheLingbao redis-err! userid:%v id:%v err:%v", userid, id, err)
		return nil
	}
	if len(result) == 0 {
		return nil
	}

	info = &model.CacheLingbao{
		ID:      util.StrToInt32(result["ID"]),
		Advance: util.StrToInt32(result["Advance"]),
		Star:    util.StrToInt32(result["Star"]),
	}
	info.VBarsFromString(result["VBars"])
	return
}

func (s *Scene) UpdateCacheLingbao(userid int64, info *model.CacheLingbao) {
	key := fmt.Sprintf(model.RedisKey_LingBao, userid, info.ID)

	data := map[string]interface{}{
		"ID":      info.ID,
		"Advance": info.Advance,
		"Star":    info.Star,
		"VBars":   info.VBarsToString(),
	}

	if _, err := s.client.HSet(context.TODO(), key, data).Result(); err != nil {
		log.Error("UpdateCacheLingbao redis-err! userid:%v id:%v err:%v", userid, info.ID, err)
	}
}

func (s *Scene) GetCachePenglai(userid int64) (info *model.CachePenglai) {
	key := fmt.Sprintf(model.RedisKey_PengLai, userid)

	str, err := s.client.Get(context.TODO(), key).Result()
	if errors.Is(err, v8.Nil) {
		return nil
	}
	if err != nil {
		log.Error("GetCachePenglai redis-err! userid:%v err:%v", userid, err)
		return nil
	}
	if str == "" {
		return nil
	}

	info = &model.CachePenglai{}
	if err := json.Unmarshal([]byte(str), info); err != nil {
		log.Error("GetCachePenglai json-err! userid:%v err:%v", userid, err)
		return nil
	}
	return info
}

func (s *Scene) UpdateCachePenglai(userid int64, info *model.CachePenglai) {
	key := fmt.Sprintf(model.RedisKey_PengLai, userid)

	buf, err := json.Marshal(info)
	if err != nil {
		log.Error("UpdateCachePenglai json-err! userid:%v err:%v", userid, err)
		return
	}

	if _, err := s.client.Set(context.TODO(), key, string(buf), 0).Result(); err != nil {
		log.Error("UpdateCachePenglai redis-err! userid:%v json:%v err:%v", userid, string(buf), err)
	}
}

func (s *Scene) GetCacheBannerInfo(userid int64) (allInfo *model.CacheBannerInfoAll) {
	key := fmt.Sprintf(model.RedisKey_BannerInUse, userid)
	resultInfos, err := s.client.HGetAll(context.TODO(), key).Result()
	if errors.Is(err, v8.Nil) {
		log.Error("[SYSTEM SCORE] GetCacheBannerInfo banner not found userid:%v err:%v", userid, err)
		return
	}

	if err != nil {
		log.Error("[SYSTEM SCORE] GetCacheBannerInfo redis-err! userid:%v err:%v", userid, err)
		return nil
	}

	var bannersInUse [4]*model.CacheBannerInfo
	for key, value := range resultInfos {
		cacheBannerInfo := &model.CacheBannerInfo{}
		err = json.Unmarshal([]byte(value), &cacheBannerInfo)
		if err != nil {
			log.Error("[SYSTEM SCORE] GetCacheBannerInfo json.Unmarshal err:%v", err)
			return nil
		}
		bannersInUse[util.StrToInt(key)] = cacheBannerInfo
	}
	allInfo = &model.CacheBannerInfoAll{}
	allInfo.BannersInUse = bannersInUse
	return allInfo
}

func (s *Scene) UpdateCacheBannerInfo(userid int64, bannerIndex int32, info *model.CacheBannerInfo) {
	key := fmt.Sprintf(model.RedisKey_BannerInUse, userid)

	data, err := json.Marshal(info)
	if err != nil {
		log.Error("UpdateCacheBannerInfo json-err! userid:%v err:%v", userid, err)
		return
	}

	_, err = s.client.HSet(context.TODO(), key, bannerIndex, string(data)).Result()
	if err != nil {
		log.Error("UpdateCacheBannerInfo redis-err! userid:%v json:%v err:%v", userid, string(data), err)
		return
	}
}

func (s *Scene) GetCacheConstellation(userid int64) (info *model.CacheConstellation) {
	key := fmt.Sprintf(model.RedisKey_Constellation, userid)

	str, err := s.client.Get(context.TODO(), key).Result()
	if errors.Is(err, v8.Nil) {
		return nil
	}
	if err != nil {
		log.Error("GetCacheConstellation redis-err! userid:%v err:%v", userid, err)
		return nil
	}
	if str == "" {
		return nil
	}

	info = &model.CacheConstellation{}
	if err := json.Unmarshal([]byte(str), info); err != nil {
		log.Error("GetCacheConstellation json-err! userid:%v err:%v", userid, err)
		return nil
	}
	return info
}

func (s *Scene) UpdateCacheConstellation(userid int64, info *model.CacheConstellation) {
	key := fmt.Sprintf(model.RedisKey_Constellation, userid)

	buf, err := json.Marshal(info)
	if err != nil {
		log.Error("UpdateCacheConstellation json-err! userid:%v err:%v", userid, err)
		return
	}

	if _, err := s.client.Set(context.TODO(), key, string(buf), 0).Result(); err != nil {
		log.Error("UpdateCacheConstellation redis-err! userid:%v json:%v err:%v", userid, string(buf), err)
	}
}
