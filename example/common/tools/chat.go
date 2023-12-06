package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/gorilla/websocket"
	"github.com/panjf2000/ants/v2"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
	pblogin "xy3-proto/login"
	pbchat "xy3-proto/new-chat"
	"xy3-proto/pkg/log"
)

var (
	chatConnectCount int64 // 连上chat ws长连接数量
	receiveMsgCount  int64 // 系统总收到消息数量 playerNum * perPlayerReceived
)

func connectCount() int64 {
	return atomic.LoadInt64(&chatConnectCount)
}

func msgCount() int64 {
	return atomic.LoadInt64(&receiveMsgCount)
}

func RunChat() {
	// 异步发送消息
	go RunSendMessage()
	// 接收消息
	go RunTestReceiveMessage()
}

// TestPlayerSendMessage
// 单个玩家发送消息
func TestPlayerSendMessage(info *pblogin.LoginRsp) {
	// 发送消息
	wg := &sync.WaitGroup{}
	p, _ := ants.NewPoolWithFunc(C, func(i interface{}) {
		defer wg.Done()
		err := sendMessage(info, i.(int32))
		if err != nil {
			Error("sendMessage:%v", err)
			return
		}

	})
	defer p.Release()
	var count int32
	for count < int32(ChatCount) {
		count++
		wg.Add(1)
		_ = p.Invoke(count)
	}
	wg.Wait()
}

func RunSendMessage() {
	RunWithLogTick("sendMessage", TestSendMessage)
}

// TestSendMessage
// 压测发送消息
func TestSendMessage() {
	if len(GameLoginResp) < 1 {
		panic("GameLoginResp must be at least 1")
	}
	var loginInfo *pblogin.LoginRsp
	for _, token := range GameLoginResp {
		loginInfo = token
		break
	}
	TestPlayerSendMessage(loginInfo)
}

func RunTestReceiveMessage() {
	RunReceiveMsgWithLogTick("receiveMessageBlock", TestReceiveMessageBlock, fmt.Sprintf("%v%v", ChatAddr, apiConnectChatPath))
}

// TestReceiveMessageBlock
// 压测接收消息 阻塞接收
func TestReceiveMessageBlock() {
	wg := new(sync.WaitGroup)
	p, _ := ants.NewPoolWithFunc(C, func(i interface{}) {
		defer wg.Done()
		receiveMsg(i.(*pblogin.LoginRsp))
	})
	defer p.Release()
	for _, v := range GameLoginResp {
		wg.Add(1)
		err := p.Invoke(v)
		if err != nil {
			Error("TestReceiveMessageBlock err:%v", err)
		}
	}
	wg.Wait()
}

func bearToken(token string) string {
	return fmt.Sprintf("Bearer %v", token)
}

// RunSetZoneServer
// 设置聊天频道
func RunSetZoneServer() {
	RunWithLogTick("chatSetZoneServer", SetZoneServer, fmt.Sprintf("%v%v", ChatAddr, apiSetZoneServerPath))
}

func SetZoneServer() {
	if len(GameLoginResp) < 1 {
		panic("GameLoginResp must be at least 1")
	}
	wg := new(sync.WaitGroup)
	wg.Add(len(GameLoginResp))
	p, _ := ants.NewPoolWithFunc(C, func(i interface{}) {
		err := setZoneServer(i.(*pblogin.LoginRsp))
		if err != nil {
			Error("setZoneServer err:%v", err)
		}
		wg.Done()
	})
	defer p.Release()
	for _, v := range GameLoginResp {
		err := p.Invoke(v)
		if err != nil {
			Error("setZoneServer err:%v", err)
		}
	}
	wg.Wait()
}

// 设置区服
func setZoneServer(info *pblogin.LoginRsp) (err error) {
	func() {
		if err != nil {
			atomic.AddInt64(&ErrCount, 1)
		}
	}()
	reqB, err := json.Marshal(pbchat.SetZoneServerReq{
		ZoneId:   "1", // 压力测试
		ServerId: "2",
		PlayerId: info.PlayerID,
	})
	if err != nil {
		return err
	}
	defer atomic.AddInt64(&RequestCount, 1)
	req, err := http.NewRequest("POST", fmt.Sprintf("%v%v", ChatAddr, apiSetZoneServerPath), bytes.NewReader(reqB))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearToken(info.PlayerToken))

	// 本地test
	if isLocal {
		req.Header.Set("userid", fmt.Sprintf("%v", info.PlayerID))
	}

	resp, err := HttpClient.Do(req)
	if err != nil {
		return err
	}
	errCodes.Store(resp.StatusCode, 1)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("setZoneServer failed, status code: %v", resp.StatusCode)
	}
	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	setZoneServerRsp := new(pbchat.SetZoneServerReply)

	if err := encoding.GetCodec("json").Unmarshal(bodyByte, setZoneServerRsp); err != nil {
		return err
	}
	return nil
}

