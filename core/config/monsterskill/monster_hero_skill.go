package monsterskill

import (
	cfg "x-server/core/config/hero"

	"xy3-proto/pkg/log"
)

// UpdateMonsterPassiveSkill 更新英雄本身技能技能,只处理Hero.xlsx:passiveSkill列的激活情况，其他系统的附加技能需要而外添加
func UpdateMonsterPassiveSkill(headIndex int32, monsterHeros []*MonsterHero) (skillID []int32) {
	skillID = make([]int32, 0)
	if headIndex < HeadIndexMin || headIndex > HeadIndexMax {
		return skillID
	}

	vtrigger := []TrigerItem{TriggerLineup, TriggerStar, TriggerHero, TriggerLevel, TriggerAwake}
	mapp := make(map[int32]bool)
	for _, vTrigger := range vtrigger {
		vv := GetMonsterActiveSkill(headIndex, monsterHeros, vTrigger)
		for _, psid := range vv {
			mapp[psid] = true
		}
	}

	for k := range mapp {
		skillID = append(skillID, k)
	}

	// h.PassiveSkillID = make([]int32, 0, len(mapp))
	// for k := range mapp {
	// 	h.PassiveSkillID = append(h.PassiveSkillID, k)
	// }
	return skillID
}

// GetMonsterActiveSkill 获取当前仙人激活的技能列表
func GetMonsterActiveSkill(headIndex int32, monsterHeros []*MonsterHero, trigerItem TrigerItem) (activeSkill []int32) {
	activeSkill = make([]int32, 0)

	if headIndex < HeadIndexMin || headIndex > HeadIndexMax {
		return activeSkill
	}

	var monsterHero *MonsterHero
	for _, mh := range monsterHeros {
		if mh.Index == headIndex {
			monsterHero = mh
		}
	}
	if monsterHero == nil {
		return activeSkill
	}

	heroCfg := cfg.GetHeroConfig(monsterHero.HeroID)
	for _, skillID := range heroCfg.PassiveSkillID {
		passiveSkillCfg := cfg.GetPassiveSkillConfig(skillID)
		if passiveSkillCfg == nil {
			// log.Error("GetActiveSkill passiveSkillConfig NotExist [%s]", skillID)
			continue
		}
		if passiveSkillCfg.TriggerItem != int32(trigerItem) {
			// log.Error("GetActiveSkill passiveSkillCfg.triggerItem Error [%s]", skillID)
			continue
		}
		isActive := false
		switch trigerItem {
		case TriggerLineup:
			isActive = triggerMonsterByLineup(monsterHero, passiveSkillCfg, monsterHeros)
		case TriggerStar:
			isActive = triggerMonsterByStar(monsterHero, passiveSkillCfg)
		case TriggerHero:
			isActive = triggerMonsterByHero(passiveSkillCfg, monsterHeros)
		case TriggerLevel:
			isActive = triggerMonsterByLevel(monsterHero, passiveSkillCfg)
		case TriggerAwake:
			isActive = triggerMonsterByAwake(monsterHero, passiveSkillCfg)
		}
		if isActive {
			activeSkill = append(activeSkill, skillID)
		}
	}
	return activeSkill
}

//nolint:all
func triggerMonsterByLineup(monsterHero *MonsterHero, config *cfg.PassiveSkill, monsterHeros []*MonsterHero) bool {
	condition := config.GetTriggerValue()
	switch condition[0] {
	case int32(type1):
		for _, heroID := range condition[1:] {
			isLineup := false
			for _, f := range monsterHeros {
				if f.HeroID == heroID {
					isLineup = true
				}
			}
			if !isLineup {
				return false
			}
		}
		return true
	case int32(type2):
		num := condition[1]
		race := condition[2]
		lineupNum := int32(0)
		for _, f := range monsterHeros {
			heroCfg := cfg.GetHeroConfig(f.HeroID)
			if heroCfg == nil {
				continue
			}
			if heroCfg.Race == race {
				lineupNum++
			}
		}
		if lineupNum >= num {
			return true
		}
		return false
	case int32(type3):
		num := condition[1]
		raceMap := make(map[int32]int32)
		for _, f := range monsterHeros {
			heroCfg := cfg.GetHeroConfig(f.HeroID)
			if heroCfg == nil {
				continue
			}
			raceMap[heroCfg.Race] = heroCfg.Race
		}
		if int32(len(raceMap)) >= num {
			return true
		}
		return false
	case int32(type4):
		num := condition[1]
		heroType := condition[2]
		lineupNum := int32(0)
		for _, f := range monsterHeros {
			heroCfg := cfg.GetHeroConfig(f.HeroID)
			if heroCfg == nil {
				continue
			}
			if heroCfg.Type == heroType {
				lineupNum++
			}
		}
		if lineupNum >= num {
			return true
		}
		return false
	case int32(type5):
		num := condition[1]
		typeMap := make(map[int32]int32)
		for _, f := range monsterHeros {
			heroCfg := cfg.GetHeroConfig(f.HeroID)
			if heroCfg == nil {
				continue
			}
			typeMap[heroCfg.Type] = heroCfg.Type
		}
		if int32(len(typeMap)) >= num {
			return true
		}
		return false
	case int32(type6):
		num := condition[1]
		lineupNum := int32(0)
		for _, f := range monsterHeros {
			heroCfg := cfg.GetHeroConfig(f.HeroID)
			if heroCfg == nil {
				continue
			}
			if heroCfg.Sex == 0 {
				lineupNum++
			}
		}
		if lineupNum >= num {
			return true
		}
		return false
	case int32(type7):
		num := condition[1]
		lineupNum := int32(0)
		for _, f := range monsterHeros {
			heroCfg := cfg.GetHeroConfig(f.HeroID)
			if heroCfg == nil {
				continue
			}
			if heroCfg.Sex == 1 {
				lineupNum++
			}
		}
		if lineupNum >= num {
			return true
		}
		return false
	case int32(type8):
		numY := int32(0)
		mapH := make(map[int32]bool)
		for i := 2; i < len(condition); i++ {
			mapH[condition[i]] = true
		}
		for _, f := range monsterHeros {
			if _, ok := mapH[f.HeroID]; ok {
				numY++
			}
		}
		return numY >= condition[1]
	default:
		log.Error("triggerByLineup Type NotExist[%s]", condition[0])
	}
	return false
}

