package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/encoding"
	"io"
	"net/http"
	"sync/atomic"
	pbAccount "xy3-proto/account"
	pbLogin "xy3-proto/login"
)

// login
// 获取登录token
func login(accountId string, accountResp *pbAccount.AccountRoleListRsp) (loginRsp *pbLogin.LoginRsp, err error) {
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
	resp, err := HttpClient.Post(fmt.Sprintf("%v%v", LoginAddr, apiLoginPath), "application/json", bytes.NewReader(reqB))
	if err != nil {
		return nil, err
	}
	errCodes.Store(resp.StatusCode, 1)
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("login failed, status code: %v body:%v", resp.StatusCode, string(b))
	}

	loginRsp = new(pbLogin.LoginRsp)

	buff := new(bytes.Buffer)
	from, err := buff.ReadFrom(resp.Body)
	if err != nil || from == 0 {
		return nil, err
	}
	defer resp.Body.Close()

	if err := encoding.GetCodec("json").Unmarshal(buff.Bytes(), loginRsp); err != nil {
		return nil, err
	}
	return loginRsp, nil
}
