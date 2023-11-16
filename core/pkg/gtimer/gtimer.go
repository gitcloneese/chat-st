package gtimer

import (
	"runtime/debug"
	"sync"
	"time"
)

var (
	once       sync.Once
	ins        *GTimer
	dumpLogger IDumpLogger
)

type GTimer struct {
	timers sync.Map
}

type IDumpLogger interface {
	Logf(fmt string, args ...interface{})
}

func Ins() *GTimer {
	once.Do(func() {
		ins = &GTimer{
			timers: sync.Map{},
		}
	})

	return ins
}

func Cancel(ch chan struct{}) {
	if nil == ch {
		return
	}

	ch <- struct{}{}
}

func Once(duration time.Duration, f func()) chan struct{} {
	return run(duration, duration, false, f, nil)
}

func OnceLocked(duration time.Duration, f func(), mu *sync.Mutex) chan struct{} {
	return run(duration, duration, false, f, mu)
}

func Forever(duration time.Duration, f func()) chan struct{} {
	return ForeverLocked(duration, f, nil)
}

func ForeverLocked(duration time.Duration, f func(), mu *sync.Mutex) chan struct{} {
	return run(duration, duration, true, f, mu)
}

func ForeverNow(duration time.Duration, f func()) chan struct{} {
	return ForeverNowLocked(duration, f, nil)
}

func ForeverNowLocked(duration time.Duration, f func(), mu *sync.Mutex) chan struct{} {
	return run(0, duration, true, f, mu)
}

func ForeverTime(durFirst, durRepeat time.Duration, f func()) chan struct{} {
	return ForeverTimeLocked(durFirst, durRepeat, f, nil)
}

func ForeverTimeLocked(durFirst, durRepeat time.Duration, f func(), mu *sync.Mutex) chan struct{} {
	return run(durFirst, durRepeat, true, f, mu)
}

func (t *GTimer) Cancel(id int) {
	Cancel(t.get(id))
}

func (t *GTimer) Once(id int, duration time.Duration, f func()) {
	_, exists := t.timers.Load(id)
	if exists {
		return
	}

	ch := Once(duration, f)
	t.timers.Store(id, ch)
}

func (t *GTimer) Forever(id int, duration time.Duration, f func()) {
	t.ForeverLocked(id, duration, f, nil)
}

func (t *GTimer) ForeverLocked(id int, duration time.Duration, f func(), mu *sync.Mutex) {
	_, exists := t.timers.Load(id)
	if exists {
		return
	}

	ch := ForeverLocked(duration, f, mu)
	t.timers.Store(id, ch)
}

func (t *GTimer) ForeverNow(id int, duration time.Duration, f func()) {
	t.ForeverNowLocked(id, duration, f, nil)
}

func (t *GTimer) ForeverNowLocked(id int, duration time.Duration, f func(), mu *sync.Mutex) {
	_, exists := t.timers.Load(id)
	if exists {
		return
	}

	ch := ForeverNowLocked(duration, f, mu)
	t.timers.Store(id, ch)
}

func (t *GTimer) ForeverTime(id int, durFirst, durRepeat time.Duration, f func()) {
	t.ForeverTimeLocked(id, durFirst, durRepeat, f, nil)
}

func (t *GTimer) ForeverTimeLocked(id int, durFirst, durRepeat time.Duration, f func(), mu *sync.Mutex) {
	_, exists := t.timers.Load(id)
	if exists {
		return
	}

	ch := ForeverTimeLocked(durFirst, durRepeat, f, mu)
	t.timers.Store(id, ch)
}

func (t *GTimer) get(id int) chan struct{} {
	timer, exists := t.timers.Load(id)
	if !exists {
		return nil
	}

	return timer.(chan struct{})
}

func run(durFirst, durRepeat time.Duration, repeated bool, f func(), mu * sync.Mutex) chan struct{} {
	//如果不加这个缓冲，在定时器回调函数中取消该定时器就会死锁
	ch := make(chan struct{},1)
	go func() {
		defer onRecover(func() {
			run(durFirst, durRepeat, repeated, f, mu)
		})
		timer := time.NewTimer(durFirst)
		for {
			select {
			case <-timer.C:
				func(){
					if mu != nil {
						mu.Lock()
						defer mu.Unlock()
					}

					f()
				}()
				if repeated {
					timer.Reset(durRepeat)
				} else {
					return
				}
			case <-ch:
				timer.Stop()
				return
			}
		}
	}()

	return ch
}

func onRecover(cb func())  {
	if err := recover();err != nil {
		if dumpLogger != nil {
			dumpLogger.Logf("Timer Recover,err:%s", err)
			dumpLogger.Logf("Trace Stack:%s", debug.Stack())
			cb()
		}
	}
}

