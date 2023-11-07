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
	"time"
	pblogin "xy3-proto/login"
	pbchat "xy3-proto/new-chat"
)

func chat(info *pblogin.LoginRsp) {
	// 设置区服
	err := setZoneServer(info)
	if err != nil {
		log.Printf("开始聊天,玩家:%v 设置区服失败 error:%v\n", info.PlayerID, err)
		return
	}
	// 发送消息
	go func() {
		var count int32
		for count < int32(ChatCount) {
			count++
			err := sendMessage(info, count)
			if err != nil {
				log.Println(err)
			}
			time.Sleep(time.Millisecond * 10)
		}
	}()
	// 接收消息
	receiveMsg(info)
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
	req.Header.Set("Authorization", info.PlayerToken)

	// 本地test
	if Local != 0 {
		req.Header.Set("userid", fmt.Sprintf("%v", info.PlayerID))
	}

	resp, err := HttpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("accountRoleList failed, status code: %v", resp.StatusCode)
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

	log.Printf("设置zoneServer信息成功: player:%v", info.PlayerID)
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%v%v", ChatAddr, apiSetZoneServerPath), bytes.NewReader(reqB))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", info.PlayerToken)

	// 本地test
	if Local != 0 {
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

	log.Printf("发送信息成功: player:%v :Msg:%v", info.PlayerID, msg)
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
		Host:   ChatAddr,
		Path:   apiConnectChatPath,
	}

	reqHeader := http.Header{}
	if Local != 0 {
		u.RawQuery = fmt.Sprintf("userid=%v", info.PlayerID)
	} else {
		reqHeader.Add("Authorization", info.PlayerToken)
	}

	log.Printf("connect to chat %v", u.String())

	c, _, err1 := websocket.DefaultDialer.Dial(u.String(), reqHeader)
	if err1 != nil {
		panic(err1)
	}
	defer c.Close()
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
