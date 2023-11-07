package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
	pbAccount "xy3-proto/account"
	pblogin "xy3-proto/login"
	"xy3-proto/pkg/log"
)

func generateAccount() string {
	y, m, d := time.Now().Date()
	h, M, s := time.Now().Clock()
	return fmt.Sprintf("%v%02v%02v-%02v:%02v:%02v-%v", y, m, d, h, M, s, atomic.AddInt32(acc, 1))
}

// accountRoleList
// 获取Account认证
func accountRoleList(accountId string) (*pbAccount.AccountRoleListRsp, error) {
	log.Info("正在获取Account认证 accountID:%v", accountId)
	reqB, err := json.Marshal(pbAccount.AccountRoleListReq{
		PlatformID:   -1, // 内部测试
		SDKAccountId: accountId,
	})
	if err != nil {
		return nil, err
	}
	resp, err := HttpClient.Post(fmt.Sprintf("%v%v", Addr, accountRoleListPath), "application/json", bytes.NewReader(reqB))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("accountRoleList failed, status code: %v", resp.StatusCode)
	}

	accountRsp := new(pbAccount.AccountRoleListRsp)

	buff := new(bytes.Buffer)
	from, err := buff.ReadFrom(resp.Body)
	if err != nil || from == 0 {
		return nil, err
	}
	resp.Body.Close()

	if err := json.Unmarshal(buff.Bytes(), accountRsp); err != nil {
		return nil, err
	}
	log.Info("请求account信息成功: account:%v", accountId)
	return accountRsp, nil
}

// getChatToken
// 获取聊天token
func getChatToken(account ...string) (token *pblogin.LoginRsp, err error) {
	var accountId string
	if len(account) > 0 {
		accountId = account[0]
	} else {
		accountId = generateAccount()
	}

	roleListRsp, err1 := accountRoleList(accountId)
	if err1 != nil {
		log.Error("1. connectChat getAccountToken failed, accountId:%v, err:%v", accountId, err)
		return nil, err1
	}

	loginRsp, err2 := login(accountId, roleListRsp)
	if err2 != nil {
		log.Error("2. connectChat getLoginToken failed, accountId:%v, err:%v", accountId, err)
		return nil, err2
	}
	return loginRsp, err
}
