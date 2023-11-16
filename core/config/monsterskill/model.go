package monsterskill

import "x-server/core/model"

const (
	// HeadIndexMin 最小阵容位
	HeadIndexMin = 1
	// HeadIndexMax 最大阵容位
	HeadIndexMax = 6
)

// TrigerItem .
type TrigerItem int32

const (
	// TriggerNone 没有触发条件，由其他系统独立获得维护
	TriggerNone TrigerItem = 0
	// TriggerLineup 阵容
	TriggerLineup TrigerItem = 1
	// TriggerStar 仙人星级
	TriggerStar TrigerItem = 2
	// TriggerHero 仙人
	TriggerHero TrigerItem = 3
	// TriggerLevel 等级
	TriggerLevel TrigerItem = 4
	// TriggerAwake 仙人觉醒等级
	TriggerAwake TrigerItem = 5
)

// FormationCondition .
type FormationCondition int32

const (
	// type1 齐上阵
	type1 FormationCondition = 1
	// type2 x个x种族齐上阵
	type2 FormationCondition = 2
	// type3 x个不同种族
	type3 FormationCondition = 3
	// type4 x个x类型齐上阵
	type4 FormationCondition = 4
	// type5 x个不同类型
	type5 FormationCondition = 5
	// type6 x个女性
	type6 FormationCondition = 6
	// type7 x个男性
	type7 FormationCondition = 7
	// type8 上阵Y个 指定的X个仙人列表
	type8 FormationCondition = 8
)

type MonsterHero struct {
	Index             int32               // 布阵位置
	HeroID            int32               // 仙人id
	Level             int32               // 仙人等级
	Star              int32               // 仙人星级
	Awaken            int32               // 觉醒
	SkillID           int32               // 英雄普攻技能
	SkillIDOrder      int32               // 英雄普攻技能order
	AngerSkillID      int32               // 英雄怒气技能
	AngerSkillIDOrder int32               // 英雄怒气技能order
	PassiveSkill      []int32             // 被动技能列表
	AttrSet           *model.AttributeSet // 属性集
}
