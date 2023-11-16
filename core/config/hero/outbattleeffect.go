package hero

import (
	"encoding/json"
	"errors"
	"strconv"

	"x-server/core/model"
	coremodel "x-server/core/model"
)

// OutType 场外技能效果类型
type OutType int32

const (
	// OutTypeAttr 影响属性
	OutTypeAttr OutType = 1
	// OutTypeReplaceNormal 替换普攻技能
	OutTypeReplaceNormal OutType = 2
	// OutTypeReplaceAnger 替换怒气技能
	OutTypeReplaceAnger OutType = 3
	// OutTypeAddInBattleBuff 增加场内buff技能(也是被动技)
	OutTypeAddInBattleBuff OutType = 4
	// OutTypeReplaceInBattleBuff 替换场内buff技能（也是被动技能）
	OutTypeReplaceInBattleBuff OutType = 5
)

var (
	outbattleeffectcfg = &outbattleeffectconfig{}
)

// outbattleeffectconfig .
type outbattleeffectconfig struct {
	maps map[int32]*OutBattleEffect
}

// OutBattleEffect .
type OutBattleEffect struct {
	EffectID              int32      `json:"effectId"`
	Type                  int32      `json:"type"`
	Attr                  [][2]int64 `json:"attr"`
	ReplaceSkill          int32      `json:"replaceSkill"`
	PassiveSkillToReplace int32      `json:"passiveSkillToReplace"`
	Order                 int32      `json:"order"` //当有多个替换技能，比如多次替换普攻，多次替换怒气，这个就是顺序

	AttrSet *coremodel.AttributeSet
}

func (c *outbattleeffectconfig) Set(text []byte) error {
	var ec = outbattleeffectconfig{maps: make(map[int32]*OutBattleEffect)}
	var data []*OutBattleEffect
	if err := json.Unmarshal(text, &data); err != nil {
		return err
	}

	for _, v := range data {
		ec.maps[v.EffectID] = v
		if v.Type != int32(OutTypeAttr) &&
			v.Type != int32(OutTypeReplaceNormal) &&
			v.Type != int32(OutTypeReplaceAnger) &&
			v.Type != int32(OutTypeAddInBattleBuff) &&
			v.Type != int32(OutTypeReplaceInBattleBuff) {
			return errors.New("out_battle_effect err! EffectID:" + strconv.Itoa(int(v.EffectID)) + " Type-err!")
		}

		v.AttrSet = model.MakeAttributeSet(v.Attr)
	}

	*c = ec
	return nil
}

// OutBattleEffectOne 通过配置ID取得一条行配置
func OutBattleEffectOne(ID int32) *OutBattleEffect {
	if outbattleeffectcfg == nil {
		return nil
	}
	if v, ok := outbattleeffectcfg.maps[ID]; ok {
		return v
	}
	return nil
}

// OutBattleEffectMap 取得整个配置表
func OutBattleEffectMap() map[int32]*OutBattleEffect {
	if outbattleeffectcfg == nil {
		return nil
	}
	return outbattleeffectcfg.maps
}
