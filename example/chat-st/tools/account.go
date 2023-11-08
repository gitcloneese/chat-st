package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/encoding"
	"io"
	"log"
	"net/http"
	"sync/atomic"
	"time"
	pbAccount "xy3-proto/account"
	pblogin "xy3-proto/login"
)

func generateAccount() string {
	y, m, d := time.Now().Date()
	h, M, s := time.Now().Clock()
	return fmt.Sprintf("%v%02v%02v-%02v%02v%02v-%v", y, m, d, h, M, s, atomic.AddInt32(acc, 1))
}

// accountRoleList
// 获取Account认证
func accountRoleList(accountId string) (*pbAccount.AccountRoleListRsp, error) {
	log.Printf("正在获取Account认证 accountID:%v", accountId)
	reqB, err := json.Marshal(pbAccount.AccountRoleListReq{
		PlatformID:   -1, // 内部测试
		SDKAccountId: accountId,
	})
	if err != nil {
		return nil, err
	}

	resp, err := HttpClient.Post(fmt.Sprintf("%v%v", AccountAddr, apiAccountRoleListPath), "application/json", bytes.NewReader(reqB))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("accountRoleList failed, status code: %v", resp.StatusCode)
	}

	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	accountRsp := new(pbAccount.AccountRoleListRsp)

	if err := encoding.GetCodec("json").Unmarshal(bodyByte, accountRsp); err != nil {
		return nil, err
	}

	log.Printf("请求account信息成功: account:%v", accountId)
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
	// 本地环境 跳过account和login
	if isLocal {
		return &pblogin.LoginRsp{
			PlayerID: atomic.AddInt64(&localPlayerIdAcc, 1),
		}, nil
	}

	roleListRsp, err1 := accountRoleList(accountId)
	if err1 != nil {
		log.Printf("1. connectChat getAccountToken failed, accountId:%v, err:%v", accountId, err1)
		return nil, err1
	}

	loginRsp, err2 := login(accountId, roleListRsp)
	if err2 != nil {
		log.Printf("2. connectChat getLoginToken failed, accountId:%v, err:%v", accountId, err2)
		return nil, err2
	}
	return loginRsp, err
}
