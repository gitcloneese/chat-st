package tools

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"xy3-proto/pkg/log"
)

func RequestNum() int64 {
	return atomic.LoadInt64(&RequestCount)
}

func connectCount() int64 {
	return atomic.LoadInt64(&chatConnectCount)
}

func msgCount() int64 {
	return atomic.LoadInt64(&receiveMsgCount)
}

func ErrNum() int64 {
	return atomic.LoadInt64(&ErrCount)
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
	errCodes = new(sync.Map)
	fmt.Printf("|||执行完毕:%20v| 总请求次数:%5v | 成功:%4v | 失败:%5v | 用时:%10.4f | qps:%10.4f | 错误码:%v |||\n", name, allRequestNum, success, errNum, latency, float64(allRequestNum)/latency, strings.Join(errCode, ","))
	time.Sleep(time.Second)
}

// RunWithLogTick
// 异步执行 每800毫秒打印一次当前任务执行次数
// endLogNum 结束打印的次数
func RunWithLogTick(name string, f func(), endLogNum ...int) {
	log.Info("开始执行:%v !!!", name)
	errCodes = new(sync.Map)
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
	signal.Notify(osStop, os.Interrupt, os.Kill)
	printLog := func() {
		allRequestNum := RequestNum() - requestCountStart
		errNum := ErrNum() - errStart
		success := allRequestNum - errNum
		latency := time.Since(startTime).Seconds()
		errCode := make([]string, 0, 10)
		errCodes.Range(func(key, value interface{}) bool {
			errCode = append(errCode, fmt.Sprintf("%v", key.(int)))
			return true
		})
		fmt.Printf("|||执行:%20v| 总请求次数:%5v | 成功:%4v | 失败:%5v | ws数量:%5v | 用时:%10.4f | qps:%10.4f | 平均延迟:%7.4f | 错误码:%v |||\n", name, allRequestNum, success, errNum, connectCount(), latency, float64(allRequestNum)/latency, latency/float64(allRequestNum), strings.Join(errCode, ","))
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

// RunReceiveMsgWithLogTick
// 打印接收消日志
func RunReceiveMsgWithLogTick(name string, f func()) {
	log.Info("开始执行:%v !!!", name)
	errCodes = new(sync.Map)
	now := time.Now()
	// 异步执行task
	go f()
	tickReceiveMsgLog(name, now)
}

// 每隔1秒钟 做一次日志打印
func tickReceiveMsgLog(name string, startTime time.Time) {
	tick := time.NewTicker(800 * time.Millisecond)
	defer tick.Stop()
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, os.Kill, syscall.SIGTERM)
	printLog := func() {
		latency := time.Since(startTime).Seconds()
		errCode := make([]string, 0, 10)
		errCodes.Range(func(key, value interface{}) bool {
			errCode = append(errCode, fmt.Sprintf("%v", key.(int)))
			return true
		})
		fmt.Printf("|||执行:%20v| ws长连接数量:%5v | 收到消息数:%5v | 用时:%10.4f |||\n", name, connectCount(), msgCount(), latency)
	}

	for {
		select {
		case <-tick.C:
			printLog()
		case <-stopCh:
			// 额外再打印一次
			printLog()
			return
		}
	}
}
