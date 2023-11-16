package tools

import (
	"crypto/tls"
	"flag"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
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

	QAcc               int64 // 用于做qps统计
	percentChatPlayers float64
	c                  int // 并发写成数
)

// 配置日志输出
func logInit() {
	// 设置日志切割 rotatelogs
	filePath := "logs/"
	fileName := ""
	file := filePath + fileName
	log.SetOutput(os.Stdout)

	// 设置日志级别。低于 Debug 级别的 Trace 将不会被打印
	log.SetLevel(log.DebugLevel)

	writer, _ := rotatelogs.New(
		file+"%Y%m%d.log",
		//日志最大保存时间
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithRotationSize(3*1024*1024),
		rotatelogs.WithLinkName("log.txt"),
	)
	writeMap := lfshook.WriterMap{
		log.PanicLevel: writer,
		log.FatalLevel: writer,
		log.ErrorLevel: writer,
		log.WarnLevel:  writer,
		log.InfoLevel:  writer,
		log.DebugLevel: writer,
	}
	// 配置 lfshook
	hook := lfshook.NewHook(writeMap, &log.TextFormatter{
		// 设置日期格式
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.AddHook(hook)

}

// 每个玩家默认1s发送一个聊天
func init() {
	logInit()
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
	fs.IntVar(&c, "c", 1, "并发携程书")

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

	log.Infof("addr:%v wsPath:%v playerNum:%v", Addr, apiConnectChatPath, PlayerNum)
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
	log.Info("===============开始准备玩家信息!!!====================")
	temp1 := atomic.LoadInt64(&QAcc)
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
		//总共发出的请求数
		temp2 := atomic.LoadInt64(&QAcc)
		qs := temp2 - temp1
		log.Infof("==============%v个玩家信息准备完成!!! 用时:%vs 请求总数:%v QPS:%v ============== ", n, latency, qs, float64(qs)/latency)
	} else {
		log.Infof("玩家信息准备失败!!! 用时:%v s ", latency)
		panic("玩家信息准备失败!!!")
	}
}

// PrepareChat
// 开始聊天
func PrepareChat() {
	time.Sleep(2 * time.Second)
	switch T {
	case 0:
		log.Info(`=====================开始压测!!!================================`)
		log.Info(`=========当前压测类型t为0:压测全部(发送消息/阻塞接收消息)=================`)
		log.Info(`===============================================================`)
		TestChat()
	case 1:
		log.Info(`=====================开始压测!!!================================`)
		log.Info(`=========当前压测类型t为1:压测发送消息==============================`)
		log.Info(`===============================================================`)
		TestSendMessage()
	case 2:
		time.Sleep(time.Second * 10)
		log.Info(`=====================开始压测!!!================================`)
		log.Info(`=========当前压测类型t为2:压测接收消息 阻塞接收 =====================`)
		log.Info(`===============================================================`)
		TestReceiveMessageBlock()
	case 3:
		log.Info(`=====================开始压测!!!=====================================`)
		log.Info(`=========当前压测类型t为3:压测接收消息 **** 不 ** 阻塞接收,只为了测连接数上限****`)
		log.Info(`====================================================================`)
		TestReceiveMessageUnBlock()
	case 4:
		// 一个玩家发送消息 求成功率， 平均耗时
		log.Info(`=====================开始压测!!!=====================================`)
		log.Info(`=========当前压测类型t为4:压测一个玩家发送消息  成功率， 耗时`)
		log.Info(`====================================================================`)
		TestOneSendMessage()
	}
}
