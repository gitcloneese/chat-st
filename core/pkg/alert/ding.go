package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type dingText struct {
	Content string `json:"content"`
}

type dingAt struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type dingMsg struct {
	Msgtype string   `json:"msgtype"`
	Text    dingText `json:"text"`
	At      dingAt   `json:"at"`
}

func SendDing(url, msg string) (err error) {
	if url == "" || msg == "" {
		err = fmt.Errorf("url or msg is nil")
		return err
	}
	ding := dingMsg{Msgtype: "text", Text: dingText{Content: msg}, At: dingAt{AtMobiles: nil, IsAtAll: false}}
	var sendData []byte
	if sendData, err = json.Marshal(ding); err != nil {
		return err
	}
	var req *http.Request
	if req, err = http.NewRequest(http.MethodPost, url, bytes.NewReader(sendData)); err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; encoding=utf-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return err
}
