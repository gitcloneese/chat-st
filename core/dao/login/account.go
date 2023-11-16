package login

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"x-server/core/dao/model"
	pberr "xy3-proto/errcode"
	"xy3-proto/pkg/log"
)

func (l *Login) SetAccountToken(ctx context.Context, accountId string, platformId, channelId int) (accountToken string, err error) {
	//
	outDateTime := time.Now().Add(model.AccountTokenTime * time.Second).Unix()
	accountToken = fmt.Sprintf(model.AccountTokenFmt, accountId, platformId, channelId, rand.Int31(), outDateTime)
	_, err = l.client.Set(ctx, model.GetAccountTokenKey(accountId), accountToken, model.AccountTokenTime*time.Second).Result()
	if err != nil {
		log.Error("SetAccountToken Set account id:%s err:%+v", accountId, err)
		return accountToken, err
	}
	return accountToken, err
}

func (l *Login) VerifyAccountToken(ctx context.Context, accountId string, accountToken string) (platformId, channelId int64, err error) {
	var res string
	res, err = l.client.Get(ctx, model.GetAccountTokenKey(accountId)).Result()
	if err != nil {
		log.Error("VerifyAccountToken Get account id:%s err:%+v", accountId, err)
		return platformId, channelId, err
	}
	if res != accountToken {
		log.Error("VerifyAccountToken account token not match id:%s in:%s db:%s", accountId, accountToken, res)
		err = pberr.TokenInvalid
		return platformId, channelId, err
	}
	arr := strings.Split(res, ":")
	if len(arr) != 5 {
		log.Error("VerifyAccountToken account token not match id:%s in:%s db:%s", accountId, accountToken, res)
		err = pberr.TokenInvalid
		return platformId, channelId, err
	}
	var outDateTime int64
	outDateTime, err = strconv.ParseInt(arr[4], 10, 64)
	if err != nil {
		log.Error("VerifyAccountToken account token not match id:%s in:%s db:%s", accountId, accountToken, res)
		err = pberr.TokenInvalid
		return platformId, channelId, err
	}
	now := time.Now().Unix()
	if outDateTime < now {
		log.Error("VerifyAccountToken account token date error id:%s in:%s db:%s out date time:%d", accountId, accountToken, res, now)
		err = pberr.TokenInvalid
		return platformId, channelId, err
	}
	platformId, err = strconv.ParseInt(arr[1], 10, 64)
	if err != nil {
		log.Error("VerifyAccountToken account token platform id error id:%s in:%s db:%s", accountId, accountToken, res)
		err = pberr.TokenInvalid
		return platformId, channelId, err
	}
	channelId, err = strconv.ParseInt(arr[2], 10, 64)
	if err != nil {
		log.Error("VerifyAccountToken account token channel id error id:%s in:%s db:%s", accountId, accountToken, res)
		err = pberr.TokenInvalid
		return platformId, channelId, err
	}
	return platformId, channelId, err
}
