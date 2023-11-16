package util

import (
	"math/rand"
	"time"
)

//nolint:staticcheck
func init() {
	rand.Seed(time.Now().Unix())
}

// RandomClosed get random num [min,max]
func RandomClosed(min int, max int) int {
	if min > max {
		min, max = max, min
	}
	return rand.Intn(max-min+1) + min
}

func RandClosed(min int32, max int32) int32 {
	if min > max {
		min, max = max, min
	}
	return rand.Int31n(max-min+1) + min
}

// RandomSemiClosed get random num [min,max)
func RandomSemiClosed(min int, max int) int {
	if min == max {
		return min
	}

	if min > max {
		min, max = max, min
	}
	return rand.Intn(max-min) + min
}

func RandSemiClosed(min int32, max int32) int32 {
	if min == max {
		return min
	}

	if min > max {
		min, max = max, min
	}
	return rand.Int31n(max-min) + min
}

func RandCheck(ratio int, base int) bool {
	if base <= 0 {
		return false
	}
	//[0,base)
	randNum := rand.Intn(base)
	return ratio > randNum
}

func RandOk(ratio int32, base int32) bool {
	if base <= 0 {
		return false
	}
	//[0,base)
	return ratio > rand.Int31n(base)
}

// Rand [0, max)
func Rand(max int) int {
	if max <= 0 {
		return 0
	}
	return rand.Intn(max)
}

// RandRange [min, max)
func RandRange(min int, max int) int {
	if min < 0 || max <= 0 {
		return 0
	}
	if min >= max {
		return min
	}
	return min + rand.Intn(max-min)
}

// RandRange int64版本 [min, max)
func RandRange64(min int64, max int64) int64 {
	if min < 0 || max <= 0 {
		return 0
	}
	if min >= max {
		return min
	}
	return rand.Int63n(max-min) + min
}

// 根据权重随机结果
func RandomIndexFromWeight(weightList []int64) int {
	total := int64(0)
	for _, v := range weightList {
		total += v
	}

	if total == 0 {
		return 0
	}

	randNum := rand.Int63n(total) + 1
	for k, v := range weightList {
		if randNum > v {
			randNum -= v
		} else {
			return k
		}
	}

	return 0
}

// 按权重随机 输入 [权重1,权重2,权重3] 返回随机权重对应的索引值
func RandomItem(tbl *[]uint32) int {
	var result int = -1
	var length int = len(*tbl)
	if length <= 0 {
		return result
	}
	var sum uint32 = 0
	var vtbl []uint32 = make([]uint32, length)
	for i, value := range *tbl {
		sum += value
		vtbl[i] = sum
	}
	if sum <= 0 {
		return result
	}

	r := uint32(rand.Int31n(int32(sum)))
	for i, value := range vtbl {
		if r < value {
			result = i
			break
		}
	}
	return result
}

// 按权重随机 M个随机N个 输入 [权重1,权重2,权重3],个数 返回随机N个权重对应的索引值
func RandomItemEx(tbl *[]uint32, num uint32) (result []uint32) {
	result = make([]uint32, 0)
	m := make(map[uint32]bool)
	for i := 0; i < int(num); i++ {
		var id_list []uint32
		var weight_list []uint32
		for k, v := range *tbl {
			if _, ok := m[uint32(k)]; ok {
				continue
			}
			id_list = append(id_list, uint32(k))
			weight_list = append(weight_list, v)
		}
		if len(id_list) <= 0 {
			return
		}
		random_index := RandomItem(&weight_list)
		if random_index == -1 {
			return // 异常情况
		}
		result = append(result, id_list[random_index])
		m[id_list[random_index]] = true
	}
	return
}

// 按权重随机 输入 [[id1,权重1],[id2,权重2],[id3,权重3] 返回随机权重对应的id
func RandomItemID(tbl *[][2]uint32) int32 {
	var id_list []uint32
	var weight_list []uint32
	for _, v := range *tbl {
		id_list = append(id_list, v[0])
		weight_list = append(weight_list, v[1])
	}
	random_index := RandomItem(&weight_list)
	if random_index != -1 {
		return int32(id_list[random_index])
	}
	return -1
}

