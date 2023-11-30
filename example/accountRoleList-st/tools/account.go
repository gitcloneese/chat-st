package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/encoding"
	log "github.com/sirupsen/logrus"
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

// accountRoleList
// 获取Account认证
func accountRoleList(platformAccount string, account *pbPlatform.LoginResp, wg *sync.WaitGroup) (*pbAccount.AccountRoleListRsp, error) {
	if wg != nil {
		defer wg.Done()
	}
	// 统计qps
	//log.Infof("正在获取Account认证 accountID:%v", accountId)
	reqB, err := json.Marshal(pbAccount.AccountRoleListReq{
		PlatformID:   1, // 内部测试 TODO 因为不同环境是根据配置来决定的 -1: 用来测试 不经过platform平台认证 1:测试服 4: k8s集群
		SDKAccountId: account.Unionid,
		SdkToken:     account.Accesstoken,
		ChannelID:    0,
	})
	if err != nil {
		return nil, err
	}

	defer atomic.AddInt64(&RequestCount, 1)
	resp, err := HttpClient.Post(fmt.Sprintf("%v%v", AccountAddr, apiAccountRoleListPath), "application/json", bytes.NewReader(reqB))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("accountRoleList failed, status code: %v", resp.StatusCode)
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

// AccountRoleList
// 测试账号角色列表
func AccountRoleList() {
	log.Info("===============开始访问accountRoleList信息!!!====================")
	now := time.Now()
	wg := &sync.WaitGroup{}
	num := len(PlatformGuestLogin)
	wg.Add(num)
	var errNum int32
	for k, v := range PlatformGuestLogin {
		go func(accountId string, accountPlatformLoginResp *pbPlatform.LoginResp, wg *sync.WaitGroup) {
			_, err := accountRoleList(accountId, accountPlatformLoginResp, wg)
			if err != nil {
				atomic.AddInt32(&errNum, 1)
				log.Errorf("accountRoleListReq account:%v, roleListResp:%v err:%v", accountId, accountPlatformLoginResp, err)
			}
		}(k, v, wg)
	}
	wg.Wait()
	latency := time.Since(now).Seconds()
	log.Infof("============== 成功:%v 失败:%v 用时:%v 请求总数:%v QPS:%v ============== ", int32(num)-errNum, errNum, latency, num, float64(num)/latency)
}
