package model

import (
	"encoding/json"
	"time"
	pblogger "xy3-proto/logger"

	"xy3-proto/pkg/log"
)

// MLogRanking 斗法排名数据， 每日00：00
type MLogRanking struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	OS          int       `gorm:"column:OS"`
	WarZone     int64     `gorm:"column:WarZone"`
	Server      int       `gorm:"column:Server"`
	PlayerId    int64     `gorm:"column:PlayerId"`
	CombatPower int64     `gorm:"column:CombatPower"`
	Ranking     int64     `gorm:"column:Ranking"`
	Time        time.Time `gorm:"column:Time"`
	Position    string    `gorm:"column:Position"`
	BattleType  int32     `gorm:"column:BattleType"`
	Points      int64     `gorm:"column:Points"`
	PlayerName  string    `gorm:"column:PlayerName"`
}

func (m *MLogRanking) TableName() string {
	return "MLogRanking"
}

func UnmarshalToMLogRanking(msg *pblogger.LogMsg) (*MLogRanking, error) {
	obj := &MLogRanking{}

	if err := json.Unmarshal([]byte(msg.Json), obj); err != nil {
		log.Error("UnmarshalToMLogRanking err! err:%v msg:%v", err, msg)
		return nil, err
	}
	obj.Time = time.Unix(msg.Time, 0)

	return obj, nil
}
