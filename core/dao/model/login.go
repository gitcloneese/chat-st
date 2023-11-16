package model

import (
	"fmt"
	"time"
)

const (
	AccessTokenKey  = "login:access-token:{%d}"
	RefreshTokenKey = "login:refresh-token:{%d}"

	/// MongoDB
	// Collection
	AccountCollection = "account"

	RedisAccount      = "account:%v"      //account:unionid set:<playerid>
	RedisRole2Account = "role2account:%v" //role2account:%v

	AllocPlayerIDKey = "prime:playerid"
	FixedPlayerID    = int64(10000000)

	// 登录设备信息
	RedisLoginInfo       = "login:loginlogs:%v"    //login:info:playerid
	RedisRegisterInfo    = "login:registerlogs:%v" //register:info:playerid
	RedisLoginInfoExpire = 24 * time.Hour
)

func GetAccessTokenKey(accID int64) string {
	return fmt.Sprintf(AccessTokenKey, accID)
}

func GetRefreshTokenKey(accID int64) string {
	return fmt.Sprintf(RefreshTokenKey, accID)
}

const (
	// 秒
	AccessTokenTime  = 7200
	RefreshTokenTime = 7200
)
