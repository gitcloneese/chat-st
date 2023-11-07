package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
	resp, err := HttpClient.Post(fmt.Sprintf("%v%v", Addr, platformPath), "application/json", bytes.NewReader(reqB))
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
	resp.Body.Close()

	if err := json.Unmarshal(buff.Bytes(), loginResp); err != nil {
		return nil, err
	}
	return loginResp, nil
}
