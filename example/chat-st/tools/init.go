package tools

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"google.golang.org/protobuf/proto"
	"net/http"
	"sync"
	"time"
	pblogin "xy3-proto/login"
	pbchat "xy3-proto/new-chat"
	"xy3-proto/pkg/log"
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
	CmdM = map[pbchat.Operation]proto.Message{
		pbchat.Operation_OP_SendChatReply:    new(pbchat.SendChatReply),
		pbchat.Operation_OP_RecvChat:         new(pbchat.ChatMessage),
		pbchat.Operation_OP_UpdateRoomList:   new(pbchat.UpdateRoomList),
		pbchat.Operation_OP_RoomList:         new(pbchat.RoomListReq),
		pbchat.Operation_OP_RoomListReply:    new(pbchat.RoomListResp),
		pbchat.Operation_OP_RoomHistoryReply: new(pbchat.ChatHistory),
	}
)

func operationMsg(op pbchat.Operation) proto.Message {
	switch op {
	case pbchat.Operation_OP_SendChatReply:
		return new(pbchat.SendChatReply)
	case pbchat.Operation_OP_RecvChat:
		return new(pbchat.ChatMessage)
	case pbchat.Operation_OP_UpdateRoomList:
		return new(pbchat.UpdateRoomList)
	case pbchat.Operation_OP_RoomList:
		return new(pbchat.RoomListReq)
	case pbchat.Operation_OP_RoomListReply:
		return new(pbchat.RoomListResp)
	case pbchat.Operation_OP_RoomHistoryReply:
		return new(pbchat.ChatHistory)
	}
	return nil
}

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
	Addr        string
	AccountAddr string
	LoginAddr   string
	ChatAddr    string
	WsPath      string
	PlayerNum   int
	ChatCount   int

	Local            int //本地环境测试
	isLocal          bool
	localPlayerIdAcc int64 = 500000000

	PlayerTokens     map[string]*pblogin.LoginRsp
	PlayerTokensLock = new(sync.RWMutex)

	apiAccountRoleListPath = accountRoleListPath
	apiLoginPath           = loginPath
	apiSetZoneServerPath   = setZoneServerPath
	apiConnectChatPath     = wsPath
	apiSendMessagePath     = sendMessagePath

	T int // 压测类型 默认全流程 1:发送消息 2:接口消息(tcp的连接上线), 3, 4

	RequestCount       int64 // 用于做qps统计
	ErrCount           int64 // 用于做失败统计
	percentChatPlayers float64
	C                  int  // 并发携程数
	Debug              bool // debug will print log, if true will print debug information
	errCodes           = new(sync.Map)

	chatConnectCount int64 // 连上chat ws长连接数量
)

// 每个玩家默认1s发送一个聊天
func init() {
	addFlag(flag.CommandLine)

}
func addFlag(fs *flag.FlagSet) {
	fs.StringVar(&Addr, "addr", addr, fmt.Sprintf("服务器地址默认:%s", addr))
	fs.StringVar(&AccountAddr, "accountAddr", "", "账户服务地址")
	fs.StringVar(&LoginAddr, "loginAddr", "", "登录服务器地址")
	fs.StringVar(&ChatAddr, "chatAddr", "", "聊天服务器地址")
	fs.IntVar(&Local, "local", 0, "是否是本地测试 本地测试 path地址不加 xy3-xxx前缀(不访问nginx)")
	fs.StringVar(&WsPath, "wsPath", wsPath, "wsPath to connect to server")
	fs.IntVar(&PlayerNum, "playerNum", 1000, fmt.Sprintf("玩家数量默认:%v", 1000))
	fs.IntVar(&ChatCount, "chatCount", 1000, fmt.Sprintf("玩家发言次数默认:%v", 1000))
	fs.IntVar(&T, "t", 0, "压测类型 不设置默认全流程 1:发送消息 2:接口消息(tcp的连接上线), 3, 4")
	fs.Float64Var(&percentChatPlayers, "percentChatPlayers", 1, "聊天玩家百分比")
	fs.IntVar(&C, "c", 10, "并发携程数量")
	fs.BoolVar(&Debug, "debug", false, "是否开启debug")

	flag.Parse()

	if AccountAddr == "" {
		AccountAddr = Addr
	}
	if LoginAddr == "" {
		LoginAddr = Addr
	}
	if ChatAddr == "" {
		ChatAddr = Addr
	}

	if Local != 0 {
		apiAccountRoleListPath = accountRoleListPathLocal
		apiLoginPath = loginPathLocal
		apiSetZoneServerPath = setZoneServerPathLocal
		apiConnectChatPath = wsPathLocal
		apiSendMessagePath = sendMessagePathLocal
		isLocal = true
	}

	log.Info("addr:%v wsPath:%v playerNum:%v\n", Addr, apiConnectChatPath, PlayerNum)
}

func generatePlayerToken() {
	account := generateAccount()
	info, err := getChatToken(account)
	if err != nil {
		return
	}
	PlayerTokensLock.Lock()
	defer PlayerTokensLock.Unlock()
	PlayerTokens[account] = info
}

// PreparePlayers
// 准备所有玩家token信息
func preparePlayers() {
	if PlayerTokens == nil {
		PlayerTokens = make(map[string]*pblogin.LoginRsp)
	}
	playerNums := PlayerNum
	wg := new(sync.WaitGroup)
	p, _ := ants.NewPool(C)
	defer p.Release()
	for playerNums > 0 {
		playerNums--
		wg.Add(1)
		_ = p.Submit(func() {
			defer wg.Done()
			generatePlayerToken()
		})
	}
	wg.Wait()
	n := len(PlayerTokens)
	if n <= 0 {
		panic("玩家信息准备失败!!!")
	}
}

func RunPreparePlayers() {
	RunWithLogTick("preparePlayers", preparePlayers)
}
