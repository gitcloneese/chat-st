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
	pbPlatform "xy3-proto/platform"
)

const (
	// platformPath
	// 获取account登录授权
	platformPath = "/auth/platform/GuestLogin"
)

var (
	acc = new(int32)
)

func generateImei() string {
	y, m, d := time.Now().Date()
	h, M, s := time.Now().Clock()
	return fmt.Sprintf("%v-%v-%v-%v-%v-%v-%v", y, m, d, h, M, s, atomic.AddInt32(acc, 1))
}

// 获取平台token
// 目前不需要获取平台token
func platformGuestLogin(imei string) (*pbPlatform.LoginResp, error) {
	if imei == "" {
		imei = generateImei()
	}
	var err error
	reqB, err := json.Marshal(pbPlatform.GuestLoginReq{
		Deviceuid: imei,
	})
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%v%v", PlatformAddr, platformPath)
	resp, err := HttpPost(url, bytes.NewReader(reqB), nil)
	loginResp := new(pbPlatform.LoginResp)
	if err := encoding.GetCodec("json").Unmarshal(resp, loginResp); err != nil {
		return nil, err
	}

	PlatformLoginLock.Lock()
	defer PlatformLoginLock.Unlock()
	PlatformGuestLogin[imei] = loginResp

	return loginResp, nil
}

// preparePlatformAccount
// 准备所有账户
func preparePlatformAccount() {
	nums := AccountNum
	wg := new(sync.WaitGroup)
	wg.Add(nums)
	p, _ := ants.NewPool(C)
	defer p.Release()
	for nums > 0 {
		err := p.Submit(func() {
			defer wg.Done()
			account := generateAccount()
			_, err := platformGuestLogin(account)
			if err != nil {
				Error("PlatformGuestLogin account:%v err:%v", account, err)
			}
		})
		if err != nil {
			print("preparePlatformAccount err: ", err)
		}
		nums--
	}
	wg.Wait()
	n := len(PlatformGuestLogin)
	if n <= 0 {
		Error("账户信息准备失败!!!")
		panic("账户信息准备失败!!!")
	}
}

func RunPlatformGuestLoginReq() {
	RunWithLogTick("platformGuestLoginReq", preparePlatformAccount, fmt.Sprintf("%v%v", PlatformAddr, platformPath))
}
