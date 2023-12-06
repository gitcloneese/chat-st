package tools

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
	pbAccount "xy3-proto/account"
	pbLogin "xy3-proto/login"
	pbPlatform "xy3-proto/platform"
)

const (
	addr = "127.0.0.1:8000"
	// wsPath
	// 连接聊天服
	wsPath      = "/xy3-cross/new-chat/Connect"
	wsPathLocal = "/new-chat/Connect"
	// platformPath
	// 获取account登录授权
	platformPath = "/auth/platform/GuestLogin"
	// accountRoleListPath
	// 获取游戏角色列表,聊天服token
	accountRoleListPath      = "/xy3-cross/account/AccountRoleList"
	accountRoleListPathLocal = "/account/AccountRoleList"
	// loginPath
	// 获取登录token
	loginPath      = "/xy3-2/login/Login"
	loginPathLocal = "/login/Login"
	// setZoneServerPath
	// 设置角色所在区服
	setZoneServerPath      = "/xy3-cross/new-chat/SetZoneServer"
	setZoneServerPathLocal = "/new-chat/SetZoneServer"
	// sendMessagePath
	// 发送消息
	sendMessagePath      = "/xy3-cross/new-chat/SendMessage"
	sendMessagePathLocal = "/new-chat/SendMessage"
)

var (
	HttpClient = http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			MaxIdleConnsPerHost: 50000,
			MaxIdleConns:        50000,
			IdleConnTimeout:     time.Second * 10,
		},
	}
)

var (
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

)

// 每个玩家默认1s发送一个聊天
func init() {
	addFlag(flag.CommandLine)
}

func addFlag(fs *flag.FlagSet) {
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

	if isLocal {
		apiAccountRoleListPath = accountRoleListPathLocal
		apiLoginPath = loginPathLocal

		// 聊天相关地址
		apiSetZoneServerPath = setZoneServerPathLocal
		apiConnectChatPath = wsPathLocal
		apiSendMessagePath = sendMessagePathLocal
	}

	fmt.Printf("platformAddr:%v accountRoleListAddr:%v accountNum:%v C:%v\n", PlatformAddr, AccountAddr, AccountNum, C)
}
