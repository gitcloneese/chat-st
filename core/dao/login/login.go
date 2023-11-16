package login

import (
	"context"
	"fmt"
	"math/rand"
	"x-server/core/dao/model"
	coremodel "x-server/core/model"
	"x-server/core/pkg/util"
	pberr "xy3-proto/errcode"
	"xy3-proto/pkg/log"

	v8 "github.com/go-redis/redis/v8"
)

type Login struct {
	client *v8.Client
}

func New(r *v8.Client) *Login {
	return &Login{
		client: r,
	}
}

func (l *Login) AllocPlayerID() int64 {
	key := model.AllocPlayerIDKey

	for i := 0; i < 3; i++ {
		n := int64(rand.Int31n(5) + 1)
		uuid, err := l.client.IncrBy(context.TODO(), key, n).Result()
		if err != nil {
			log.Error("AllocPlayerID IncrBy redis-err:%v", err)
			continue
		}
		if uuid < model.FixedPlayerID {
			_, err = l.client.Set(context.TODO(), key, model.FixedPlayerID, 0).Result()
			if err != nil {
				log.Error("AllocPlayerID Set redis-err:%v", err)
			}
		} else {
			return uuid
		}
	}
	return 0
}

func (l *Login) AddUserID(sdkuuid string, userid int64) (err error) {
	pipe := l.client.Pipeline()
	key := fmt.Sprintf(model.RedisAccount, sdkuuid)
	ctx := context.Background()
	if err = pipe.SAdd(ctx, key, userid).Err(); err != nil {
		log.Error("AddAccountUserID sdkuuid:%v userid:%v err:%v", sdkuuid, userid, err)
		return pberr.RedisError
	}

	role2AccountKey := fmt.Sprintf(model.RedisRole2Account, userid)
	if err = pipe.Set(ctx, role2AccountKey, sdkuuid, 0).Err(); err != nil {
		log.Error("AddAccountUserID AddRedisRole2Account sdkuuid:%v userid:%v err:%v", sdkuuid, userid, err)
		return pberr.RedisError
	}

	// 添加账户id
	if err = pipe.HSet(ctx, fmt.Sprintf(model.RedisKey_Player, userid), "UnionId", sdkuuid).Err(); err != nil {
		log.Error("AddCacheRoleUuid userid:%v uuid:%v sdkuuid:%v err:%v", userid, sdkuuid, err)
		return pberr.RedisError
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Error("AddAccountUserID AddRedisRole2Account Exec sdkuuid:%v userid:%v err:%v", sdkuuid, userid, err)
		return pberr.RedisError
	}
	return
}

func (l *Login) AddLoginDeviceInfo(playerId int64, data map[string]interface{}) (err error) {
	pipe := l.client.Pipeline()
	key := fmt.Sprintf(model.RedisLoginInfo, playerId)

	_, err = pipe.HSet(context.TODO(), key, data).Result()
	if err != nil {
		log.Error("AddLoginDeviceInfo HSet err, userid:%v err:%v", playerId, err)
		return
	}

	_, err = pipe.Expire(context.TODO(), key, model.RedisLoginInfoExpire).Result()
	if err != nil {
		log.Error("AddLoginDeviceInfo set expire err, userid:%v err:%v", playerId, err)
		return pberr.RedisError
	}

	_, err = pipe.Exec(context.TODO())
	if err != nil {
		log.Error("AddLoginDeviceInfo Exec err, userid:%v err:%v", playerId, err)
		return pberr.RedisError
	}

	return
}

func (l *Login) GetLoginInfo(playerid int64) (loginInfo *coremodel.TLogFields, err error) {
	key := fmt.Sprintf(model.RedisLoginInfo, playerid)
	result, err := l.client.HGetAll(context.TODO(), key).Result()
	if err != nil {
		log.Error("GetLoginInfo HGet err, userid:%v err:%v", playerid, err)
		return
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("GetLoginInfo len(result) == 0, userid:%v", playerid)
	}

	loginInfo = &coremodel.TLogFields{
		UnionId:   result["UnionId"],
		PlayerId:  util.StrToInt64(result["PlayerId"]),
		DeviceUID: result["DeviceUID"],
		OS:        util.StrToInt(result["OS"]),
		LoginType: util.StrToInt(result["Type"]),
		LoginTime: util.StrToInt64(result["LoginTime"]),
		Expired:   util.StrToBool(result["Expired"]),
	}
	return
}
