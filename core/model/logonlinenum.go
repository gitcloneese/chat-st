package model

import (
	"encoding/json"
	"time"
	pblogger "xy3-proto/logger"

	"xy3-proto/pkg/log"
)

type TLogOnline struct {
	ID      int       `gorm:"primaryKey;autoIncrement"`
	OS      int       `gorm:"column:OS"`
	WarZone int       `gorm:"column:WarZone"`
	Server  int       `gorm:"column:Server"`
	Time    time.Time `gorm:"column:Time"`
	Online  int       `gorm:"column:Online"`
}

func (t *TLogOnline) TableName() string {
	return "MLogOnline"
}

func UnmarshalToTLogOnline(msg *pblogger.LogMsg) (*TLogOnline, error) {
	obj := &TLogOnline{}

	if err := json.Unmarshal([]byte(msg.Json), obj); err != nil {
		log.Error("UnmarshalToTLogOnline err! err:%v msg:%v", err, msg)
		return nil, err
	}

	return obj, nil
}
