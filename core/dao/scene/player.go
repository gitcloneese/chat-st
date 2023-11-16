package scene

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"x-server/core/dao/model"
	util2 "x-server/core/pkg/util"
	"xy3-proto/pkg/log"

	"encoding/base64"
	"time"
	"xy3-proto/pkg/conf/env"

	v8 "github.com/go-redis/redis/v8"
)

func (s *Scene) IsExistRole(userid int64) bool {
	n, _ := s.client.Exists(context.TODO(), fmt.Sprintf(model.RedisKey_Player, userid)).Result()
	return n == 1
}

func (s *Scene) GetCacheRole(userid int64) (info *model.CacheRole) {
	result, err := s.client.HGetAll(context.TODO(), fmt.Sprintf(model.RedisKey_Player, userid)).Result()
	if err != nil || len(result) == 0 {
		log.Error("GetCacheRole HGetAll RedisKey_Player err:%v", err)
		return
	}
	info = &model.CacheRole{
		ID:         util2.StrToInt64(result["ID"]),
		Nick:       result["Nick"],
		Sex:        util2.StrToInt32(result["Sex"]),
		Level:      util2.StrToInt32(result["Level"]),
		Exp:        util2.StrToInt64(result["Exp"]),
		Power:      util2.StrToInt64(result["Power"]),
		HeadID:     util2.StrToInt32(result["HeadID"]),
		FrameID:    util2.StrToInt32(result["FrameID"]),
		DrawID:     util2.StrToInt32(result["DrawID"]),
		Title:      util2.StrToInt32(result["Title"]),
		LoginTime:  util2.StrToInt64(result["LoginTime"]),
		LogoutTime: util2.StrToInt64(result["LogoutTime"]),
		ServerID:   util2.StrToInt64(result["ServerID"]),
		LoginDays:  util2.StrToInt32(result["LoginDays"]),
		DailyTime:  util2.StrToInt64(result["DailyTime"]),
		IsRobot:    util2.StrToInt32(result["IsRobot"]),
	}
	return
}

// Internal function to get the role's info by field
// 通过字段获取角色信息的内部函数
func (s *Scene) getCacheInfoByField(userid int64, field model.CacheRoleFieldEnum) string {
	result, err := s.client.HGet(context.TODO(), fmt.Sprintf(model.RedisKey_Player, userid), field.String()).Result()
	if err != nil {
		log.Error("getCacheInfoByField redis-player err! userid:%v err:%v", userid, err)
		return ""
	}
	return result
}

// Uses GetCacheInfoByField to get the role's level
// 通过GetCacheInfoByField获取角色等级
func (s *Scene) GetCacheRoleLevel(userid int64) int {
	levelString := s.getCacheInfoByField(userid, model.CacheRoleFieldLevel)
	if levelString == "" {
		log.Error("[PLAYER] GetCacheRoleLevel getCacheInfoByField empty levelString, userId:%v ", userid)
		return -1
	}
	// convert string to int
	intLevel, _ := strconv.Atoi(levelString)

	return intLevel
}

// 获取多个玩家基本信息
func (s *Scene) GetMutliPlayerInfo(ids []int64) (infos map[int64]*model.CacheRole) {
	pipe := s.client.Pipeline()
	for _, ID := range ids {
		pipe.HGetAll(context.TODO(), fmt.Sprintf(model.RedisKey_Player, ID))
	}
	result, err := pipe.Exec(context.TODO())
	if err != nil {
		log.Error("GetMutliPlayerInfo Exec err:%v", err)
		return nil
	}
	infos = make(map[int64]*model.CacheRole)
	for _, cmder := range result {
		// fmt.Println("strMap", cmder)
		var m map[string]string
		m, err = cmder.(*v8.StringStringMapCmd).Result()
		if err != nil {
			log.Error("GetMutliPlayerInfo Result err:%v", err)
			//nolint:all
			err = nil
			continue
		}

		if m["ID"] == "" {
			continue
		}

		info := &model.CacheRole{
			ID:         util2.StrToInt64(m["ID"]),
			Nick:       m["Nick"],
			Sex:        util2.StrToInt32(m["Sex"]),
			Level:      util2.StrToInt32(m["Level"]),
			Exp:        util2.StrToInt64(m["Exp"]),
			Power:      util2.StrToInt64(m["Power"]),
			HeadID:     util2.StrToInt32(m["HeadID"]),
			FrameID:    util2.StrToInt32(m["FrameID"]),
			DrawID:     util2.StrToInt32(m["DrawID"]),
			Title:      util2.StrToInt32(m["Title"]),
			LoginTime:  util2.StrToInt64(m["LoginTime"]),
			LogoutTime: util2.StrToInt64(m["LogoutTime"]),
			ServerID:   util2.StrToInt64(m["ServerID"]),
			LoginDays:  util2.StrToInt32(m["LoginDays"]),
			DailyTime:  util2.StrToInt64(m["DailyTime"]),
			IsRobot:    util2.StrToInt32(m["IsRobot"]),
			RegTime:    util2.StrToInt64(m["RegTime"]),
			LockTime:   util2.StrToInt64(m["LockTime"]),
			Os:         util2.StrToInt32(m["Os"]),
			UnionId:    util2.ToString(m["UnionId"]),
		}
		infos[info.ID] = info
	}
	return infos
}

