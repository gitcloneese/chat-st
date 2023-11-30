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
	pbAccount "xy3-proto/account"
	pbLogin "xy3-proto/login"
)

// login
// 获取登录token
func login(accountId string, accountResp *pbAccount.AccountRoleListRsp) (*pbLogin.LoginRsp, error) {
	var err error
	defer func() {
		if err != nil {
			atomic.AddInt64(&ErrCount, 1)
		}
	}()
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
		atomic.AddInt64(&ErrCount, 1)
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

// getLoginToken
// 访问单个服务的login 获得登录权限
func getLoginToken() {
	wg := &sync.WaitGroup{}
	length := len(AccountRoleListResp)
	wg.Add(length)
	for k, v := range AccountRoleListResp {
		go func(accountId string, accountResp *pbAccount.AccountRoleListRsp, wg *sync.WaitGroup) {
			defer wg.Done()
			_, err := login(PlatformGuestLogin[accountId].Unionid, accountResp)
			if err != nil {
				Error("GetLoginToken failed, accountId: %v, accountResp:%+v err:%v", accountId, accountResp, err)
			}
		}(k, v, wg)
	}
	wg.Wait()
}

func RunGameLogin() {
	RunWithLog("getLoginToken", getLoginToken)
}
