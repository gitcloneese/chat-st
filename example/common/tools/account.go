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
	pbPlatform "xy3-proto/platform"
)

func generateAccount() string {
	y, m, d := time.Now().Date()
	h, M, s := time.Now().Clock()
	return fmt.Sprintf("%v%02v%02v-%02v%02v%02v-%v", y, m, d, h, M, s, atomic.AddInt32(acc, 1))
}

// account
// 获取Account认证
func accountRoleListRequest(info *accountPlatformLoginInfo) (*pbAccount.AccountRoleListRsp, error) {
	platformAccount := info.account
	platformLoginResp := info.loginInfo
	var err error
	reqB, err := json.Marshal(pbAccount.AccountRoleListReq{
		PlatformID:   1, // 内部测试 TODO 因为不同环境是根据配置来决定的 -1: 用来测试 不经过platform平台认证 1:测试服 4: k8s集群
		SDKAccountId: platformLoginResp.Unionid,
		SdkToken:     platformLoginResp.Accesstoken,
		ChannelID:    0,
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
	resp, err := HttpClient.Post(fmt.Sprintf("%v%v", AccountAddr, apiAccountRoleListPath), "application/json", bytes.NewReader(reqB))
	if err != nil {
		return nil, err
	}
	errCodes.Store(resp.StatusCode, 1)
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("accountRoleList failed, status code: %v", resp.StatusCode)
		return nil, err
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

	AccountRoleListLock.Lock()
	defer AccountRoleListLock.Unlock()
	AccountRoleListResp[platformAccount] = accountRsp

	return accountRsp, nil
}

type accountPlatformLoginInfo struct {
	account   string
	loginInfo *pbPlatform.LoginResp
}

// accountRoleList
// 测试账号角色列表
func accountRoleList() {
	wg := &sync.WaitGroup{}
	p, _ := ants.NewPoolWithFunc(C, func(i interface{}) {
		defer wg.Done()
		platformLoginInfo := i.(*accountPlatformLoginInfo)
		_, err := accountRoleListRequest(platformLoginInfo)
		if err != nil {
			Error("accountRoleListReq account:%v, err:%v\n", platformLoginInfo, err)
		}
	})
	defer p.Release()

	for k, v := range PlatformGuestLogin {
		wg.Add(1)
		err := p.Invoke(&accountPlatformLoginInfo{k, v})
		if err != nil {
			print("accountRoleList", err)
		}
	}
	wg.Wait()

	if len(AccountRoleListResp) == 0 {
		panic("accountRoleList all error")
	}
}

func RunAccountRoleListReq() {
	RunWithLogTick("accountRoleListReq", accountRoleList, fmt.Sprintf("%v%v", AccountAddr, apiAccountRoleListPath))
}
