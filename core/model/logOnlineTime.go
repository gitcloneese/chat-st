package model

import (
	"encoding/json"
	"time"
	pblogger "xy3-proto/logger"

	"xy3-proto/pkg/log"
)

type MLogOnlineTime struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	PlayerId    int64     `gorm:"column:PlayerId"`
	OS          int       `gorm:"column:OS"`
	WarZone     int       `gorm:"column:WarZone"`
	Server      int       `gorm:"column:Server"`
	OnlineTime  int32     `gorm:"column:OnlineTime"`
	OfflineTime time.Time `gorm:"column:OfflineTime"`
	Time        int64     `gorm:"-"` // 用于记录玩家的登出时间戳， 写日志时 转换为OfflineTime
}

func (t *MLogOnlineTime) TableName() string {
	return "MLogOnlineTime"
}

func UnmarshalToLogOnlineTime(msg *pblogger.LogMsg) (*MLogOnlineTime, error) {
	obj := &MLogOnlineTime{}

	if err := json.Unmarshal([]byte(msg.Json), obj); err != nil {
		log.Error("UnmarshalToLogOnlineTime err! err:%v msg:%v", err, msg)
		return nil, err
	}

	obj.OS = int(msg.Os)
	obj.OfflineTime = time.Unix(msg.Time, 0)

	return obj, nil
}
