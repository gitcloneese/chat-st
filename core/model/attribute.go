package model

import (
	"fmt"
	"xy3-proto/pkg/conf/env"
	"xy3-proto/pkg/log"

	battle "xy3-proto/battle"
)

const (
	Permillage = 1000 // Permillage 千分比值
)

type AttributeSet struct {
	attrs map[battle.AttributeKey]int64
}

// NewAttributeSet 构造一个空的属性集合
func NewAttributeSet() *AttributeSet {
	p := &AttributeSet{attrs: make(map[battle.AttributeKey]int64)}
	return p
}

// MakeAttributeSet 基于数组构造一个属性集合
func MakeAttributeSet(attr [][2]int64) *AttributeSet {
	p := &AttributeSet{attrs: make(map[battle.AttributeKey]int64)}
	for _, v := range attr {
		p.SetAttr(battle.AttributeKey(v[0]), v[1])
	}
	return p
}

func MakeAttributeSetFromBattleAttribute(attr []*battle.Attrbute) *AttributeSet {
	p := &AttributeSet{attrs: make(map[battle.AttributeKey]int64)}
	for _, v := range attr {
		p.SetAttr(v.Key, int64(v.Value))
	}
	return p
}

func MakeAttributeSetFromMap(attrMap map[int32]int64) *AttributeSet {
	p := &AttributeSet{attrs: make(map[battle.AttributeKey]int64)}
	for key, value := range attrMap {
		p.SetAttr(battle.AttributeKey(key), value)
	}
	return p
}

// Clear 清理重置
func (p *AttributeSet) Clear() {
	p.attrs = make(map[battle.AttributeKey]int64)
}

func (p *AttributeSet) Zero() *AttributeSet {
	for k := range p.attrs {
		p.attrs[k] = 0
	}
	return p
}

func (p *AttributeSet) Data() map[battle.AttributeKey]int64 {
	return p.attrs
}

// Copy 克隆一个属性集合镜像
func (p *AttributeSet) Copy() *AttributeSet {
	pCopy := NewAttributeSet()
	for k, v := range p.attrs {
		pCopy.attrs[k] = v
	}
	return pCopy
}

// 生成只有仙人属性的属性对象
func (p *AttributeSet) ToSetHero() {
	for k := range p.attrs {
		if battle.AttributeKey_ACTORNONE < k && k < battle.AttributeKey(100) {
			p.attrs[k] = 0
		}
	}
}

// 生成只有主角属性的属性对象
func (p *AttributeSet) ToSetActor() {
	for k := range p.attrs {
		if k < battle.AttributeKey_ACTORNONE || battle.AttributeKey(100) < k {
			p.attrs[k] = 0
		}
	}
}

// AddActorAttr
// 添加主角属性
func (p *AttributeSet) AddActorAttr(t *AttributeSet) {
	p.AddAttrSet(t)
	p.ToSetActor()
}

// Full 使用另一个属性集合覆盖填充
func (p *AttributeSet) Full(other *AttributeSet) {
	p.Clear()
	for k, v := range other.attrs {
		p.attrs[k] = v
	}
}

func (p *AttributeSet) Format() []*battle.Attrbute {
	result := make([]*battle.Attrbute, 0)

	for i := battle.AttributeKey_HP; i < battle.AttributeKey_MAX; i++ {
		if v, ok := p.attrs[i]; ok {
			attr := &battle.Attrbute{
				Key:   i,
				Value: uint64(v),
			}
			result = append(result, attr)
		}
		// } else {
		// 	attr := &battle.Attrbute{
		// 		Key:   i,
		// 		Value: uint64(0),
		// 	}
		// 	result = append(result, attr)
		// }
	}

	return result
}

func (p *AttributeSet) ToBattleAttribute(list []*battle.Attrbute) []*battle.Attrbute {
	for k, v := range p.attrs {
		if v == 0 {
			continue
		}
		list = append(list, &battle.Attrbute{Key: k, Value: uint64(v)})
	}
	return list
}

func (p *AttributeSet) ToMap() (m map[int32]uint64) {
	m = make(map[int32]uint64)
	for k, v := range p.attrs {
		if v == 0 {
			continue
		}
		m[int32(k)] = uint64(v)
	}
	return
}

// new - old
func (p *AttributeSet) DiffToMap(obj *AttributeSet) (m map[int32]uint64) {
	m = make(map[int32]uint64)
	for k, v := range p.attrs {
		if v != obj.GetAttr(k) {
			m[int32(k)] = uint64(v)
		}
	}
	for k := range obj.attrs {
		if _, ok := p.attrs[k]; !ok {
			m[int32(k)] = 0
		}
	}
	return
}

