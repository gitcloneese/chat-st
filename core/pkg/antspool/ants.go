package antspool

import (
	"errors"
	"sync/atomic"
	"time"
	"xy3-proto/pkg/log"

	ants "github.com/panjf2000/ants/v2"
)

type TaskPool struct {
	p               *ants.Pool
	waitTaskQueueCh chan func()
	stopCh          chan bool
	isDisPatch      int32
}

func InitAntPool(size int) (tp *TaskPool) {
	if size <= 0 {
		size = 100
	}
	tp = &TaskPool{
		isDisPatch: 1,
	}
	var err error
	// ants will pre-malloc the whole capacity of pool when you invoke this function
	tp.p, err = ants.NewPool(size,
		ants.WithPreAlloc(true),
		ants.WithNonblocking(true),
		ants.WithExpiryDuration(time.Second),
		ants.WithLogger(&AntsLogger{}))
	if err != nil {
		panic(err)
	}

	tp.waitTaskQueueCh = make(chan func(), 2000*size)
	tp.stopCh = make(chan bool, 1)

	go tp._tDispatchWork()
	return
}

func (tp *TaskPool) _tDispatchWork() {
	tk := time.NewTicker(80 * time.Millisecond)
	for {
		select {
		case <-tk.C:
			tp.dispatchWork()
		case <-tp.stopCh:
			log.Warn("[ants pool] task dispatch work quit")
			return
		}
	}
}

func (tp *TaskPool) dispatchWork() {
	if !atomic.CompareAndSwapInt32(&tp.isDisPatch, 1, 2) {
		return
	}
	defer atomic.StoreInt32(&tp.isDisPatch, 1)

	for t := range tp.waitTaskQueueCh {
		// 派发
		err := tp.submit(t)
		if err != nil {
			return
		}
		// 是否满
		if tp.p.Running() >= tp.p.Cap() {
			return
		}
	}
}

func (tp *TaskPool) AddTask(t func()) {
	select {
	case tp.waitTaskQueueCh <- t:
	default:
		log.Error("[ants pool] addWaitQueue full. ch len:%d", len(tp.waitTaskQueueCh))
	}
}

func (tp *TaskPool) submit(t func()) (err error) {
	err = tp.p.Submit(t)
	if errors.Is(err, ants.ErrPoolOverload) { // 满了 放到队列中
		tp.AddTask(t)
		return err
	}
	if err != nil {
		log.Error("[ants pool] Submit err:%+v", err)
		return err
	}
	return err
}

func (tp *TaskPool) ReleaseAntPool() {
	// 先判断是否执行完 再关闭
	tk := time.NewTicker(500 * time.Millisecond)
	for range tk.C {
		if tp.p.Running() > 0 {
			continue
		}
		tp.stopCh <- true
		tp.p.Release()
		return
	}
}
