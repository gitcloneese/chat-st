package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/panjf2000/ants/v2"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
	pbPlatform "xy3-proto/platform"
)

var (
	acc = new(int32)
)

func generateImei() string {
	y, m, d := time.Now().Date()
	h, M, s := time.Now().Clock()
	return fmt.Sprintf("%v-%v-%v-%-%-%v-%v", y, m, d, h, M, s, atomic.AddInt32(acc, 1))
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
	defer func() {
		if err != nil {
			atomic.AddInt64(&ErrCount, 1)
		}
	}()
	defer atomic.AddInt64(&RequestCount, 1)
	// 设置延迟
	defer SetLatency()(time.Now())
	resp, err := HttpClient.Post(fmt.Sprintf("%v%v", PlatformAddr, platformPath), "application/json", bytes.NewReader(reqB))
	if err != nil {
		return nil, err
	}
	errCodes.Store(resp.StatusCode, 1)
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("platform login failed, status code: %v", resp.StatusCode)
		return nil, err
	}

	loginResp := new(pbPlatform.LoginResp)

	buff := new(bytes.Buffer)
	from, err := buff.ReadFrom(resp.Body)
	if err != nil || from == 0 {
		return nil, err
	}
	defer resp.Body.Close()

	if err := encoding.GetCodec("json").Unmarshal(buff.Bytes(), loginResp); err != nil {
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