// CacheRoleServer
// 第一次创建playerId时缓存server
func (s *Scene) CacheRoleServer(playerId, server int64) (err error) {
	pipe := s.client.Pipeline()
	key := fmt.Sprintf(model.RedisKey_Player, playerId)
	data := make(map[string]interface{})
	data["ID"] = playerId
	data["ServerID"] = server
	_, err = pipe.HSet(context.TODO(), key, data).Result()
	if err != nil {
		return err
	}
	return nil
}

// CacheRoleInfo 缓存角色数据
func (s *Scene) CacheRoleInfo(info *model.CacheRole) (err error) {
	pipe := s.client.Pipeline()
	key := fmt.Sprintf(model.RedisKey_Player, info.ID)
	data := make(map[string]interface{})
	data["ID"] = info.ID
	data["Nick"] = info.Nick
	data["Sex"] = info.Sex
	data["Level"] = info.Level
	data["Exp"] = info.Exp
	data["Power"] = info.Power
	data["HeadID"] = info.HeadID
	data["FrameID"] = info.FrameID
	data["DrawID"] = info.DrawID
	data["Title"] = info.Title
	data["LoginTime"] = info.LoginTime
	data["LogoutTime"] = info.LogoutTime
	data["ServerID"] = info.ServerID
	data["LoginDays"] = info.LoginDays
	data["DailyTime"] = info.DailyTime
	data["IsRobot"] = info.IsRobot
	if info.ClientIp != "" {
		data["ClientIp"] = info.ClientIp
	}
	_, err = pipe.HSet(context.TODO(), key, data).Result()
	if err != nil {
		return err
	}

	if info.IsRobot == 0 {
		// 将角色按照等级缓存
		key = fmt.Sprintf(model.RedisKey_Player_Level, info.Level)
		pipe.SAdd(context.TODO(), key, info.ID)

		// 将角色按照名字缓存
		key = fmt.Sprintf(model.RedisKey_Player_Name, info.Nick)
		pipe.SAdd(context.TODO(), key, info.ID)
	}
	_, err = pipe.Exec(context.TODO())

	return err
}

// SaveRegTime
// 保存注册时间
func (s *Scene) SaveRegTime(userId, time int64) error {
	err := s.client.HSet(context.Background(), fmt.Sprintf(model.RedisKey_Player, userId), "RegTime", time).Err()
	if err != nil {
		log.Error("SaveRegTime player:%v time:%v err:%v", userId, time, err)
		return err
	}
	return nil
}

// SetCacheRoleOs
// 保存登录os
func (s *Scene) SetCacheRoleOs(userId int64, os int32) error {
	err := s.client.HSet(context.Background(), fmt.Sprintf(model.RedisKey_Player, userId), "Os", os).Err()
	if err != nil {
		log.Error("SetCacheRoleOs player:%v os:%v err:%v", userId, os, err)
		return err
	}
	return nil
}

// SetCacheRoleUnionId
// 保存UnionId
func (s *Scene) SetCacheRoleUnionId(userId int64, unionId string) error {
	err := s.client.HSet(context.Background(), fmt.Sprintf(model.RedisKey_Player, userId), "UnionId", unionId).Err()
	if err != nil {
		log.Error("SetCacheRoleUnionId player:%v unionId:%v err:%v", userId, unionId, err)
		return err
	}
	return nil
}

