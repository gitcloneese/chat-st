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
	pbAccount "xy3-proto/account"
	pbchat "xy3-proto/new-chat"
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
	Addr         string
	AccountAddr  string
	PlatformAddr string
	AccountNum   int

	Local            int //本地环境测试
	isLocal          bool
	localPlayerIdAcc int64 = 500000000

	AccountLoginResp = make(map[string]*pbPlatform.LoginResp)
	AccountLoginLock = new(sync.RWMutex)

	AccountRoleListResp = make(map[string]*pbAccount.AccountRoleListRsp)
	AccountRoleListLock = new(sync.RWMutex)

	apiAccountRoleListPath = accountRoleListPath

	T int // 压测类型 默认全流程 1:发送消息 2:接口消息(tcp的连接上线), 3, 4

	QAcc int64 // 用于做qps统计
	c    int   // 并发携程数
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
	fs.StringVar(&PlatformAddr, "platformAddr", "", "platform账户服务地址")
	fs.IntVar(&Local, "local", 0, "是否是本地测试 本地测试 path地址不加 xy3-xxx前缀(不访问nginx)")
	fs.IntVar(&AccountNum, "accountNum", 1000, fmt.Sprintf("账户数量默认:%v", 1000))

	flag.Parse()

	if PlatformAddr == "" {
		PlatformAddr = Addr
	}
	if AccountAddr == "" {
		AccountAddr = Addr
	}

	if Local != 0 {
		apiAccountRoleListPath = accountRoleListPathLocal
		isLocal = true
	}

	log.Infof("platformAddr:%v accountRoleListAddr:%v accountNum:%v", PlatformAddr, AccountAddr, AccountNum)
}

func generatePlayerToken(ws *sync.WaitGroup, errNum int32) {
	defer ws.Done()
	account := generateAccount()
	tokenInfo, err := platformGuestLogin(account)
	if err != nil {
		atomic.AddInt32(&errNum, 1)
		return
	}
	AccountLoginLock.Lock()
	defer AccountLoginLock.Unlock()
	AccountLoginResp[account] = tokenInfo
}

// PreparePlatformAccount
// 准备所有账户
func PreparePlatformAccount() {
	now := time.Now()
	log.Info("===============开始准备账户信息!!!====================")
	temp1 := atomic.LoadInt64(&QAcc)
	nums := AccountNum
	wg := new(sync.WaitGroup)
	wg.Add(int(nums))
	var errNum int32
	for nums > 0 {
		nums--
		go generatePlayerToken(wg, errNum)
	}
	wg.Wait()
	latency := time.Since(now).Seconds()
	n := len(AccountLoginResp)
	if n > 0 {
		//总共发出的请求数
		temp2 := atomic.LoadInt64(&QAcc)
		qs := temp2 - temp1
		log.Infof("==============%v个账户信息准备完成!!! 成功：%v 失败:%v 用时:%vs 请求总数:%v QPS:%v ============== ", n, int32(AccountNum)-errNum, errNum, latency, qs, float64(qs)/latency)
	} else {
		log.Infof("账户信息准备失败!!! 用时:%v s ", latency)
		panic("账户信息准备失败!!!")
	}
}
