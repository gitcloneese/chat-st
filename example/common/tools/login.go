package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/panjf2000/ants/v2"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
	pbAccount "xy3-proto/account"
	pbLogin "xy3-proto/login"
)

// login
// 获取登录token
func login(info *GameAccountResp) (*pbLogin.LoginRsp, error) {
	accountId := info.UnionId
	accountResp := info.AccountResp
	var err error
	reqB, err := json.Marshal(pbLogin.LoginReq{
		AccountToken: accountResp.AccountToken,
		SDKAccountId: accountId,
	})
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			atomic.AddInt64(&ErrCount, 1)
		}
	}()

	defer atomic.AddInt64(&RequestCount, 1)
	// 设置延迟
	defer SetLatency()(time.Now())
	resp, err := HttpClient.Post(fmt.Sprintf("%v%v", AccountAddr, apiLoginPath), "application/json", bytes.NewReader(reqB))
	if err != nil {
		atomic.AddInt64(&ErrCount, 1)
		return nil, err
	}
	errCodes.Store(resp.StatusCode, 1)
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		err = fmt.Errorf("login failed, status code: %v body:%v", resp.StatusCode, string(b))
		return nil, err
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

type GameAccountResp struct {
	UnionId     string
	AccountResp *pbAccount.AccountRoleListRsp
}

// getLoginToken
// 访问单个服务的login 获得登录权限
func getLoginToken() {
	wg := &sync.WaitGroup{}
	length := len(AccountRoleListResp)
	wg.Add(length)
	p, _ := ants.NewPoolWithFunc(C, func(i interface{}) {
		defer wg.Done()
		req := i.(*GameAccountResp)
		_, err := login(req)
		if err != nil {
			Error("GetLoginToken failed, accountId: %v, accountResp:%#v err:%v", req.UnionId, req.AccountResp, err)
		}
	})
	defer p.Release()
	for k, v := range AccountRoleListResp {
		err := p.Invoke(&GameAccountResp{PlatformGuestLogin[k].Unionid, v})
		if err != nil {
			print("getLoginToken failed, accountId: %v, err:%v\n", k, err)
		}
	}
	wg.Wait()
}

func RunGameLoginReq() {
	RunWithLogTick("getLoginTokenReq", getLoginToken, fmt.Sprintf("%v%v", AccountAddr, apiLoginPath))
}