// GetCacheRoleUnionId
// 获取账户id
func (s *Scene) GetCacheRoleUnionId(userId int64) (string, error) {
	id, err := s.client.HGet(context.Background(), fmt.Sprintf(model.RedisKey_Player, userId), "UnionId").Result()
	if err != nil {
		log.Error("GetCacheRoleOs player:%v err:%v", userId, err)
		return "", err
	}
	return id, nil
}

// GetCacheRoleOs
// 获取登录os
func (s *Scene) GetCacheRoleOs(userId int64) (int32, error) {
	osStr, err := s.client.HGet(context.Background(), fmt.Sprintf(model.RedisKey_Player, userId), "Os").Result()
	if err != nil {
		log.Error("GetCacheRoleOs player:%v err:%v", userId, err)
		return 0, err
	}
	os, err0 := strconv.Atoi(osStr)
	if err0 != nil {
		log.Error("GetCacheRoleOs Atoi player:%v osStr:%v err:%v", userId, osStr, err)
	}

	return int32(os), nil
}

// CacheMutliRoleInfo 存储一批玩家信息(给机器人用)
func (s *Scene) CacheMutliRoleInfo(infos []*model.CacheRole) (err error) {
	pipe := s.client.Pipeline()
	for _, info := range infos {
		key := fmt.Sprintf(model.RedisKey_Player, info.ID)
		data := make(map[string]interface{})
		data["ID"] = info.ID
		data["Nick"] = info.Nick
		data["Sex"] = info.Sex
		data["Level"] = info.Level
		data["Exp"] = info.Exp
		data["Power"] = info.Power
		data["HeadID"] = info.HeadID
		data["FrameID"] = info.FrameID
		data["DrawID"] = info.DrawID
		data["Title"] = info.Title
		data["LoginTime"] = info.LoginTime
		data["LogoutTime"] = info.LogoutTime
		data["ServerID"] = info.ServerID
		data["LoginDays"] = info.LoginDays
		data["DailyTime"] = info.DailyTime
		data["IsRobot"] = info.IsRobot
		_, err = pipe.HSet(context.TODO(), key, data).Result()
		if err != nil {
			log.Error("CacheMutliRoleInfo Hset err:%v", err)
			return
		}
	}
	_, err = pipe.Exec(context.TODO())
	return
}

func (s *Scene) UpdateLevel(oldLevel, newLevel int32, roleID int64) (err error) {
	pipe := s.client.Pipeline()
	oldKey := fmt.Sprintf(model.RedisKey_Player_Level, oldLevel)
	pipe.SRem(context.TODO(), oldKey, roleID)
	newKey := fmt.Sprintf(model.RedisKey_Player_Level, newLevel)
	pipe.SAdd(context.TODO(), newKey, roleID)
	_, err = pipe.Exec(context.TODO())
	return
}

func (s *Scene) UpdateName(name string, serverid int32) {
	b64name := base64.StdEncoding.EncodeToString([]byte(name))
	key := fmt.Sprintf(model.RedisNameTable, b64name)

	if _, err := s.client.SAdd(context.TODO(), key, fmt.Sprintf("%v", serverid)).Result(); err != nil {
		log.Error("UpdateName redis-err! name:%v serverid:%v err:%v", name, serverid, err)
	}
}

func (s *Scene) DeleteName(name string, serverid int32) {
	b64name := base64.StdEncoding.EncodeToString([]byte(name))
	key := fmt.Sprintf(model.RedisNameTable, b64name)

	if _, err := s.client.SRem(context.TODO(), key, fmt.Sprintf("%v", serverid)).Result(); err != nil {
		log.Error("DeleteName redis-err! name:%v serverid:%v err:%v", name, serverid, err)
	}
}

func (s *Scene) AddRoleID(serverid int32, userid int64) {
	key := fmt.Sprintf(model.RedisRoleTable, serverid)

	if _, err := s.client.SAdd(context.TODO(), key, fmt.Sprintf("%v", userid)).Result(); err != nil {
		log.Error("AddRoleID redis-err! serverid:%v userid:%v err:%v", serverid, userid, err)
	}
}

