package model

import (
	"encoding/json"
	"time"
	pblogger "xy3-proto/logger"
	"xy3-proto/pkg/log"
)

// MLogTutorialProgress 教程进度
type MLogTutorialProgress struct {
	ID       int64     `gorm:"primaryKey;autoIncrement"`
	PlayerID int64     `gorm:"column:PlayerID"`
	OS       int       `gorm:"column:OS"`
	WarZone  int       `gorm:"column:WarZone"`
	Server   int       `gorm:"column:Server"`
	StepID   int32     `gorm:"column:StepID"`
	Time     time.Time `gorm:"column:Time"`
}

func (t *MLogTutorialProgress) TableName() string {
	return "MLogTutorialProgress"
}

func UnmarshalToTLogTutorialProgress(msg *pblogger.LogMsg) (*MLogTutorialProgress, error) {
	obj := &MLogTutorialProgress{}

	if err := json.Unmarshal([]byte(msg.Json), obj); err != nil {
		log.Error("UnmarshalToTLogBattle err! err:%v msg:%v", err, msg)
		return nil, err
	}

	obj.OS = int(msg.Os)
	obj.Time = time.Unix(msg.Time, 0)
	return obj, nil
}
