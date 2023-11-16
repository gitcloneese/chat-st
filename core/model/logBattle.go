package model

import (
	"encoding/json"
	"time"
	pblogger "xy3-proto/logger"
	"xy3-proto/pkg/log"
)

const (
	LEADER_POS    = 0
	BOSS_POS      = 7
	FRIENDLY_CAMP = 0
	ENEMY_CAMP    = 1
)

// MLogBattle 战斗记录
type MLogBattle struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"`
	PlayerId     int64     `gorm:"column:PlayerId"`
	OS           int       `gorm:"column:OS"`
	WarZone      int       `gorm:"column:WarZone"`
	Server       int       `gorm:"column:Server"`
	BattleType   int       `gorm:"column:BattleType"` //0::降魔  1::登天塔 2::斗法  3::诸天斗法 4::天梯斗法 5::三界试练  6::方寸山
	Pos1         int32     `gorm:"column:Pos1"`
	Pos2         int32     `gorm:"column:Pos2"`
	Pos3         int32     `gorm:"column:Pos3"`
	Pos4         int32     `gorm:"column:Pos4"`
	Pos5         int32     `gorm:"column:Pos5"`
	Pos6         int32     `gorm:"column:Pos6"`
	Win          int       `gorm:"column:Win"`         // 0失败 1胜利
	Level        int       `gorm:"column:Level"`       // 等级
	CombatPower  int64     `gorm:"column:CombatPower"` // 战力
	Time         time.Time `gorm:"column:Time"`
	LevelBefore  int32     `gorm:"column:LevelBefore"`  //登天塔 战斗前层数， 斗法战斗前排名
	LevelAfter   int32     `gorm:"column:LevelAfter"`   // --
	Points       int32     `gorm:"column:Points"`       //登天塔积分变化，斗法积分变化
	ResultPoints int32     `gorm:"column:ResultPoints"` // --
	Position     string    `gorm:"column:Position"`     // 斗法当前段位 降魔关卡进度
}

func (t *MLogBattle) TableName() string {
	return "MLogBattle"
}

// MLogBattleParam 战斗参数
type MLogBattleParam struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	LogBattleID int64     `gorm:"column:LogBattleID"` //战斗记录的自增ID: MLogBattle.ID
	Time        time.Time `gorm:"column:Time"`
	BattleParam string    `gorm:"column:BattleParam"` // 战斗参数
}

func (t *MLogBattleParam) TableName() string {
	return "MLogBattleParam"
}

// MLogBattleResult 每个战斗者的战斗结果
type MLogBattleResult struct {
	ID          int64 `gorm:"primaryKey;autoIncrement"`
	LogBattleID int64 `gorm:"column:LogBattleID"`
	Camps       int32 `gorm:"column:Camps"`
	HeroID      int32 `gorm:"column:HeroID"`
	Position    int32 `gorm:"column:Pos"`
	DamageDealt int64 `gorm:"column:Hurt"`
	LostHp      int64 `gorm:"column:LostHp"`
	Healing     int64 `gorm:"column:Healing"`
}

func (t *MLogBattleResult) TableName() string {
	return "MLogBattleResult"
}

type MLogBattleInfo struct {
	Battle        MLogBattle
	BattleParam   MLogBattleParam
	BattleResults []MLogBattleResult
}

func UnmarshalToTLogBattle(msg *pblogger.LogMsg) (*MLogBattleInfo, error) {
	obj := &MLogBattleInfo{}

	if err := json.Unmarshal([]byte(msg.Json), obj); err != nil {
		log.Error("UnmarshalToTLogBattle err! err:%v msg:%v", err, msg)
		return nil, err
	}

	obj.Battle.OS = int(msg.Os)
	obj.Battle.Time = time.Unix(msg.Time, 0)
	obj.BattleParam.Time = obj.Battle.Time
	return obj, nil
}
