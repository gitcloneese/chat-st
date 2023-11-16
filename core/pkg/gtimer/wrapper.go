package gtimer

import (
	"sync"
	"sync/atomic"
	"time"
)

//保留10000以下的timer id作为全局id
const INIT_TIMER_ID = 10000

const INT32_MAX = 0x7fffffff

var (
	genID    int32
	timerMap sync.Map
	mutex    sync.Mutex
)

func init() {
	genID = INIT_TIMER_ID
}

//DoOnceLocked 执行一次
//传入timerID小于等于0时生成新timerID
//返回timerID
func DoOnce(delay time.Duration, callback func()) int32 {
	mutex.Lock()
	defer mutex.Unlock()
	timerID := genTimerID()
	ch := Once(delay, func() {
		delTimer(timerID)
		callback()
	})
	addTimer(timerID, ch)
	return timerID
}

//DoOnceLocked 执行一次带加锁
//传入timerID小于等于0时生成新timerID
//返回timerID
func DoOnceLocked(delay time.Duration, callback func(), lock *sync.Mutex) int32 {
	mutex.Lock()
	defer mutex.Unlock()
	timerID := genTimerID()
	ch := OnceLocked(delay, func() {
		delTimer(timerID)
		callback()
	}, lock)
	addTimer(timerID, ch)
	return timerID
}

//Loop 重复执行
//传入timerID小于等于0时生成新timerID
//返回timerID
func Loop(timerID int32, duration time.Duration, callback func()) int32 {
	mutex.Lock()
	defer mutex.Unlock()
	if timerID <= 0 {
		timerID = genTimerID()
	} else {
		if timerHandler := getTimer(timerID); timerHandler != nil {
			timerHandler <- struct{}{}
			delTimer(timerID)
		}
	}
	ch := Forever(duration, callback)
	addTimer(timerID, ch)
	return timerID
}

//LoopLock 重复执行带加锁
//传入timerID小于等于0时生成新timerID
//返回timerID
func LoopLock(timerID int32, duration time.Duration, callback func(), lock *sync.Mutex) int32 {
	mutex.Lock()
	defer mutex.Unlock()
	if timerID <= 0 {
		timerID = genTimerID()
	} else {
		if timerHandler := getTimer(timerID); timerHandler != nil {
			timerHandler <- struct{}{}
			delTimer(timerID)
		}
	}
	ch := ForeverLocked(duration, callback, lock)
	addTimer(timerID, ch)
	return timerID
}

//CancelTimer
func CancelTimer(timerID int32) bool {
	mutex.Lock()
	defer mutex.Unlock()
	if timerHandler := getTimer(timerID); timerHandler == nil {
		return false
	} else {
		timerHandler <- struct{}{}
		delTimer(timerID)
		return true
	}
}

//SetDumpLogger 设置发生recover时调用的logger
func SetDumpLogger(l IDumpLogger)  {
	dumpLogger = l
}

func genTimerID() int32 {
	return atomicAdd(&genID)
}

func addTimer(timerID int32, ch chan struct{}) {
	timerMap.Store(timerID, ch)
}

func getTimer(timerID int32) chan struct{} {
	timerHandler, ok := timerMap.Load(timerID)
	if !ok {
		return nil
	}
	return timerHandler.(chan struct{})
}

func delTimer(timerID int32) {
	timerMap.Delete(timerID)
}

func atomicAdd(addr *int32) int32 {
	atomic.CompareAndSwapInt32(addr, INT32_MAX, 1)
	return atomic.AddInt32(addr, 1)
}