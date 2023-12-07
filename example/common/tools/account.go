package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/panjf2000/ants/v2"
	"sync"
	"sync/atomic"
	"time"
	pbAccount "xy3-proto/account"
	pbPlatform "xy3-proto/platform"
)

const (
	// accountRoleListPath
	// 获取游戏角色列表,聊天服token
	accountRoleListPath      = "/xy3-cross/account/AccountRoleList"
	accountRoleListPathLocal = "/account/AccountRoleList"
	// setZoneServerPath
)

func generateAccount() string {
	y, m, d := time.Now().Date()
	h, M, s := time.Now().Clock()
	return fmt.Sprintf("%v%02v%02v-%02v%02v%02v-%v", y, m, d, h, M, s, atomic.AddInt32(acc, 1))
}

func accountRoleListRequest(info *accountPlatformLoginInfo) (*pbAccount.AccountRoleListRsp, error) {
	platformAccount := info.account
	platformLoginResp := info.loginInfo
	var err error
	reqB, err := json.Marshal(pbAccount.AccountRoleListReq{
		PlatformID:   int32(PlatformId), // 内部测试 TODO 因为不同环境是根据配置来决定的 -1: 用来测试 不经过platform平台认证 1:测试服 4: k8s集群
		SDKAccountId: platformLoginResp.Unionid,
		SdkToken:     platformLoginResp.Accesstoken,
		ChannelID:    0,
	})
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("%v%v", AccountAddr, apiAccountRoleListPath)
	bodyByte, err := HttpPost(path, bytes.NewReader(reqB), nil)
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
