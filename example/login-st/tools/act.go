package tools

import (
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

func Error(format string, args ...interface{}) {
	if Debug {
		log.Error(format, args...)
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
	log.Info("|||执行完毕:%20v| 总请求次数:%5v | 成功:%4v | 失败:%5v | 用时:%10.4f | qps:%10.4f |||", name, allRequestNum, success, errNum, latency, float64(allRequestNum)/latency)
	time.Sleep(time.Second)
}
