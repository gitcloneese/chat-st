package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/encoding"
	"net/http"
	pbAccount "xy3-proto/account"
	pbLogin "xy3-proto/login"
)

// login
// 获取登录token
func login(accountId string, accountResp *pbAccount.AccountRoleListRsp) (*pbLogin.LoginRsp, error) {
	//log.Printf("正在获取登录token... account:%v", accountId)
	reqB, err := json.Marshal(pbLogin.LoginReq{
		AccountToken: accountResp.AccountToken,
		SDKAccountId: accountId,
	})
	if err != nil {
		return nil, err
	}
	resp, err := HttpClient.Post(fmt.Sprintf("%v%v", LoginAddr, apiLoginPath), "application/json", bytes.NewReader(reqB))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed, status code: %v", resp.StatusCode)
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
	//log.Printf("请求登录成功！accountId: %v player:%v", accountId, loginRsp.PlayerID)
	return loginRsp, nil
}
