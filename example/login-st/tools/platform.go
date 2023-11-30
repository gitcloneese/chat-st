package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/encoding"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
	"xy3-proto/pkg/log"
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

	reqB, err := json.Marshal(pbPlatform.GuestLoginReq{
		Deviceuid: imei,
	})
	if err != nil {
		return nil, err
	}
	defer atomic.AddInt64(&RequestCount, 1)
	resp, err := HttpClient.Post(fmt.Sprintf("%v%v", PlatformAddr, platformPath), "application/json", bytes.NewReader(reqB))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("platform login failed, status code: %v", resp.StatusCode)
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

// PreparePlatformAccount
// 准备所有账户
func PreparePlatformAccount() {
	now := time.Now()
	log.Info("===============开始准备账户信息!!!====================")
	temp1 := atomic.LoadInt64(&RequestCount)
	nums := AccountNum
	wg := new(sync.WaitGroup)
	wg.Add(nums)
	var errCount int32
	for nums > 0 {
		go func() {
			defer wg.Done()
			account := generateAccount()
			_, err := platformGuestLogin(account)
			if err != nil {
				atomic.AddInt32(&errCount, 1)
				log.Error("PlatformGuestLogin account:%v err:%v", account, err)
			}
		}()
		nums--
	}
	wg.Wait()
	latency := time.Since(now).Seconds()
	n := len(PlatformGuestLogin)
	if n > 0 {
		//总共发出的请求数
		temp2 := atomic.LoadInt64(&RequestCount)
		qs := temp2 - temp1
		log.Info("==============%v个账户信息准备完成!!! 成功：%v 失败:%v 用时:%vs 请求总数:%v QPS:%v ============== ", n, int32(AccountNum)-errCount, errCount, latency, qs, float64(qs)/latency)
	} else {
		log.Info("账户信息准备失败!!! 用时:%v s ", latency)
		panic("账户信息准备失败!!!")
	}
}
