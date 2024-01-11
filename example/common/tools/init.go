package tools

import (
	"flag"
	"fmt"
	"sync"
	pbAccount "xy3-proto/account"
	pbLogin "xy3-proto/login"
	pbPlatform "xy3-proto/platform"
)

const (
	addr = "127.0.0.1:8000"
)

var (
	ServerId               int
	PlatformId             int // 平台id
	Addr                   string
	AccountAddr            string
	PlatformAddr           string
	ChatAddr               string
	WsPath                 string
	ChatCount              int
	apiSetZoneServerPath   = setZoneServerPath
	apiConnectChatPath     = wsPath
	apiSendMessagePath     = sendMessagePath
	AccountNum             int
	isLocal                bool
	PlatformGuestLogin     = make(map[string]*pbPlatform.LoginResp)
	PlatformLoginLock      = new(sync.RWMutex)
	AccountRoleListResp    = make(map[string]*pbAccount.AccountRoleListRsp)
	AccountRoleListLock    = new(sync.RWMutex)
	GameLoginResp          = make(map[string]*pbLogin.LoginRsp)
	GameLoginLock          = new(sync.RWMutex)
	apiAccountRoleListPath = accountRoleListPath
	apiLoginPath           = loginPath
	T                      int64
	Debug                  bool // debug模式将会打印error日志
	C                      int  // 并发携程数

	// TestOne 压一个玩家的所有接口， 每个接口执行N次, 为true时, accountNum 参数不生效
	TestOne   bool   //
	AccountId string // 账户 指定testOne时 有需要则设置
	N         int
	Data      string // 自定义json数据
)

// 每个玩家默认1s发送一个聊天
func init() {
	addFlag(flag.CommandLine)
}

func addFlag(fs *flag.FlagSet) {
	fs.IntVar(&ServerId, "serverId", 2, "服务器id: 1, 2, 现在默认第二服务器")
	fs.IntVar(&PlatformId, "platformId", 1, "平台id 默认是1, -1:用来测试 不经过platform平台认证  1:测试服 4: k8s集群")
	fs.StringVar(&Addr, "addr", addr, fmt.Sprintf("服务器地址默认:%s", addr))
	fs.StringVar(&AccountAddr, "accountAddr", "", "账户服务地址")
	fs.StringVar(&PlatformAddr, "platformAddr", "", "platform账户服务地址")
	fs.IntVar(&AccountNum, "accountNum", 1000, fmt.Sprintf("账户数量默认:%v", 1000))
	fs.BoolVar(&isLocal, "isLocal", false, "默认false 不是本地测试 false:不是本地测试 true:本地测试")
	fs.Int64Var(&T, "t", 0, "测试流程：默认跑全程， 1:跑PlatformLogin 2:跑AccountRoleList 3:跑GetGameLoginToken")
	fs.BoolVar(&Debug, "debug", false, "debug模式，将会打印error日志")
	fs.IntVar(&C, "c", 10, "携程数")

	// 聊天相关
	fs.StringVar(&ChatAddr, "chatAddr", "", "聊天服务器地址")
	fs.StringVar(&WsPath, "wsPath", wsPath, "wsPath to connect to server")
	fs.IntVar(&ChatCount, "chatCount", 1000, fmt.Sprintf("玩家发言次数默认:%v", 1000))

	// 压一个玩家的所有接口， 每个接口执行N次
	fs.BoolVar(&TestOne, "testOne", false, "压一个玩家的所有接口， 每个接口执行N次 需要设置 -n=xxx")
	fs.StringVar(&AccountId, "accountId", "", "指定账户名")
	fs.StringVar(&Data, "data", "", `指定请求内容:'{"xxx":"xxx"}'`)
	fs.IntVar(&N, "n", 1000, "压一个玩家的所有接口， 每个接口执行N次 需要设置 -n=xxx")

	flag.Parse()

	if PlatformAddr == "" {
		PlatformAddr = Addr
	}

	if AccountAddr == "" {
		AccountAddr = Addr
	}

	if ChatAddr == "" {
		ChatAddr = Addr
	}

	// 服务器设置相关路基设置
	{
		initLoginPath()
		initManualPath()
	}

	if isLocal {
		apiAccountRoleListPath = accountRoleListPathLocal
		apiLoginPath = loginPathLocal

		// 聊天相关地址
		apiSetZoneServerPath = setZoneServerPathLocal
		apiConnectChatPath = wsPathLocal
		apiSendMessagePath = sendMessagePathLocal
		apiLoginPath = loginPathLocal
	}
	if TestOne {
		AccountNum = 1
	}

	fmt.Printf("-testOne=bool设置, 压一个玩家的所有接口 一个接口压n次 -n=int设置, 默认1000\n")
	fmt.Printf("platformAddr:%v accountRoleListAddr:%v TestOne:%v accountNum:%v C:%v\n", PlatformAddr, AccountAddr, TestOne, AccountNum, C)
}
