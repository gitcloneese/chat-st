//file 		: function.go
//date 		: 2020/11/28 12:15
//author 	: looyer@sina.com
//purpose 	: 辅助函数库

package util

import (
	"container/list"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// MakeGUID 通过K标识生成存库永久唯一的guid
func MakeGUID(k int) int64 {
	return 0
}

// TempGUID 生成本进程单次运行的唯一guid
func TempGUID() int64 {
	return 0
}

func GetCpuPercent() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0]
}

func GetMemPercent() float64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.UsedPercent
}

func GetServiceName(ID string) string {
	index := strings.LastIndex(ID, "-")
	return ID[:index]
}

// MaxInt .
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// MinInt .
func MinInt(a, b int) int {
	if a > b {
		return b
	}
	return a
}

// MinInt32 .
func MinInt32(a, b int32) int32 {
	if a > b {
		return b
	}
	return a
}

// MaxInt64 .
func MaxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// MinInt64 .
func MinInt64(a, b int64) int64 {
	if a > b {
		return b
	}
	return a
}

func AtomicAdd(addr *int32) int32 {
	atomic.CompareAndSwapInt32(addr, 0x7fffffff, 1)
	return atomic.AddInt32(addr, 1)
}

func MakeError(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

func ListForEach(list list.List, visitor func(data interface{})) {
	for iter := list.Front(); iter != nil; iter = iter.Next() {
		visitor(iter.Value)
	}
}

// Foreach 遍历容器元素执行visitor container必须是Map/Slice/Array
func ForEach(container interface{}, visitor func(k interface{}, v interface{})) {
	if !AssertContainer(container) {
		return
	}
	cValue := reflect.ValueOf(container)
	switch cValue.Type().Kind() {
	case reflect.Map:
		iter := cValue.MapRange()
		for iter.Next() {
			visitor(iter.Key().Interface(), iter.Value().Interface())
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < cValue.Len(); i++ {
			visitor(i, cValue.Index(i).Interface())
		}
	}
}

// AssertContainer 判断传入参数是否Map/Slice/Array
func AssertContainer(container interface{}) bool {
	if container == nil {
		return false
	}
	kind := reflect.ValueOf(container).Type().Kind()
	if kind != reflect.Map && kind != reflect.Slice && kind != reflect.Array {
		return false
	}
	return true
}

// AssertArray 判断传入参数是否Slice/Array
func AssertArray(container interface{}) bool {
	kind := reflect.ValueOf(container).Type().Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return false
	}
	return true
}

type IWheelItem interface {
	GetWeight() int
}

// GetRandomElement 使用轮盘算法随机从容器中获取一个元素 元素必须实现IWheelItem接口
func GetRandomElement(container interface{}) interface{} {
	if !AssertContainer(container) {
		return nil
	}
	cValue := reflect.ValueOf(container)
	cType := cValue.Type()
	if cValue.Len() == 0 {
		return nil
	}
	eleType := cType.Elem()
	_, succ := reflect.New(eleType).Interface().(IWheelItem)
	if !succ {
		return nil
	}
	total := 0
	ForEach(container, func(k interface{}, v interface{}) {
		total += v.(IWheelItem).GetWeight()
	})
	randNum := RandomClosed(0, total-1)
	tmpNum := 0
	var ret interface{} = nil
	ForEach(container, func(k interface{}, v interface{}) {
		tmpNum += v.(IWheelItem).GetWeight()
		if tmpNum > randNum && ret == nil {
			ret = v
		}
	})
	return ret
}

// M个里取N个
func GetNFromM(container interface{}, num int) []interface{} {
	ret := []interface{}{}
	if !AssertContainer(container) {
		return ret
	}
	cValue := reflect.ValueOf(container)
	if cValue.Len() == 0 {
		return ret
	}
	index := 0
	//蓄水池
	ForEach(container, func(k interface{}, v interface{}) {
		if index < num {
			ret = append(ret, v)
		} else {
			if RandCheck(num, index+1) {
				swapIndex := RandomClosed(0, num-1)
				ret[swapIndex] = v
			}
		}
		index++
	})
	return ret
}

// Filter 筛选Map/Slice/Array
func Filter(container interface{}, filter func(k interface{}, v interface{}) bool) []interface{} {
	ret := []interface{}{}
	if !AssertContainer(container) {
		return ret
	}
	cValue := reflect.ValueOf(container)
	if cValue.Len() == 0 {
		return ret
	}
	ForEach(container, func(k interface{}, v interface{}) {
		if filter(k, v) {
			ret = append(ret, v)
		}
	})
	return ret
}

// find 筛选Map/Slice/Array
func Finder(container interface{}, filter func(v interface{}) bool) interface{} {
	if !AssertContainer(container) {
		return nil
	}

	cValue := reflect.ValueOf(container)
	if cValue.Len() == 0 {
		return nil
	}

	switch cValue.Type().Kind() {
	case reflect.Map:
		iter := cValue.MapRange()
		for iter.Next() {
			if filter(iter.Value()) {
				return iter.Value().Interface()
			}
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < cValue.Len(); i++ {
			if filter(cValue.Index(i)) {
				return cValue.Index(i).Interface()
			}
		}
	}

	return nil
}

func Shuffle(array interface{}) {
	if !AssertArray(array) {
		return
	}
	cValue := reflect.ValueOf(array)
	length := cValue.Len()
	for i := 0; i < length-1; i++ {
		randNum := RandomClosed(i, length-1)
		if randNum != i {
			tmp := cValue.Index(i).Interface()
			cValue.Index(i).Set(reflect.ValueOf(cValue.Index(randNum).Interface()))
			cValue.Index(randNum).Set(reflect.ValueOf(tmp))
		}
	}
}

func GetInt32Max() int32 {
	return 0x7fffffff
}

func ExtractConfigContentsFromString(input string) []string {
	//regular expression to match content inside curly braces
	regex := regexp.MustCompile(`\{([^}]+)\}`)
	matches := regex.FindAllStringSubmatch(input, -1)
	var contents []string
	// Extract the content from the matches and populate the array
	for _, match := range matches {
		if len(match) >= 2 {
			content := match[1]
			contents = append(contents, content)
		}
	}
	return contents
}
