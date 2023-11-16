package dao

import (
	"context"
	"fmt"
	"x-server/core/dao/model"
	"x-server/core/pkg/util"
	"xy3-proto/pkg/log"

	v8 "github.com/go-redis/redis/v8"
)

// 查询角色
func (d *dao) QueryPlayer(param string) (playerIDs []int64, err error) {
	key := fmt.Sprintf(model.RedisKey_Player_Name, param)
	result, err := d.client.SRandMemberN(context.TODO(), key, 10).Result()
	if err != nil {
		return
	}
	playerIDs = make([]int64, 0)
	// 如果通过参数查询名字找不到这个玩家,则将用这个参数去查找对应的id
	if len(result) == 0 {
		id := util.StrToInt64(param)
		if id == 0 {
			return
		}
		info := d.Scene.GetCacheRole(id)
		if err != nil || info == nil {
			return
		}
		playerIDs = append(playerIDs, info.ID)
	} else {
		for _, res := range result {
			playerIDs = append(playerIDs, util.StrToInt64(res))
		}
	}
	return
}

// QueryPlayerNameKey
// param 为redis的key    param :: player_name:钟水饺
func (d *dao) QueryPlayerNameKey(redisKey string) (playerIDs []int64, err error) {
	result, err := d.client.SRandMemberN(context.TODO(), redisKey, 1).Result()
	if err != nil {
		return
	}
	playerIDs = make([]int64, 0)
	// 如果通过参数查询名字找不到这个玩家,则将用这个参数去查找对应的id
	if len(result) == 0 {
		id := util.StrToInt64(redisKey)
		if id == 0 {
			return
		}
		info := d.Scene.GetCacheRole(id)
		if err != nil || info == nil {
			return
		}
		playerIDs = append(playerIDs, info.ID)
	} else {
		for _, res := range result {
			playerIDs = append(playerIDs, util.StrToInt64(res))
		}
	}
	return
}

// MatchPlayerNames
// redis 模糊查询玩家名称
func (d *dao) MatchPlayerNames(name string) []string {
	key := fmt.Sprintf("%v*", fmt.Sprintf(model.RedisKey_Player_Name, name))
	allKeys := make([]string, 0, 2) // 最多查询10个把
	var count, execCount = 10, 0
	var cursor uint64
	var err error
	var keys []string
	for len(allKeys) < count || execCount < 10 {
		keys, cursor, err = d.client.Scan(context.Background(), cursor, key, 500).Result()
		if err != nil {
			log.Error("QueryPlayer patten:%v err:%v", name, err)
			return allKeys
		}
		allKeys = append(allKeys, keys...)
		execCount++
	}
	log.Error("MatchPlayerNames name:%v keys:%v", name, allKeys)
	// 去重
	if len(allKeys) == 0 {
		return allKeys
	}
	disKey := make(map[string]bool)
	disKey1 := make([]string, 0, 4)
	for _, v := range allKeys {
		disKey[v] = true
	}
	for k := range disKey {
		disKey1 = append(disKey1, k)
	}
	return disKey1
}

// MatchPlayerNameIds
// redis 模糊查询玩家名称 -》 id
func (d *dao) MatchPlayerNameIds(name string) []int64 {
	names := d.MatchPlayerNames(name)
	if len(names) == 0 {
		return []int64{}
	}

	users := make([]int64, 0, 10)
	for _, v := range names {
		user, err := d.QueryPlayerNameKey(v)
		if err != nil {
			log.Error("MatchPlayerNameIds patten:%v name:%v", names, v)
			continue
		}
		users = append(users, user...)
	}
	return users
}

// Role2Account
// 玩家id到账号查询
func (d *dao) Role2Account(userId []int64) map[int64]string {
	resp := make(map[int64]string)
	if len(userId) == 0 {
		return resp
	}
	cmds := make(map[int64]*v8.StringCmd)
	pipe := d.client.Pipeline()
	ctx := context.Background()
	for _, v := range userId {
		key := fmt.Sprintf(model.RedisRole2Account, v)
		cmds[v] = pipe.Get(ctx, key)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Error("Role2Account userId:%v err:%v", userId, err)
		return resp
	}
	for k, v := range cmds {
		res, err1 := v.Result()
		if err1 != nil {
			log.Error("Role2Account user:%v result err:%v", k, err1)
			continue
		}
		resp[k] = res
	}
	return resp
}
