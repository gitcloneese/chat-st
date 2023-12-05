package tools

import (
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"xy3-proto/pkg/log"
)

func RequestNum() int64 {
	return atomic.LoadInt64(&RequestCount)
}

func ErrNum() int64 {
	return atomic.LoadInt64(&ErrCount)
}

func allLatency() float64 {
	AllLatencyLock.RLock()
	defer AllLatencyLock.RUnlock()
	return AllLatency
}

// 最小请求延迟时间
func minLatency() float64 {
	MinLatencyLock.RLock()
	defer MinLatencyLock.RUnlock()
	return MinLatency
}

func setMinLatency(latency float64) {
	MinLatencyLock.Lock()
	defer MinLatencyLock.Unlock()
	if MinLatency <= 0 && latency > 0 {
		MinLatency = latency
		return
	}

	if latency < MinLatency {
		MinLatency = latency
	}
}

func setAllLatency(latency float64) {
	AllLatencyLock.Lock()
	defer AllLatencyLock.Unlock()
	AllLatency += latency
}

func resetMinRequestLatency() {
	MinLatencyLock.Lock()
	defer MinLatencyLock.Unlock()
	MinLatency = 0
}
func resetAllRequestLatency() {
	AllLatencyLock.Lock()
	defer AllLatencyLock.Unlock()
	AllLatency = 0
}

// 最大请求延迟
func maxLatency() float64 {
	MaxLatencyLock.RLock()
	defer MaxLatencyLock.RUnlock()
	return MaxLatency
}

func setMaxLatency(latency float64) {
	MaxLatencyLock.Lock()
	defer MaxLatencyLock.Unlock()
	if latency > MaxLatency {
		MaxLatency = latency
	}
}

func resetMaxRequestLatency() {
	MaxLatencyLock.Lock()
	defer MaxLatencyLock.Unlock()
	MaxLatency = 0
}

// ResetLatency
// 重置延迟记录
func ResetLatency() {
	resetMinRequestLatency()
	resetMaxRequestLatency()
	resetAllRequestLatency()
}

// SetLatency
// 设置延迟时间 , 这是一个defer方法
func SetLatency() func(time.Time) {
	return func(startTime time.Time) {
		latency := time.Since(startTime).Seconds()
		setMinLatency(latency)
		setMaxLatency(latency)
		setAllLatency(latency)
	}
}

func Error(format string, args ...interface{}) {
	if Debug {
		if len(args) > 0 {
			if strings.HasSuffix(format, "\n") {
				fmt.Printf(format, args...)
			} else {
				fmt.Printf(format+"\n", args...)
			}
		} else {
			fmt.Println(format)
		}
	}
}

func RunWithLog(name string, f func()) {
	log.Info("开始执行:%v !!!", name)
	now := time.Now()
	errCodes = new(sync.Map)
	request1 := RequestNum()
	errNum1 := ErrNum()
	f()
	allRequestNum := RequestNum() - request1
	errNum := ErrNum() - errNum1
	success := allRequestNum - errNum
	latency := time.Since(now).Seconds()
	errCode := make([]string, 0, 10)
	errCodes.Range(func(key, value interface{}) bool {
		errCode = append(errCode, fmt.Sprintf("%v", key.(int)))
		return true
	})
	fmt.Printf("|||执行完毕:%20v| 总请求次数:%5v | 成功:%4v | 失败:%5v | 用时:%10.4f | qps:%10.4f | 错误码:%v |||\n", name, allRequestNum, success, errNum, latency, float64(allRequestNum)/latency, strings.Join(errCode, ","))
	time.Sleep(time.Second)
}

// RunWithLogTick
// 异步执行 每800毫秒打印一次当前任务执行次数
// endLogNum 结束打印的次数
func RunWithLogTick(name string, f func(), apiPath ...string) {
	fmt.Println("------------------------------------------")
	if len(apiPath) > 0 {
		path := apiPath[0]
		fmt.Printf("开始执行:%v 接口地址:%v 线程数:%v!!!\n", name, path, C)
	} else {
		fmt.Printf("开始执行:%v 线程数:%v!!!\n", name, C)
	}
	// 清除错误码
	errCodes = new(sync.Map)
	// 重置延迟时间
	ResetLatency()
	now := time.Now()
	request := RequestNum()
	errNum := ErrNum()
	stopCh := make(chan os.Signal, 1)
	// 异步执行task
	go func(ch chan os.Signal, f func()) {
		defer close(ch)
		f()
	}(stopCh, f)
	tickLog(name, now, errNum, request, stopCh)
}

// 每隔1秒钟 做一次日志打印
func tickLog(name string, startTime time.Time, errStart, requestCountStart int64, stop chan os.Signal) {
	tick := time.NewTicker(800 * time.Millisecond)
	defer tick.Stop()
	osStop := make(chan os.Signal, 1)
	defer close(osStop)
	signal.Notify(osStop, os.Interrupt, os.Kill)
	printLog := func() {
		allRequestNum := RequestNum() - requestCountStart
		errNum := ErrNum() - errStart
		success := allRequestNum - errNum
		latency := time.Since(startTime).Seconds()
		errCode := make([]int, 0, 10)
		errCodes.Range(func(key, value interface{}) bool {
			errCode = append(errCode, key.(int))
			return true
		})
		sort.Ints(errCode)
		codes := ""
		for _, v := range errCode {
			if codes == "" {
				codes = fmt.Sprintf("%v", v)
			} else {
				codes = fmt.Sprintf("%v,%v", codes, v)
			}
		}
		fmt.Printf("|||执行:%20v| 总请求次数:%7v | 成功:%7v | 失败:%7v | 用时:%10.4f | qps:%10.4f | 平均延迟:%7.4f | 最大延迟:%7.4v | 最小延迟:%7.4v | 错误码:%v |||\n", name, allRequestNum, success, errNum, latency, float64(allRequestNum)/latency, allLatency()/float64(allRequestNum), maxLatency(), minLatency(), codes)
	}

	for {
		select {
		case <-tick.C:
			printLog()
		case <-osStop:
			// 额外再打印一次
			printLog()
			return
		case <-stop:
			// 额外再打印一次
			printLog()
			return
		}
	}
}