func triggerMonsterByStar(monsterHero *MonsterHero, config *cfg.PassiveSkill) bool {
	condition := config.GetTriggerValue()
	return monsterHero.Star >= condition[0]
}

func triggerMonsterByHero(config *cfg.PassiveSkill, monsterHeros []*MonsterHero) bool {
	condition := config.GetTriggerValue()
	for _, heroID := range condition {
		isTrigger := false
		for _, f := range monsterHeros {
			if heroID == f.HeroID {
				isTrigger = true
				break
			}
		}
		if !isTrigger {
			return false
		}
	}
	return true
}

func triggerMonsterByLevel(monsterHero *MonsterHero, config *cfg.PassiveSkill) bool {
	condition := config.GetTriggerValue()
	return monsterHero.Level >= condition[0]
}

func triggerMonsterByAwake(monsterHero *MonsterHero, config *cfg.PassiveSkill) bool {
	condition := config.GetTriggerValue()
	return monsterHero.Awaken >= condition[0]
}

// ApplyMonsterPassiveSkill 仙人场外技能生效
func ApplyMonsterPassiveSkill(monsterHero *MonsterHero, list []int32) {
	if len(list) <= 0 {
		return
	}

	skillsMap := make(map[int32]bool)
	for _, v := range monsterHero.PassiveSkill {
		skillsMap[v] = true
	}

	for _, id := range list {
		// 去重校验
		if _, ok := skillsMap[id]; ok {
			continue
		}

		monsterHero.PassiveSkill = append(monsterHero.PassiveSkill, id)
		skillsMap[id] = true

		skillCfg := cfg.GetPassiveSkillConfig(id)
		if skillCfg == nil {
			continue // 配置有误
		}

		if skillCfg.EffectType != int32(cfg.EffectTypeOutBattle) {
			continue // 不是场外技能
		}

		outSkillCfg := cfg.OutBattleEffectOne(id)
		if outSkillCfg == nil {
			continue // 配置有误
		}

		switch outSkillCfg.Type {
		case int32(cfg.OutTypeAttr): // 直接转化为固定属性 1
			monsterHero.AttrSet.AddAttrSet(outSkillCfg.AttrSet)

		case int32(cfg.OutTypeReplaceNormal): // 普攻技能 2
			if outSkillCfg.Order > monsterHero.SkillIDOrder {
				monsterHero.SkillID = outSkillCfg.ReplaceSkill
				monsterHero.SkillIDOrder = outSkillCfg.Order
			}

		case int32(cfg.OutTypeReplaceAnger): // 怒气技能 3
			if outSkillCfg.Order > monsterHero.AngerSkillIDOrder {
				monsterHero.AngerSkillID = outSkillCfg.ReplaceSkill
				monsterHero.AngerSkillIDOrder = outSkillCfg.Order
			}

		case int32(cfg.OutTypeAddInBattleBuff): // 场内技能 4
			skillsMap[outSkillCfg.ReplaceSkill] = true

		case int32(cfg.OutTypeReplaceInBattleBuff): // 场内技能 5
			skillsMap[outSkillCfg.PassiveSkillToReplace] = false
			skillsMap[outSkillCfg.ReplaceSkill] = true
		}
	}
}