func generateMessage(playerId int64, accId int32) string {
	return fmt.Sprintf("我是%v, 这是第%v次发言。", playerId, accId)
}

// 发送消息
func sendMessage(info *pblogin.LoginRsp, chatNums int32) (err error) {
	defer func() {
		if err != nil {
			atomic.AddInt64(&ErrCount, 1)
		}
	}()
	msg := generateMessage(info.PlayerID, chatNums)
	reqB, err := json.Marshal(pbchat.SendChat{
		Msg:        msg,
		FromPlayer: info.PlayerID,
		RoomType:   pbchat.RoomType_RT_World,
	})
	if err != nil {
		return err
	}
	defer atomic.AddInt64(&RequestCount, 1)
	req, err := http.NewRequest("POST", fmt.Sprintf("%v%v", ChatAddr, apiSendMessagePath), bytes.NewReader(reqB))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearToken(info.PlayerToken))

	// 本地test
	if isLocal {
		req.Header.Set("userid", fmt.Sprintf("%v", info.PlayerID))
	}

	resp, err := HttpClient.Do(req)
	if err != nil {
		return err
	}
	errCodes.Store(resp.StatusCode, 1)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("sendMessage failed, status code: %v", resp.StatusCode)
	}

	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	msgRsp := new(pbchat.SendChatReply)

	if err := encoding.GetCodec("json").Unmarshal(bodyByte, msgRsp); err != nil {
		return err
	}

	return nil
}

var (
	chatUrl     *url.URL
	chatUrlOnce sync.Once
)

// 接收消息
func receiveMsg(info *pblogin.LoginRsp) {
	if chatUrl == nil {
		chatUrlOnce.Do(func() {
			uri, err := url.Parse(ChatAddr)
			if err != nil {
				panic("receiveMsg parse url failed" + err.Error())
			}
			scheme := "ws"
			if uri.Scheme == "https" {
				scheme = "wss"
			}
			chatUrl = &url.URL{
				Scheme: scheme,
				Host:   uri.Host,
				Path:   apiConnectChatPath,
			}
		})
	}
	if chatUrl == nil {
		Error("receiveMsg chatUrl is nil")
		return
	}
	u := &url.URL{
		Scheme: chatUrl.Scheme,
		Host:   chatUrl.Host,
		Path:   chatUrl.Path,
	}

	reqHeader := http.Header{}
	if isLocal {
		u.RawQuery = fmt.Sprintf("userid=%v", info.PlayerID)
	} else {
		reqHeader.Add("Authorization", bearToken(info.PlayerToken))
	}

	c, _, err1 := websocket.DefaultDialer.Dial(u.String(), reqHeader)
	atomic.AddInt64(&RequestCount, 1)
	if err1 != nil {
		Error("receiveMsg websocket Dail err:%v", err1)
		atomic.AddInt64(&ErrCount, 1)
		return
	}

	go func(c *websocket.Conn) {
		defer c.Close()
		// 记录长练级数量
		atomic.AddInt64(&chatConnectCount, 1)
		defer atomic.AddInt64(&chatConnectCount, -1)

		// 加一个ping
		stopChan := make(chan struct{}, 1)
		defer close(stopChan)
		go func(stop chan struct{}) {
			buf := new(bytes.Buffer)
			buf.WriteByte(byte(3))
			tick := time.NewTicker(time.Second * 10)
			defer tick.Stop()
			for {
				select {
				case <-stop:
					return
				case <-tick.C:
					if err := c.WriteMessage(websocket.PingMessage, buf.Bytes()); err != nil {
						Error("receiveMsg websocket ping err:%v", err)
					}
				}
			}
		}(stopChan)

		// 这里阻塞接口 收到的消息
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				Error("read websocket err: %v", err)
				break
			}
			atomic.AddInt64(&receiveMsgCount, 1)
			flag := message[0]
			if flag == 0 {
				distribute(message[1:])
			}
		}
	}(c)
}

func distribute(data []byte) {
	msg := &pbchat.Message{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		log.Error("distribute proto Unmarshal err:%v", err)
		return
	}
	cmdLogic(msg.Ops, msg.Data)
}

func cmdLogic(ops pbchat.Operation, data []byte) {
	v := operationMsg(ops)
	if v == nil {
		log.Error("ops %v not find parse message", ops)
		return
	}
	err := proto.Unmarshal(data, v)
	if err != nil {
		log.Error("ops %v message parse err %v", ops, err)
		return
	}

	buf, err := json.MarshalIndent(v, "\t", "    ")
	if err != nil {
		log.Error("json marshal indent err %v", err)
		return
	}

	log.Info("ops %v message  %v", ops, string(buf))
	fmt.Println(string(buf))
}

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
