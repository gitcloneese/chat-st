package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
	pblogin "xy3-proto/login"
	pbchat "xy3-proto/new-chat"
)

func TestChat() {
	// 异步发送消息
	go TestSendMessage()
	// 接收消息
	TestReceiveMessageBlock()
}

// TestPlayerSendMessage
// 单个玩家发送消息
func TestPlayerSendMessage(info *pblogin.LoginRsp, wg *sync.WaitGroup, delayMillSecond ...int) {
	if wg != nil {
		defer wg.Done()
	}
	// 发送消息
	var count int32
	for count < int32(ChatCount) {
		count++
		err := sendMessage(info, count)
		if err != nil {
			log.Errorf("玩家:%v 发送聊天失败:%v ", info.PlayerID, err)
		}
		if len(delayMillSecond) > 0 {
			time.Sleep(time.Millisecond * time.Duration(delayMillSecond[0]))
		}
	}
}

// TestSendMessage
// 压测发送消息
func TestSendMessage() {
	playerNums := len(PlayerTokens)
	chatNum := int32(float64(playerNums) * percentChatPlayers)
	log.Infof("================开始压测发送消息!!! 玩家数量:%v 每个玩家发送:%v次================ ", chatNum, ChatCount)

	now := time.Now()
	temp1 := atomic.LoadInt64(&QAcc)

	wg := &sync.WaitGroup{}
	var nowCount int32
	wg.Add(playerNums)
	for _, v := range PlayerTokens {
		go TestPlayerSendMessage(v, wg)
		nowCount++
		// 聊天玩家百分比
		if nowCount >= chatNum {
			break
		}
	}
	wg.Wait()
	temp2 := atomic.LoadInt64(&QAcc)
	qs := temp2 - temp1

	latency := time.Since(now).Seconds()
	log.Infof("================压测发送消息完成!!! 用时:%vs 请求总数:%v QPS:%v ================ ", latency, qs, float64(qs)/latency)
}

// TestOneSendMessage
// 一个玩家发送消息
func TestOneSendMessage() {
	playerNums := len(PlayerTokens)
	chatNum := int32(float64(playerNums) * percentChatPlayers)
	log.Infof("================开始压测发送消息!!! 玩家数量:%v 每个玩家发送:%v次================ ", chatNum, ChatCount)

	now := time.Now()
	temp1 := atomic.LoadInt64(&QAcc)

	var one *pblogin.LoginRsp
	for k := range PlayerTokens {
		one = PlayerTokens[k]
		break
	}
	count := 0
	wg := &sync.WaitGroup{}
	wg.Add(c)
	perNum := ChatCount / c
	for count < c {
		TestGoroutineSendMessage(one, wg, perNum)
		count++
	}
	wg.Wait()
	temp2 := atomic.LoadInt64(&QAcc)
	qs := temp2 - temp1

	latency := time.Since(now).Seconds()
	log.Infof("================压测发送消息完成!!! 用时:%vs 请求总数:%v QPS:%v ================ ", latency, qs, float64(qs)/latency)
}

// TestGoroutineSendMessage 单个携程发送
func TestGoroutineSendMessage(info *pblogin.LoginRsp, wg *sync.WaitGroup, n int) {
	if wg != nil {
		defer wg.Done()
	}
	var i int
	for i < n {
		TestPlayerSendMessage(info, nil, 2)
		i++
	}
}

// TestReceiveMessageUnBlock
// 压测接收消息 建立连接后就返回
// 只为测ws连接数量 玩家数量
func TestReceiveMessageUnBlock() {
	playerNums := len(PlayerTokens)
	log.Infof("================开始压测接收消息!!! 不阻塞接收 只为测ws连接数量 玩家数量:%v================ ", playerNums)

	now := time.Now()
	wg := new(sync.WaitGroup)
	for _, v := range PlayerTokens {
		wg.Add(1)
		go receiveMsg(v, wg)
	}
	wg.Wait()
	latency := time.Since(now).Seconds()
	log.Infof("================结束压测接收消息!!! 玩家数量:%v 用时:%vs ================ ", playerNums, latency)
}

// TestReceiveMessageBlock
// 压测接收消息 阻塞接收
func TestReceiveMessageBlock() {
	playerNums := len(PlayerTokens)
	log.Infof("================开始压测接收消息!!! 阻塞接收 玩家数量:%v================ ", playerNums)

	for _, v := range PlayerTokens {
		go receiveMsg(v, nil)
	}
}

func bearToken(token string) string {
	return fmt.Sprintf("Bearer %v", token)
}

// 设置区服
func setZoneServer(info *pblogin.LoginRsp) error {
	reqB, err := json.Marshal(pbchat.SetZoneServerReq{
		ZoneId:   "test-st", // 压力测试
		ServerId: "100",
		PlayerId: info.PlayerID,
	})
	if err != nil {
		return err
	}
	defer atomic.AddInt64(&QAcc, 1)
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

	//log.Infof("设置zoneServer信息成功: player:%v", info.PlayerID)
	return nil
}

func generateMessage(playerId int64, accId int32) string {
	return fmt.Sprintf("我是%v, 这是第%v次发言。", playerId, accId)
}

// 发送消息
func sendMessage(info *pblogin.LoginRsp, chatNums int32) error {
	msg := generateMessage(info.PlayerID, chatNums)
	reqB, err := json.Marshal(pbchat.SendChat{
		Msg:        msg,
		FromPlayer: info.PlayerID,
		RoomType:   pbchat.RoomType_RT_World,
	})
	if err != nil {
		return err
	}
	defer atomic.AddInt64(&QAcc, 1)
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

	//log.Infof("发送信息成功: player:%v :Msg:%v", info.PlayerID, msg)
	return nil
}

// 接收消息
func receiveMsg(info *pblogin.LoginRsp, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

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

	//log.Infof("connect to chat %v", u.String())

	c, _, err1 := websocket.DefaultDialer.Dial(u.String(), reqHeader)
	if err1 != nil {
		log.Errorf("receiveMsg websocket Dail err:%v", err1)
		return
	}
	defer c.Close()

	if wg != nil {
		return // 这里只测试ws连接数， 不做其他处理
	} else {
		// 这里阻塞接口 收到的消息
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Errorf("read websocket err: %v", err)
				break
			}
			flag := message[0]
			if flag == 0 {
				distribute(message[1:])
			}
		}
	}
}
