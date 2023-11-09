package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
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
func TestPlayerSendMessage(info *pblogin.LoginRsp, wg *sync.WaitGroup) {
	defer wg.Done()
	// 发送消息
	var count int32
	for count < int32(ChatCount) {
		count++
		err := sendMessage(info, count)
		if err != nil {
			log.Printf("玩家:%v 发送聊天失败:%v\n", info.PlayerID, err)
		}
		time.Sleep(time.Millisecond * 10)
	}
}

// TestSendMessage
// 压测发送消息
func TestSendMessage() {
	playerNums := len(PlayerTokens)
	log.Printf("================开始压测发送消息!!! 玩家数量:%v 每个玩家发送:%v次================\n", playerNums, ChatCount)

	now := time.Now()

	wg := &sync.WaitGroup{}
	wg.Add(playerNums)
	for _, v := range PlayerTokens {
		go TestPlayerSendMessage(v, wg)
	}
	wg.Wait()

	latency := time.Since(now).Seconds()
	log.Printf("================压测发送消息完成!!! 用时:%v s================\n", latency)
}

// TestReceiveMessageUnBlock
// 压测接收消息 建立连接后就返回
// 只为测ws连接数量 玩家数量
func TestReceiveMessageUnBlock() {
	playerNums := len(PlayerTokens)
	log.Printf("================开始压测接收消息!!! 不阻塞接收 只为测ws连接数量 玩家数量:%v================\n", playerNums)

	now := time.Now()
	wg := new(sync.WaitGroup)
	for _, v := range PlayerTokens {
		wg.Add(1)
		go receiveMsg(v, wg)
	}
	wg.Wait()
	latency := time.Since(now).Seconds()
	log.Printf("================结束压测接收消息!!! 玩家数量:%v 用时:%vs ================\n", playerNums, latency)
}

// TestReceiveMessageBlock
// 压测接收消息 阻塞接收
func TestReceiveMessageBlock() {
	playerNums := len(PlayerTokens)
	log.Printf("================开始压测接收消息!!! 阻塞接收 玩家数量:%v================\n", playerNums)

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

	//log.Printf("设置zoneServer信息成功: player:%v", info.PlayerID)
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

	//log.Printf("发送信息成功: player:%v :Msg:%v", info.PlayerID, msg)
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

	//log.Printf("connect to chat %v", u.String())

	c, _, err1 := websocket.DefaultDialer.Dial(u.String(), reqHeader)
	if err1 != nil {
		panic(err1)
	}
	defer c.Close()

	if wg != nil {
		return // 这里只测试ws连接数， 不做其他处理
	} else {
		// 这里阻塞接口 收到的消息
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read websocket err: ", err)
				break
			}
			flag := message[0]
			if flag == 0 {
				distribute(message[1:])
			}
		}
	}
}