// 按权重随机 M个随机N个 随机N个 输入 [[id1,权重1],[id2,权重2],[id3,权重3] 返回N个id(不可重复)
func RandomItemIDEx(tbl *[][2]uint32, num uint32) (result []uint32) {
	result = make([]uint32, 0)
	m := make(map[uint32]bool)
	for i := 0; i < int(num); i++ {
		var id_list []uint32
		var weight_list []uint32
		for _, v := range *tbl {
			if _, ok := m[v[0]]; ok {
				continue
			}
			id_list = append(id_list, v[0])
			weight_list = append(weight_list, v[1])
		}
		if len(id_list) <= 0 {
			return
		}
		random_index := RandomItem(&weight_list)
		if random_index == -1 {
			return // 异常情况
		}
		result = append(result, id_list[random_index])
		m[id_list[random_index]] = true
	}
	return
}

// 按权重随机 M个随机N个 随机N个 输入 [[id1,权重1],[id2,权重2],[id3,权重3] 返回N个id(可重复)
func RandomItemIDEx2(tbl *[][2]uint32, num uint32) (result []uint32) {
	result = make([]uint32, 0)
	var id_list []uint32
	var weight_list []uint32
	for _, v := range *tbl {
		id_list = append(id_list, v[0])
		weight_list = append(weight_list, v[1])
	}

	for i := 0; i < int(num); i++ {
		random_index := RandomItem(&weight_list)
		if random_index == -1 {
			return // 异常情况
		}
		result = append(result, id_list[random_index])
	}
	return
}

func NewRandomUtil(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}

// 按权重随机 输入 [权重1,权重2,权重3] 返回随机权重对应的索引值
func RandomItemWithRand(rnd *rand.Rand, tbl *[]uint32) int {
	var result int = -1
	var length int = len(*tbl)
	if length <= 0 {
		return result
	}
	var sum uint32 = 0
	var vtbl []uint32 = make([]uint32, length)
	for i, value := range *tbl {
		sum += value
		vtbl[i] = sum
	}
	if sum <= 0 {
		return result
	}

	r := uint32(rnd.Int31n(int32(sum)))
	for i, value := range vtbl {
		if r < value {
			result = i
			break
		}
	}
	return result
}

// 按权重随机 M个随机N个 输入 [权重1,权重2,权重3],个数 返回随机N个权重对应的索引值
func RandomItemExWithRand(rnd *rand.Rand, tbl *[]uint32, num uint32) (result []uint32) {
	result = make([]uint32, 0)
	m := make(map[uint32]bool)
	for i := 0; i < int(num); i++ {
		var id_list []uint32
		var weight_list []uint32
		for k, v := range *tbl {
			if _, ok := m[uint32(k)]; ok {
				continue
			}
			id_list = append(id_list, uint32(k))
			weight_list = append(weight_list, v)
		}
		if len(id_list) <= 0 {
			return
		}
		random_index := RandomItemWithRand(rnd, &weight_list)
		if random_index == -1 {
			return // 异常情况
		}
		result = append(result, id_list[random_index])
		m[id_list[random_index]] = true
	}
	return
}

// 按权重随机 M个随机N个 随机N个 输入 [[id1,权重1],[id2,权重2],[id3,权重3] 返回N个id(不可重复)
func RandomItemIDExWithRand(rnd *rand.Rand, tbl *[][2]uint32, num uint32) (result []uint32) {
	result = make([]uint32, 0)
	m := make(map[uint32]bool)
	for i := 0; i < int(num); i++ {
		var id_list []uint32
		var weight_list []uint32
		for _, v := range *tbl {
			if _, ok := m[v[0]]; ok {
				continue
			}
			id_list = append(id_list, v[0])
			weight_list = append(weight_list, v[1])
		}
		if len(id_list) <= 0 {
			return
		}
		random_index := RandomItemWithRand(rnd, &weight_list)
		if random_index == -1 {
			return // 异常情况
		}
		result = append(result, id_list[random_index])
		m[id_list[random_index]] = true
	}
	return
}

// UniqueRandomInt32 随机`num`个范围内[minValue,maxValue]不重复的数字
func UniqueRandomInt32(num, minValue, maxValue int32) []int32 {
	length := maxValue - minValue + 1
	seed := make([]int32, length)
	for i := int32(0); i < length; i++ {
		seed[i] = minValue + i
	}

	if num > length {
		num = length
	}

	ranArr := make([]int32, num)
	for i := int32(0); i < num; i++ {
		j := rand.Intn(int(length - i))
		ranArr[i] = seed[j]
		seed[j] = seed[length-i-1]
	}
	return ranArr
}

// RandomHit 根据概率配置返回随机命中的id
// [[id,概率],...]
func RandomHit(base int32, randomMap [][]int32) int32 {
	rate := rand.Int31n(base)
	total := int32(0)
	for _, r := range randomMap {
		if r[1] == 0 {
			continue
		}
		total += r[1]
		if total >= rate {
			return r[0]
		}
	}
	return -1
}
