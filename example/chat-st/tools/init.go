package tools

import (
	"crypto/tls"
	"flag"
	"fmt"
	"google.golang.org/protobuf/proto"
	"net/http"
	"time"
	pblogin "xy3-proto/login"
	chat "xy3-proto/new-chat"
	"xy3-proto/pkg/log"
)

const (
	addr = "127.0.0.1:8000"
	// wsPath
	// 连接聊天服
	wsPath = "/xy3-cross/new-chat/Connect"
	// platformPath
	// 获取account登录授权
	platformPath = "/auth/platform/GuestLogin"
	// accountRoleListPath
	// 获取游戏角色列表,聊天服token
	accountRoleListPath = "/xy3-cross/account/AccountRoleList"
	// loginPath
	// 获取登录token
	loginPath = "/xy3-2/login/Login"
	// setZoneServerPath
	// 设置角色所在区服
	setZoneServerPath = "/xy3-cross/new-chat/SetZoneServer"
	// sendMessagePath
	// 发送消息
	sendMessage = "/xy3-cross/new-chat/SendMessage"
)

var (
	CmdM = map[chat.Operation]proto.Message{
		chat.Operation_OP_SendChatReply:    new(chat.SendChatReply),
		chat.Operation_OP_RecvChat:         new(chat.ChatMessage),
		chat.Operation_OP_UpdateRoomList:   new(chat.UpdateRoomList),
		chat.Operation_OP_RoomList:         new(chat.RoomListReq),
		chat.Operation_OP_RoomListReply:    new(chat.RoomListResp),
		chat.Operation_OP_RoomHistoryReply: new(chat.ChatHistory),
	}
)

var (
	HttpClient = http.Client{
		Timeout:   time.Second * 10,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
)

var (
	Addr      string
	WsPath    string
	PlayerNum int

	PlayerTokens map[string]*pblogin.LoginRsp
)

// 每个玩家默认1s发送一个聊天
func init() {
	flag.StringVar(&Addr, "addr", addr, fmt.Sprintf("服务器地址默认:%s", addr))
	flag.StringVar(&WsPath, "wsPath", wsPath, "wsPath to connect to server")
	flag.IntVar(&PlayerNum, "playerNum", 1000, fmt.Sprintf("玩家数量默认:%v", 1000))

	log.Info("addr:%v wsPath:%v ", Addr, WsPath)
	preparePlayers()
}

// 准备所有玩家token信息
func preparePlayers() {
	now := time.Now()
	log.Info("准备玩家信息!!!")
	if PlayerTokens == nil {
		PlayerTokens = make(map[string]*pblogin.LoginRsp)
	}
	playerNums := PlayerNum
	for playerNums > 0 {
		playerNums--
		account := generateAccount()
		info, err := getChatToken(account)
		if err != nil {
			continue
		}
		time.Sleep(time.Millisecond * 10)
		PlayerTokens[account] = info
	}
	ts := time.Since(now).Seconds()
	n := len(PlayerTokens)
	if n > 0 {
		log.Info("%v个玩家信息准备完成!!! latency:%v", n, ts)
	} else {
		panic("玩家信息准备失败!!!")
	}

}
