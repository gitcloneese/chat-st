package hero

import (
	"encoding/json"
	battle "xy3-proto/battle"
)

var (
	heroCfg = &heroConfig{}
)

// 基本配置
type heroConfig struct {
	heroConfigs map[int32]*Hero
	raceMap     map[int32]int32
}

// 仙人配置
type Hero struct {
	HeroID           int32       `json:"heroId"`           //仙人ID
	UnitName         string      `json:"unitName"`         //仙人名称
	Sex              int32       `json:"sex"`              //性别
	InitQuality      int32       `json:"initQuality"`      //初始品质
	FragmentID       int32       `json:"fragmentId"`       //仙人对应到物品表中的碎片ID
	StarCoeff        [10]float64 `json:"starcoefficient"`  //用于战力计算的星级系数
	Race             int32       `json:"race"`             //种族 仙 佛 巫 妖
	Type             int32       `json:"actorType"`        //类型 1防御 2 辅助 3攻击 4突击
	Attack           int64       `json:"attack"`           //攻击
	Pdef             int64       `json:"pDef"`             //物理防御
	Mdef             int64       `json:"mDef"`             //魔法防御
	Hp               int64       `json:"hp"`               //生命
	HighAttr         [][2]int64  `json:"attribute"`        //高级属性
	SkillID          int32       `json:"skillId"`          //普攻技能ID
	AngerSkillID     int32       `json:"angerSkillID"`     //怒气技能ID
	PassiveSkillID   []int32     `json:"passiveSkill"`     //怒气技能ID列表
	RecruitSoulCount int32       `json:"recruitSoulCount"` //招募需要的魂魄数

	AttrSet map[battle.AttributeKey]int64
}

func (c *heroConfig) Set(text []byte) error {
	var ec = heroConfig{heroConfigs: make(map[int32]*Hero), raceMap: make(map[int32]int32)}
	var data []*Hero
	if err := json.Unmarshal(text, &data); err != nil {
		return err
	}
	for _, v := range data {
		ec.heroConfigs[v.HeroID] = v

		v.AttrSet = map[battle.AttributeKey]int64{}
		v.AttrSet[battle.AttributeKey_HP] = v.Hp
		v.AttrSet[battle.AttributeKey_ATK] = v.Attack
		v.AttrSet[battle.AttributeKey_PDEF] = v.Pdef
		v.AttrSet[battle.AttributeKey_MDEF] = v.Mdef

		ec.raceMap[v.Race] = 1
	}

	*c = ec
	return nil
}

func GetHeroConfig(id int32) *Hero {
	if heroCfg == nil {
		return nil
	}
	if v, ok := heroCfg.heroConfigs[id]; ok {
		return v
	}
	return nil
}

func GetAllHeroConfig() map[int32]*Hero {
	if heroCfg == nil {
		return nil
	}
	return heroCfg.heroConfigs
}

func GetAllRace() map[int32]int32 {
	return heroCfg.raceMap
}
