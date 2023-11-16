package model

import (
	"encoding/json"

	pblogger "xy3-proto/logger"

	"xy3-proto/pkg/log"
)

const (
	MailOverdue_Ungeted int32 = 3
	MailOverdue_Geted   int32 = 4
)

type TLogMail struct {
	LogEmbed

	UserID     int64  `json:"userid"`
	Operate    int32  `json:"operate"`
	UUID       int64  `json:"uuid"`
	ID         int32  `json:"id"`
	State      int32  `json:"state"` //0未读，1-已读未领取，2-已读已领取，3-过期未领取，4-过期已领取
	CreateTime int64  `json:"createtime"`
	EndTime    int64  `json:"endtime"` //-1:永久
	Param      string `json:"param"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	ItemList   string `json:"itemlist"`
}

func (t *TLogMail) TableName() string {
	return "MLogTransfer"
}

func UnmarshalToTLogMail(msg *pblogger.LogMsg) (obj *TLogMail) {
	obj = &TLogMail{}

	if err := json.Unmarshal([]byte(msg.Json), obj); err != nil {
		log.Error("UnmarshalToTLogMail err! err:%v msg:%v", err, msg)
		return
	}

	obj.LogEmbed.CreatedAt = msg.Time
	obj.LogEmbed.OS = int32(msg.Os)
	return
}
