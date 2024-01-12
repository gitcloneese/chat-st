package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/panjf2000/ants/v2"
	"sync"
	pbAccount "xy3-proto/account"
	pbLogin "xy3-proto/login"
)

// loginPath
// 获取登录token
var (
	loginPath      = "/xy3-%v/login/Login"
	loginPathLocal = "/login/Login"
)

func initLoginPath() {
	apiLoginPath = fmt.Sprintf(loginPath, ServerId)
}

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
	url := fmt.Sprintf("%v%v", AccountAddr, apiLoginPath)
	resp, err := HttpPost(url, bytes.NewReader(reqB), nil)
	if err != nil {
		return nil, err
	}
	loginRsp := new(pbLogin.LoginRsp)
	if err := encoding.GetCodec("json").Unmarshal(resp, loginRsp); err != nil {
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
			Error("GetLoginToken failed, accountId: %v, err:%v", req.UnionId, err)
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
	if len(GameLoginResp) > 0 {
		SetDbPlayer(GameLoginResp)
	}
}

func RunGameLoginReq() {
	RunWithLogTick("getLoginTokenReq", getLoginToken, fmt.Sprintf("%v%v", AccountAddr, apiLoginPath))
}
