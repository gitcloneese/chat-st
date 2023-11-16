package model

import (
	"encoding/json"
	"time"
	pblogger "xy3-proto/logger"

	"xy3-proto/pkg/log"
)

type TLogResource struct {
	LogEmbed
	UserID   int64  `json:"userid"`
	Operate  int32  `json:"operate"`
	Type     int32  `json:"type"` //0-资源，1-道具，2-仙人，5-头像类，6-灵宝，7-仙器
	ID       int32  `json:"id"`
	PreNum   int64  `json:"prenum"`   //操作之前的数值
	Num      int64  `json:"num"`      //改变量
	NowNum   int64  `json:"nownum"`   //操作之后的数值，
	Source   int32  `json:"source"`   //物品来源
	Extend   int64  `json:"extend"`   //扩展记录
	UserName string `json:"userName"` //
	WarZone  int    `json:"warzone"`
	Server   int    `json:"server"`
	Level    int32  `json:"level"` //
}

type MLogTransfer struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`
	OS         int32     `gorm:"column:OS"`
	WarZone    int       `gorm:"column:WarZone"`
	Server     int       `gorm:"column:Server"`
	PlayerId   int64     `gorm:"column:PlayerId"`
	PlayerName string    `gorm:"column:PlayerName"`
	ItemID     int32     `gorm:"column:ItemId"`
	ItemNum    int64     `gorm:"column:ItemNum"`
	ResultNum  int64     `gorm:"column:ResultNum"`
	Source     int32     `gorm:"column:Source"`
	Time       time.Time `gorm:"column:Time"`
	ItemTypeId int32     `gorm:"column:ItemTypeId"` //0-资源，1-道具，2-仙人，5-头像类，6-灵宝，7-仙器
	Level      int32     `gorm:"column:Level"`      // 玩家等级
}

func (t *MLogTransfer) TableName() string {
	return "MLogTransfer"
}

func UnmarshalToTLogResource(msg *pblogger.LogMsg) (*TLogResource, error) {
	obj := &TLogResource{}

	if err := json.Unmarshal([]byte(msg.Json), obj); err != nil {
		log.Error("UnmarshalToTLogMail err! err:%v msg:%v", err, msg)
		return nil, err
	}

	obj.LogEmbed.CreatedAt = msg.Time
	obj.LogEmbed.OS = int32(msg.Os)

	return obj, nil
}
