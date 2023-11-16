package model

import (
	"encoding/json"
	"time"
	pblogger "xy3-proto/logger"

	"xy3-proto/pkg/log"
)

type MLogTask struct {
	ID       int64     `gorm:"primaryKey;autoIncrement"`
	PlayerId int64     `gorm:"column:PlayerId"`
	OS       int       `gorm:"column:OS"`
	WarZone  int       `gorm:"column:WarZone"`
	Server   int       `gorm:"column:Server"`
	Type     int       `gorm:"column:Type"` // 1::每日日常 2:: 每日限时 3 :: 每周周常任务
	TaskId   int32     `gorm:"column:TaskId"`
	Points   int32     `gorm:"column:Points"`
	Time     time.Time `gorm:"column:Time"`
	Reward   string    `gorm:"column:Reward"`
}

func (t *MLogTask) TableName() string {
	return "MLogTask"
}

func UnmarshalToTLogTask(msg *pblogger.LogMsg) (*MLogTask, error) {
	obj := &MLogTask{}

	if err := json.Unmarshal([]byte(msg.Json), obj); err != nil {
		log.Error("UnmarshalToTLogTask err! err:%v msg:%v", err, msg)
		return nil, err
	}

	obj.Time = time.Unix(msg.Time, 0)
	obj.OS = int(msg.Os)

	return obj, nil
}
