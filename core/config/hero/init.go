package hero

import (
	"xy3-proto/pkg/conf/paladin"
)

// Init 初始化
func Init(mm map[string]paladin.Setter) {
	mm["hero.json"] = heroCfg
	mm["out_battle_effect.json"] = outbattleeffectcfg
	mm["passive_skill.json"] = passiveSkillCfg
}