func (p *AttributeSet) Diff(oldAttr *AttributeSet) []*battle.Attrbute {
	result := make([]*battle.Attrbute, 0)
	for k, v := range p.attrs {
		if v != oldAttr.GetAttr(k) {
			df := &battle.Attrbute{
				Key:   k,
				Value: uint64(v),
			}
			result = append(result, df)
		}
	}
	return result
}

func (p *AttributeSet) FormatStr() string {
	if env.Env == env.STEnv {
		return "cost much memory, only show in dev env"
	}

	m := make(map[int32]int64)
	for k, v := range p.attrs {
		if v > 0 {
			m[int32(k)] = v
		}
	}

	return fmt.Sprintf("%v", m)
}

func (p *AttributeSet) GetAttr(key battle.AttributeKey) int64 {
	if v, ok := p.attrs[key]; ok {
		return v
	}
	return 0
}

// SetAttr 直接设置单个属性 .
func (p *AttributeSet) SetAttr(attrDef battle.AttributeKey, value int64) {
	if attrDef <= battle.AttributeKey_NONE || attrDef >= battle.AttributeKey_MAX {
		return
	}

	if attrDef >= battle.AttributeKey_ATKRATE || attrDef <= battle.AttributeKey_MDEFRATE {
		p.attrs[attrDef] = value // 属性百分比可以为负数
	} else {
		if value > 0 {
			p.attrs[attrDef] = value
		} else {
			p.attrs[attrDef] = 0
		}
	}
}

// AddAttr 叠加增量
func (p *AttributeSet) AddAttr(attrDef battle.AttributeKey, value int64) (actualValue int64) {
	if attrDef <= battle.AttributeKey_NONE || attrDef >= battle.AttributeKey_MAX {
		return
	}
	if value == 0 {
		return 0
	}

	if value > 0 {
		if _, ok := p.attrs[attrDef]; !ok {
			p.attrs[attrDef] = 0
		}
		p.attrs[attrDef] += value
		actualValue = value
	} else {
		log.Warn("AddAttr: key:%v name:%v, value:%v 为负数 !!!", attrDef, battle.AttributeKey_name[int32(attrDef)], value)
		actualValue = value
		if p.attrs[attrDef] < value {
			if attrDef >= battle.AttributeKey_ATKRATE && attrDef <= battle.AttributeKey_MDEFRATE { // CHANGELOG  || =>> &&
				p.attrs[attrDef] -= value // 属性百分比可以为负数
			} else {
				actualValue = p.attrs[attrDef]
				p.attrs[attrDef] = 0
			}
		} else {
			p.attrs[attrDef] -= value
		}
	}

	return actualValue
}

func (p *AttributeSet) AddOne(k battle.AttributeKey, v int64) {
	p.attrs[k] += v
}

// AddAttrSet 直接加属性表
func (p *AttributeSet) AddAttrSet(pOther *AttributeSet) {
	if pOther == nil {
		return
	}
	for k, v := range pOther.attrs {
		p.AddAttr(k, v)
	}
}

// Scale 缩放属性
func (p *AttributeSet) Scale(ra int64, rb int64) {
	if ra <= 0 || rb <= 0 {
		return
	}

	for k, v := range p.attrs {
		p.attrs[k] = v * ra / rb
	}
}

// FinalScale 缩放属性(固定属性*属性百分比)
func (p *AttributeSet) FinalScale() {
	pairs := make(map[battle.AttributeKey]battle.AttributeKey)

	pairs[battle.AttributeKey_HP] = battle.AttributeKey_HPRATE
	pairs[battle.AttributeKey_ATK] = battle.AttributeKey_ATKRATE
	pairs[battle.AttributeKey_PDEF] = battle.AttributeKey_PDEFRATE
	pairs[battle.AttributeKey_MDEF] = battle.AttributeKey_MDEFRATE

	pairs[battle.AttributeKey_ACTORATTACK] = battle.AttributeKey_ACTORATKRATE
	pairs[battle.AttributeKey_ACTORHP] = battle.AttributeKey_ACTORHPRATE
	pairs[battle.AttributeKey_ACTORPDEF] = battle.AttributeKey_ACTORPDEFRATE
	pairs[battle.AttributeKey_ACTORMDEF] = battle.AttributeKey_ACTORMDEFRATE

	for k, v := range pairs {
		val := p.GetAttr(k)
		rate := p.GetAttr(v)
		if val == 0 || rate == 0 {
			continue
		}
		num := val*rate/Permillage + val
		p.SetAttr(k, num)
	}
}

func (p *AttributeSet) AddSlice(s [][2]int32) {
	if s == nil {
		return
	}
	for _, v := range s {
		if int32(battle.AttributeKey_NONE) < v[0] && v[0] < int32(battle.AttributeKey_MAX) {
			p.AddAttr(battle.AttributeKey(v[0]), int64(v[1]))
		}
	}
}
