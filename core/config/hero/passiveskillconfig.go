package hero

import (
	"encoding/json"
)

// PassiveEffectType 被动技能效果类型
type PassiveEffectType int32

const (
	// EffectTypeInBattle 战斗中处理
	EffectTypeInBattle PassiveEffectType = 1
	// EffectTypeOutBattle 战斗外处理
	EffectTypeOutBattle PassiveEffectType = 2
)

var (
	passiveSkillCfg = &passiveSkillConfig{}
)

// 基本配置
type passiveSkillConfig struct {
	passiveSkillConfigs map[int32]*PassiveSkill
}

//技能激活配置
type PassiveSkill struct {
	SkillID      int32   `json:"SkillId"`      //技能id
	TriggerItem  int32   `json:"triggerItem"`  //触发类型
	TriggerValue []int32 `json:"triggerValue"` //触发类型
	EffectType   int32   `json:"effectType"`   //效果类型：
	SateCoeff    int32   `json:"satecoeff"`    //战力计算补偿系数
	ScoreCoeff   int32   `json:"scorecoeff"`   //系统评分补偿系数
}

func (c *passiveSkillConfig) Set(text []byte) error {
	var ec = passiveSkillConfig{passiveSkillConfigs: make(map[int32]*PassiveSkill)}
	var data []*PassiveSkill
	if err := json.Unmarshal(text, &data); err != nil {
		return err
	}
	for _, v := range data {
		ec.passiveSkillConfigs[v.SkillID] = v
	}

	*c = ec
	return nil
}

func GetPassiveSkillConfig(id int32) *PassiveSkill {
	if passiveSkillCfg == nil {
		return nil
	}
	if v, ok := passiveSkillCfg.passiveSkillConfigs[id]; ok {
		return v
	}
	return nil
}

func GetAllPassiveSkillConfig() map[int32]*PassiveSkill {
	if passiveSkillCfg == nil {
		return nil
	}
	return passiveSkillCfg.passiveSkillConfigs
}

func (config *PassiveSkill) GetTriggerValue() []int32 {
	return config.TriggerValue
}

// PassiveSkillInBattle 取得场内生效的被动技能
func PassiveSkillInBattle(list []int32) (inb []int32) {
	inb = make([]int32, 0)

	for _, psid := range list {
		cc := GetPassiveSkillConfig(psid)
		if cc == nil {
			continue
		}

		if cc.EffectType == int32(EffectTypeInBattle) {
			inb = append(inb, psid)
		}
	}
	return
}
