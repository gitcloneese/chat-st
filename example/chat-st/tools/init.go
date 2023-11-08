package tools

import (
	"crypto/tls"
	"flag"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
	"sync"
	"time"
	pblogin "xy3-proto/login"
	pbchat "xy3-proto/new-chat"
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
		Timeout:   time.Second * 10,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
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
	localPlayerIdAcc int64 = 100000000

	PlayerTokens     map[string]*pblogin.LoginRsp
	PlayerTokensLock = new(sync.RWMutex)

	apiAccountRoleListPath = accountRoleListPath
	apiLoginPath           = loginPath
	apiSetZoneServerPath   = setZoneServerPath
	apiConnectChatPath     = wsPath
	apiSendMessagePath     = sendMessagePath

	T int // 压测类型 默认全流程 1:发送消息 2:接口消息(tcp的连接上线), 3, 4
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

	log.Printf("addr:%v wsPath:%v playerNum:%v", Addr, apiConnectChatPath, PlayerNum)
}

func generatePlayerToken(ws *sync.WaitGroup) {
	defer ws.Done()
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
func PreparePlayers() {
	now := time.Now()
	log.Println("===============开始准备玩家信息!!!====================")
	if PlayerTokens == nil {
		PlayerTokens = make(map[string]*pblogin.LoginRsp)
	}
	playerNums := PlayerNum
	wg := new(sync.WaitGroup)
	wg.Add(playerNums)
	for playerNums > 0 {
		playerNums--
		go generatePlayerToken(wg)
	}
	wg.Wait()
	latency := time.Since(now).Seconds()
	n := len(PlayerTokens)
	if n > 0 {
		log.Printf("==============%v个玩家信息准备完成!!! 用时:%v s==============\n", n, latency)
	} else {
		log.Printf("玩家信息准备失败!!! 用时:%v s\n", latency)
		panic("玩家信息准备失败!!!")
	}

	time.Sleep(time.Second * 2)
	log.Printf("=================开始设置区服信息!!!=====================")

	//设置区服
	now1 := time.Now()
	for k, v := range PlayerTokens {
		err1 := setZoneServer(v)
		if err1 != nil {
			log.Printf("玩家:%v设置区服失败 error:%v\n", v.PlayerID, err1)
			delete(PlayerTokens, k)
			continue
		}
	}

	latency = time.Since(now1).Seconds()
	n = len(PlayerTokens)
	if n > 0 {
		log.Printf("==============%v个玩家设置区服完成!!! 用时:%v================\n", n, latency)
	} else {
		log.Printf("玩家设置区服失败!!! 用时:%v\n", latency)
		panic("玩家信息准备失败!!!")
	}
}

// PrepareChat
// 开始聊天
func PrepareChat() {
	time.Sleep(2 * time.Second)
	switch T {
	case 0:
		log.Println(`=====================开始压测!!!================================`)
		log.Println(`=========当前压测类型t为0:压测全部(发送消息/阻塞接收消息)=================`)
		log.Println(`===============================================================`)
		TestChat()
	case 1:
		log.Println(`=====================开始压测!!!================================`)
		log.Println(`=========当前压测类型t为1:压测发送消息==============================`)
		log.Println(`===============================================================`)
		TestSendMessage()
	case 2:
		log.Println(`=====================开始压测!!!================================`)
		log.Println(`=========当前压测类型t为2:压测接收消息 阻塞接收 =====================`)
		log.Println(`===============================================================`)
		TestReceiveMessageBlock()
	case 3:
		log.Println(`=====================开始压测!!!=====================================`)
		log.Println(`=========当前压测类型t为3:压测接收消息 **** 不 ** 阻塞接收,只为了测连接数上限****`)
		log.Println(`====================================================================`)
		TestReceiveMessageUnBlock()
	}
}
