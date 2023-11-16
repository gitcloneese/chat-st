package login

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"x-server/core/dao/model"
	pberr "xy3-proto/errcode"

	"xy3-proto/pkg/log"
)

type tokenUser struct {
	UserId int64 `json:"userid"`
}

// 登录生成token
// token =  base64({userid:1001}).md5(userid+ts+basestr)
// 同一个accountId 每次都生成的是相同的token
func (l *Login) GetToken(ctx context.Context, playerId int64) (accessToken, refreshToken string, err error) {
	// 玩家信息 base64({userid:1001})
	userInfo := &tokenUser{
		UserId: playerId,
	}
	_data, err := json.Marshal(userInfo)
	if err != nil {
		log.Error("json Marshal err:%v", err)
		err = pberr.TokenInvalid
		return accessToken, refreshToken, err
	}
	encoded := base64.StdEncoding.EncodeToString(_data)

	// md5 md5(userid+ts+basestr)
	timeNow := time.Now().Nanosecond()
	str := fmt.Sprintf("%d_%d_accessToken", playerId, timeNow)
	h := md5.New()
	h.Write([]byte(str)) // 需要加密的字符串
	cipherStr := h.Sum(nil)
	md5Str := hex.EncodeToString(cipherStr)
	accessToken = encoded + "." + md5Str

	// 生成refresh_token
	str2 := fmt.Sprintf("%d_%d_refresh_token", playerId, timeNow)
	h.Write([]byte(str2))
	cipherStr = h.Sum(nil)
	md5Str2 := hex.EncodeToString(cipherStr)
	refreshToken = encoded + "." + md5Str2

	luaLogic := `local ok = redis.pcall("setex", KEYS[1], ARGV[2], ARGV[1])
	if type(ok) == "table" and #ok == 0 then return redis.pcall("setex", KEYS[2], ARGV[4], ARGV[3]) end 
	return ok`
	aKey := model.GetAccessTokenKey(playerId)
	rKey := model.GetRefreshTokenKey(playerId)
	accessTokenTime := int64(2 * 3600)
	refreshTokenTime := int64(2 * 3600)
	if model.AccessTokenTime > 0 {
		accessTokenTime = model.AccessTokenTime
	}
	if model.RefreshTokenTime > 0 {
		refreshTokenTime = model.RefreshTokenTime
	}

	_, err = l.client.Eval(ctx, luaLogic, []string{aKey, rKey},
		accessToken, accessTokenTime,
		refreshToken, refreshTokenTime).Result()
	if err != nil {
		log.Error("Set %v token err:%s", playerId, err.Error())
		return accessToken, refreshToken, err
	}

	return accessToken, refreshToken, err
}

func (l *Login) RefreshToken(ctx context.Context, refreshToken string) (accessToken, refreshTokenNew string, err error) {
	//活跃用户在token过期时，在用户无感知的情况下动态刷新token，做到一直在线状态
	//不活跃用户在token过期时，直接定向到登录页
	//前端检测到token过期后，携带refreshToken访问后台刷新token
	//假如refreshToken也过期了，则刷新token失败
	//假如refreshToken还在有效期，则签发新的token返回

	split := strings.Split(refreshToken, ".")
	if len(split) < 1 {
		log.Info("[TOKEN] refresh_token 错误")
		err = pberr.TokenInvalid
		return accessToken, refreshTokenNew, err
	}

	userbase64, err := base64.StdEncoding.DecodeString(split[0])
	if err != nil {
		log.Error("RefreshToken DecodeString Error: %v", err)
	}
	user := tokenUser{}
	if err = json.Unmarshal(userbase64, &user); err != nil {
		log.Info("[TOKEN] refresh_token 错误")
		err = pberr.TokenInvalid
		return accessToken, refreshTokenNew, err
	}

	playerId := user.UserId
	rKey := model.GetRefreshTokenKey(playerId)
	redisToken, err := l.client.Get(ctx, rKey).Result()
	if err != nil {
		log.Error("[TOKEN] err:%v", err)
		return accessToken, refreshTokenNew, err
	}

	if redisToken == "" || redisToken != refreshToken {
		//refresh_token都过期了
		log.Info(fmt.Sprintf("playerId:%d [refresh token] refresh_token out of date", playerId))
		err = pberr.TokenInvalid
		return accessToken, refreshTokenNew, err
	}

	accessToken, refreshTokenNew, err = l.GetToken(ctx, playerId)
	if err != nil {
		log.Error("RefreshToken GetToken Error: %v", err)
	}
	if err != nil {
		log.Error("[TOKEN] err:%v", err)
		return accessToken, refreshTokenNew, err
	}

	return accessToken, refreshTokenNew, err
}
