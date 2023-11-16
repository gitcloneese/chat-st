//go:build !windows
// +build !windows

// xcfg 扩展paladin加载配置，加载配置出错时打印日志，返回上层错误wrap了配置文件名
package xcfg

import (
	"xy3-proto/pkg/conf/paladin"
	"xy3-proto/pkg/log"
)

// type wrapError struct {
// 	file string
// 	err  error
// }

// func (e *wrapError) Error() string {
// 	if e.err == nil {
// 		return "no err"
// 	}
// 	return e.err.Error()
// }

// func (e *wrapError) Unwrap() error {
// 	return e.err
// }

// func (e *wrapError) File() string {
// 	return e.file
// }

// func wrapCfgError(file string, err error) error {
// 	return &wrapError{file, errors.Wrapf(err, "cfg file %s ", file)}
// }

// func isCfgError(err error) bool {
// 	_, ok := err.(*wrapError)
// 	return ok
// }

// func errorCfgName(err error) string {
// 	e, ok := err.(*wrapError)
// 	if !ok {
// 		return ""
// 	}
// 	return e.File()
// }

func Panic(err error) {
	if err == nil {
		return
	}
	log.Error("配置加载错误 %s", err.Error())

	panic(err)
}

// type wrapCfg struct {
// 	file   string
// 	setter paladin.Setter
// }

// func (w *wrapCfg) Set(text []byte) (err error) {
// 	if w.setter == nil {
// 		return nil
// 	}
// 	err = w.setter.Set(text)
// 	if err != nil {
// 		err = wrapCfgError(w.File(), err)
// 		return
// 	}
// 	return nil
// }

// func (w *wrapCfg) File() string {
// 	return w.file
// }

func Watch(keys []string, mm map[string]paladin.Setter) {
	paladin.Watch(keys, mm)
}
