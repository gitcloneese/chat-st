package model

import (
	battle "xy3-proto/battle"
	pbfriend "xy3-proto/friend"
)

// 好友信息
type FriendInfo struct {
	RoleID          int64 // 好友ID
	ServerID        int64 // 服务器ID
	GiveFriendPoint bool  // 是否已赠送好友点 fale.未赠送 true.已赠送
	GetFriendPoint  int32 // 是否已领取好友点 0.未赠送 1.已赠送 2.已领取
	Level           int32 // 亲密度等级
	Exp             int32 // 亲密度经验值
	OperationTime   int64 // 操作时间
}

// 好友点操作存储对象
type FriendPoint struct {
	RoleID         int64 // 角色ID
	GetPointCount  int32 // 领取好友数量
	GivePointCount int32 // 赠送好友点数量
}

// 好友请求
type FriendRequest struct {
	RoleID int64 // 好友id
	Time   int64 // 时间
}

////////////////////////////////////////////////////////////////////////////////////////////
// 租借英雄

const (
	ExclusiveWeapon_BarMax  = 3
	ExclusiveWeapon_RuneMax = 2
)

// 租借的英雄
type LeaseHero struct {
	OwnerID         int64                    // 拥有者id
	LeaseID         int64                    // 借出者id
	LeaseTime       int64                    // 到期的时间
	Hero            *LeaseHeroObj            // 英雄
	ExclusiveWeapon *LeaseExclusiveWeaponObj // 专属
	CD              int64                    // 冷却时间
}

func (l *LeaseHero) Format() (leaseHero *pbfriend.LeaseHero) {
	leaseHero = &pbfriend.LeaseHero{
		Hero: &pbfriend.Hero{
			Uuid:          l.Hero.UUID,
			Id:            l.Hero.HeroID,
			Star:          l.Hero.HeroStar,
			SubStar:       l.Hero.SubStar,
			Quality:       l.Hero.Quality,
			Awaken:        l.Hero.Awaken,
			Skin:          l.Hero.Skin,
			LingGenClicks: l.Hero.LinggenClicks,
			Power:         l.Hero.Power,
		},
		ExclusiveWeapon: &pbfriend.ExclusiveWeapon{
			Id:      l.ExclusiveWeapon.ExclusiveWeaponID,
			Star:    l.ExclusiveWeapon.ExclusiveWeaponStar,
			HoleNum: l.ExclusiveWeapon.HoleNum,
			SuitId:  l.ExclusiveWeapon.SuitID,
		},
		LeaseTime: l.LeaseTime,
	}

	leaseHero.ExclusiveWeapon.Rune = make([]int32, 0)
	for _, r := range l.ExclusiveWeapon.Rune {
		leaseHero.ExclusiveWeapon.Rune = append(leaseHero.ExclusiveWeapon.Rune, r)
	}
	leaseHero.ExclusiveWeapon.VBar = make([]*pbfriend.ExclusiveWeaponBar, 0)
	for _, vbar := range l.ExclusiveWeapon.VBar {
		if vbar == nil {
			continue
		}
		pbBar := &pbfriend.ExclusiveWeaponBar{
			Key:    vbar.Key,
			Num:    vbar.Num,
			Max:    vbar.Max,
			Stage:  vbar.Stage,
			Locked: vbar.Locked,
		}
		leaseHero.ExclusiveWeapon.VBar = append(leaseHero.ExclusiveWeapon.VBar, pbBar)
	}
	return leaseHero
}

type LeaseHeroObj struct {
	UUID          int64   // 英雄唯一id
	HeroID        int32   // 英雄配置id
	HeroStar      int32   // 英雄星级
	SubStar       int32   // 英雄星级子阶段
	Quality       int32   // 英雄品质
	Awaken        int32   // 英雄觉醒等级
	Skin          int32   // 英雄皮肤
	LinggenClicks []int32 // 玩家对该仙人灵根的点击列表
	Power         int64   //仙人战斗力
}

type LeaseExclusiveWeaponObj struct {
	ExclusiveWeaponID   int32                                               // 专属武器ID
	ExclusiveWeaponStar int32                                               // 专属武器星级(6星表示觉醒1段，7星表示觉醒2段)
	HoleNum             int32                                               // 专属武器当前开孔数
	SuitID              int32                                               // 专属武器镶嵌符文激活的套装ID
	Rune                [ExclusiveWeapon_RuneMax]int32                      // 专属武器镶嵌情况
	VBar                [ExclusiveWeapon_BarMax]*LeaseExclusiveWeaponBarObj // 专属武器当前洗练条
}

type LeaseExclusiveWeaponBarObj struct {
	Key    battle.AttributeKey `bson:"key"`    // 属性类型
	Num    uint32              `bson:"num"`    // 属性当前值
	Max    uint32              `bson:"max"`    // 属性最大值
	Stage  uint32              `bson:"stage"`  // 属性阶段
	Locked bool                `bson:"locked"` // 是否锁定
}

type LeaseRequest struct {
	ID     int64 // 唯一id
	RoleID int64 // 申请人id
	HeroID int32 // 申请的英雄id
	Time   int64 // 时间
	State  int32 // 状态 1.同意 2.归还
}

type LeaseTask struct {
	TaskID int32 // 任务id
	Num    int32 // 完成数量
	State  int32 // 0.未完成 1.已完成 2.已领取
}
