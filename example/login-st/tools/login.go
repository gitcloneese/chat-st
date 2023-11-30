package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/encoding"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
	pbAccount "xy3-proto/account"
	pbLogin "xy3-proto/login"
	"xy3-proto/pkg/log"
)

// login
// 获取登录token
func login(accountId string, accountResp *pbAccount.AccountRoleListRsp) (*pbLogin.LoginRsp, error) {
	reqB, err := json.Marshal(pbLogin.LoginReq{
		AccountToken: accountResp.AccountToken,
		SDKAccountId: accountId,
	})
	if err != nil {
		return nil, err
	}
	defer atomic.AddInt64(&RequestCount, 1)
	resp, err := HttpClient.Post(fmt.Sprintf("%v%v", AccountAddr, apiLoginPath), "application/json", bytes.NewReader(reqB))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("login failed, status code: %v body:%v", resp.StatusCode, string(b))
	}

	loginRsp := new(pbLogin.LoginRsp)

	buff := new(bytes.Buffer)
	from, err := buff.ReadFrom(resp.Body)
	if err != nil || from == 0 {
		return nil, err
	}
	defer resp.Body.Close()

	if err := encoding.GetCodec("json").Unmarshal(buff.Bytes(), loginRsp); err != nil {
		return nil, err
	}

	GameLoginLock.Lock()
	defer GameLoginLock.Unlock()
	GameLoginResp[accountId] = loginRsp

	return loginRsp, nil
}

// GetLoginToken
// 访问单个服务的login 获得登录权限
func GetLoginToken() {
	wg := &sync.WaitGroup{}
	length := len(AccountRoleListResp)
	wg.Add(length)
	var errCount int32
	now := time.Now()
	log.Info("===============开始访问GetLoginToken信息!!!====================")
	for k, v := range AccountRoleListResp {
		go func(accountId string, accountResp *pbAccount.AccountRoleListRsp, wg *sync.WaitGroup) {
			defer wg.Done()
			_, err := login(PlatformGuestLogin[accountId].Unionid, accountResp)
			if err != nil {
				atomic.AddInt32(&errCount, 1)
				log.Error("GetLoginToken failed, accountId: %v, accountResp:%+v err:%v", accountId, accountResp, err)
			}
		}(k, v, wg)
	}
	wg.Wait()
	latency := time.Since(now).Seconds()
	log.Info("============== 成功:%v 失败:%v 用时:%v 请求总数:%v QPS:%v ============== ", int32(length)-errCount, errCount, latency, length, float64(length)/latency)
}