func (s *Scene) IsHaveName(name string, serverid int32) bool {
	b64name := base64.StdEncoding.EncodeToString([]byte(name))
	key := fmt.Sprintf(model.RedisNameTable, b64name)

	val, err := s.client.Exists(context.TODO(), key).Result()
	if err != nil {
		log.Error("IsHaveName redis-err0! name:%v serverid:%v err:%v", name, serverid, err)
		return true
	}
	if val == 0 { //不存在
		return false
	}

	ok, err := s.client.SIsMember(context.TODO(), key, fmt.Sprintf("%v", serverid)).Result()
	if err != nil {
		log.Error("IsHaveName redis-err1! name:%v serverid:%v err:%v", name, serverid, err)
		return true
	}

	return ok
}

func (s *Scene) GetScenePrime() (info *model.CacheScenePrime) {
	key := fmt.Sprintf(model.RedisScenePrime, env.Namespace)

	result, err := s.client.HGetAll(context.TODO(), key).Result()
	if err != nil {
		log.Error("GetScenePrime rediserr:%v", err)
		return nil
	}
	if len(result) == 0 {
		return &model.CacheScenePrime{}
	}

	info = &model.CacheScenePrime{
		NowMaxLevel: util2.StrToInt32(result["NowMaxLevel"]),
	}
	return
}

func (s *Scene) UpdateScenePrime(info *model.CacheScenePrime) {
	if info == nil {
		return
	}

	//TODO, 加分布式读写锁

	key := fmt.Sprintf(model.RedisScenePrime, env.Namespace)
	data := map[string]interface{}{
		"NowMaxLevel": info.NowMaxLevel,
	}
	if _, err := s.client.HSet(context.TODO(), key, data).Result(); err != nil {
		log.Error("UpdateScenePrime err:%v", err)
	}
}

func (s *Scene) UpdateUserHeartBeat(playerID uint64) {
	err := s.client.Set(context.TODO(), fmt.Sprintf(model.RedisUserStateKey, playerID), util2.GetTimeStamp(), 0).Err()
	if err != nil {
		log.Error("UpdateUserHeartBeat Set err:%v", err)
		return
	}
}

func (s *Scene) DelAccountUserID(sdkuuid string, userid int64) (err error) {
	key := fmt.Sprintf(model.RedisKey_Account, sdkuuid)
	_, err = s.client.SRem(context.TODO(), key, interface{}([]interface{}{userid})).Result()
	if err != nil {
		log.Error("DelAccountUserID redis err! userid:%v sdkuuid:%v err:%v", userid, sdkuuid, err)
	}
	return
}

func (s *Scene) PlayerLine(userid int64) (appid string, err error) {
	key := fmt.Sprintf(model.RedisKey_Player_Line, userid)
	appid, err = s.client.Get(context.TODO(), key).Result()
	return
}

func (s *Scene) PlayerLineId(userid int64) int {
	appid, err := s.PlayerLine(userid)
	if err != nil {
		log.Error("PlayerLineId player:%v err:%v", userid, err)
		return 0
	}
	if appid != "" {
		l := strings.Split(appid, "-")
		id, err1 := strconv.Atoi(l[len(l)-1])
		if err1 != nil {
			log.Error("PlayerLineId player:%v err:%v", userid, err)
			return id
		}
	}
	return 0
}

func (s *Scene) UpdatePlayerLine(userid int64, appid string) {
	key := fmt.Sprintf(model.RedisKey_Player_Line, userid)
	if _, err := s.client.Set(context.TODO(), key, appid, time.Minute*8).Result(); err != nil {
		log.Error("UpdatePlayerLine redis err! userid:%v line:%v err:%v", userid, appid, err)
	}
}

func (s *Scene) DeletePlayerLine(userid int64) {
	key := fmt.Sprintf(model.RedisKey_Player_Line, userid)
	if _, err := s.client.Del(context.TODO(), key).Result(); err != nil {
		log.Error("DeletePlayerLine redis err! userid:%v err:%v", userid, err)
	}
}

func UpdatePlayerOS(redis *v8.Client, userid int64, ostype int32) {
	key := fmt.Sprintf(model.RedisKey_Player, userid)
	if _, err := redis.HSet(context.TODO(), key, "TheOS", fmt.Sprintf("%v", ostype)).Result(); err != nil {
		log.Error("UpdatePlayerOS redis err! userid:%v err:%v ostype:%v", userid, err, ostype)
	}
}

