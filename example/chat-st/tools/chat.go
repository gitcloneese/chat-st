package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/gorilla/websocket"
	"github.com/panjf2000/ants/v2"
	"io"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	pblogin "xy3-proto/login"
	pbchat "xy3-proto/new-chat"
	"xy3-proto/pkg/log"
)

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
	if len(PlayerTokens) < 1 {
		panic("PlayerTokens must be at least 1")
	}
	var loginInfo *pblogin.LoginRsp
	for _, token := range PlayerTokens {
		loginInfo = token
		break
	}
	TestPlayerSendMessage(loginInfo)
}

// TestOneSendMessage
// 一个玩家发送消息
func TestOneSendMessage() {
	if len(PlayerTokens) < 1 {
		log.Error("压测发送消息, 玩家数量不足")
		return
	}
	wg := &sync.WaitGroup{}
	wg.Add(ChatCount)

	var info *pblogin.LoginRsp
	for _, v := range PlayerTokens {
		info = v
		break
	}

	var chatNum int
	for chatNum < ChatCount {
		go func() {
			err := sendMessage(info, int32(chatNum))
			if err != nil {
			}
		}()
		chatNum++
	}
	wg.Wait()
}

func RunTestReceiveMessage() {
	RunReceiveMsgWithLogTick("receiveMessageBlock", TestReceiveMessageBlock)
}

// TestReceiveMessageBlock
// 压测接收消息 阻塞接收
func TestReceiveMessageBlock() {
	wg := new(sync.WaitGroup)
	p, _ := ants.NewPoolWithFunc(C, func(i interface{}) {
		go receiveMsg(i.(*pblogin.LoginRsp))
		wg.Done()
	})
	defer p.Release()
	for _, v := range PlayerTokens {
		wg.Add(1)
		err := p.Invoke(v)
		if err != nil {
		}
	}
	wg.Wait()
}

func bearToken(token string) string {
	return fmt.Sprintf("Bearer %v", token)
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

// 接收消息
func receiveMsg(info *pblogin.LoginRsp) {
	uri, err := url.Parse(ChatAddr)
	if err != nil {
		panic("receiveMsg parse url failed" + err.Error())
	}

	scheme := "ws"
	if uri.Scheme == "https" {
		scheme = "wss"
	}
	u := url.URL{
		Scheme: scheme,
		Host:   uri.Host,
		Path:   apiConnectChatPath,
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
	defer c.Close()

	// 记录长练级数量
	atomic.AddInt64(&chatConnectCount, 1)
	defer atomic.AddInt64(&chatConnectCount, -1)

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
}
