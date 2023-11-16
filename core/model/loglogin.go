package model

import (
	"encoding/json"
	"time"
	pblogger "xy3-proto/logger"

	"xy3-proto/pkg/log"
)

type TLogFields struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	UnionId   string    `gorm:"column:UnionId"`
	PlayerId  int64     `gorm:"column:PlayerId"`
	Level     int       `gorm:"column:Level"`
	IP        string    `gorm:"column:IP"`
	DeviceUID string    `gorm:"column:DeviceUID"`
	OS        int       `gorm:"column:OS"`
	WarZone   int       `gorm:"column:WarZone"`
	Server    int       `gorm:"column:Server"`
	Time      time.Time `gorm:"column:Time"`
	LoginType int       `gorm:"column:LoginType"`
	LoginTime int64     `gorm:"-"`
	Expired   bool      `gorm:"-"`
}

type TLogLogin struct {
	TLogFields
}

func (t *TLogLogin) TableName() string {
	return "MLogLogin"
}

func UnmarshalToTLogLogin(msg *pblogger.LogMsg) (*TLogLogin, error) {
	obj := &TLogLogin{}

	if err := json.Unmarshal([]byte(msg.Json), obj); err != nil {
		log.Error("UnmarshalToTLogLogin err! err:%v msg:%v", err, msg)
		return nil, err
	}

	return obj, nil
}

type TLogRegister struct {
	TLogFields
}

func (t *TLogRegister) TableName() string {
	return "MLogRegister"
}

func UnmarshalToTLogRegister(msg *pblogger.LogMsg) (*TLogRegister, error) {
	obj := &TLogRegister{}

	if err := json.Unmarshal([]byte(msg.Json), obj); err != nil {
		log.Error("UnmarshalToTLogRegister err! err:%v msg:%v", err, msg)
		return nil, err
	}

	return obj, nil
}
