package util

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"
	"xy3-proto/pkg/log"
)

// recover日志记录，用于本地调试输出日志
func Crash(name string) {
	//log.Printf("Crash begin.")
	t := time.Now()
	strFileName := fmt.Sprintf("crash-%s-%04d-%02d-%02d_%02d_%02d_%02d_%d.log",
		name,
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		os.Getpid())

	f, errFile := os.OpenFile(strFileName, os.O_CREATE|os.O_WRONLY, 0644)
	if errFile != nil {
		//log.Printf("OpenFile begin.")
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	} else {
		//log.Printf("recover begin.")
		if err := recover(); err != nil {
			//	log.Printf("recover Stack.")
			_, err := f.Write(debug.Stack())
			if err != nil {
				log.Error("Crash Write Error: %v", err)
			}
		}
	}

	defer f.Close()
}

// 从异常中恢复协程，用于线上。
func RecoverFromError(cb func()) {
	if e := recover(); e != nil {
		log.Error("Recover => %s:%s\n", e, debug.Stack())
		if nil != cb {
			cb()
		}
	}
}

// SafeDo 捕获err
func SafeDo(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("SafeDo err: :%v, \n stack: %s", err, debug.Stack())
			}
		}()
		f()
	}()
}

// ErrorLoop 循环执行函数
func ErrorLoop(f func() error, n ...int) {
	// 默认执行三次
	defaultTimes := 3
	maxTimes := 10 // 最大执行次数

	if len(n) > 0 && n[0] <= maxTimes {
		defaultTimes = n[0]
	}
	defer func() {
		if err := recover(); err != nil {
			log.Error("Loop err: %s", err)
		}
	}()

	for defaultTimes > 0 {
		defaultTimes--
		if err := f(); err == nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func Loop(f func(), n ...int) {
	// 默认执行三次
	defaultTimes := 3
	maxTimes := 10 // 最大执行次数

	if len(n) > 0 && n[0] <= maxTimes {
		defaultTimes = n[0]
	}
	defer func() {
		if err := recover(); err != nil {
			log.Error("Loop err: %s", err)
		}
	}()

	for defaultTimes > 0 {
		defaultTimes--
		f()
		time.Sleep(10 * time.Millisecond)
	}
}