func PlayerOS(redis *v8.Client, userid int64) int32 {
	key := fmt.Sprintf(model.RedisKey_Player, userid)
	str, err := redis.HGet(context.TODO(), key, "TheOS").Result()
	if err != nil {
		log.Error("PlayerOS redis err! userid:%v err:%v", userid, err)
	}
	return util2.StrToInt32(str)
}

func (s *Scene) SetUserDisable(userid int64, disableTime int64, tips string) (err error) {
	key := fmt.Sprintf(model.RedisKey_Player_Disable, userid)
	data := make(map[string]interface{})
	data["userid"] = userid
	data["disableTime"] = disableTime
	data["tips"] = tips

	if _, err = s.client.HSet(context.TODO(), key, data).Result(); err != nil {
		log.Error("SetUserDisable redis err! userid:%v err:%v distime:%v tips:%v", userid, err, disableTime, tips)
	}
	return
}

func (s *Scene) GetUserDisable(userid int64) int64 {
	key := fmt.Sprintf(model.RedisKey_Player_Disable, userid)
	result, err := s.client.HGet(context.TODO(), key, "disableTime").Result()
	if err != nil {
		if !errors.Is(err, v8.Nil) {
			log.Error("GetUserDisable redis userid:%v err:%v", userid, err)
		}
		return 0
	}
	return util2.StrToInt64(result)
}

func (s *Scene) SetForbidChat(userid int64, channel int32, unixtm int64, tips string) (err error) {
	key := fmt.Sprintf(model.RedisKey_Player_ForbidChat, userid)
	data := make(map[string]interface{})
	data["userid"] = userid
	data[fmt.Sprintf("unixtm_%v", channel)] = unixtm
	data[fmt.Sprintf("tips_%v", channel)] = tips

	if _, err = s.client.HSet(context.TODO(), key, data).Result(); err != nil {
		log.Error("ForbidChat redis err! userid:%v err:%v channel:%v unixtm:%v tips:%v", userid, err, channel, unixtm, tips)
	}
	return
}

func (s *Scene) GetForbidChat(userid int64, channel int32) (unixtm int64, tips string) {
	key := fmt.Sprintf(model.RedisKey_Player_ForbidChat, userid)

	strtm, err := s.client.HGet(context.TODO(), key, fmt.Sprintf("unixtm_%v", channel)).Result()
	if err != nil {
		log.Error("GetForbidChat redis err! userid:%v err:%v", userid, err)
		return 0, ""
	}
	unixtm = util2.StrToInt64(strtm)

	tips, _ = s.client.HGet(context.TODO(), key, fmt.Sprintf("tips_%v", channel)).Result()
	return
}

func (s *Scene) ChangeName(playerID int64, oldName, newName string) error {
	log.Debug("core dao change name method called for player: %d, oldName: %s, newName: %s\n", playerID, oldName, newName)
	ctx := context.TODO()
	res := s.client.HMGet(ctx, fmt.Sprintf(model.RedisKey_Player, playerID), "ServerID", "IsRobot").Val()
	serverIDStr, isRobotStr := res[0].(string), res[1].(string)
	serverID, _ := strconv.Atoi(serverIDStr)
	isRobot, _ := strconv.ParseBool(isRobotStr)
	if isRobot || serverID <= 0 {
		// is robot
		return fmt.Errorf("not allow to change name for player: %d due to is robot", playerID)
	}
	pl := s.client.Pipeline()
	// del old name
	b64OldName := base64.StdEncoding.EncodeToString([]byte(oldName))
	// delete name from pool  logics
	pl.SRem(ctx, fmt.Sprintf(model.RedisNameTable, b64OldName), serverID)
	pl.SRem(ctx, fmt.Sprintf(model.RedisKey_Player_Name, oldName), playerID)
	// update name pool logics,
	b64NewName := base64.StdEncoding.EncodeToString([]byte(newName))
	pl.SAdd(ctx, fmt.Sprintf(model.RedisNameTable, b64NewName), serverID)
	pl.SAdd(ctx, fmt.Sprintf(model.RedisKey_Player_Name, newName), playerID)

	// 更新redis玩家昵称
	pl.HSet(ctx, fmt.Sprintf(model.RedisKey_Player, playerID), "Nick", newName)

	if _, err := pl.Exec(ctx); err != nil {
		return fmt.Errorf("failed to changeName for player %d with oldName: %s, newName: %s, err: %w", playerID, oldName, newName, err)
	}

	return nil
}
