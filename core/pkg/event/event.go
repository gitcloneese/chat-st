// package event 事件系统
package event

import (
	"errors"
	"reflect"
	"runtime/debug"

	"xy3-proto/pkg/log"
)

type EventDispatcher interface {
	AddEventListener(eventType string, f interface{}) error
	RemoveEventListener(eventType string)
	DispatchEvent(eventType string, params ...interface{}) bool
}

type eventDispatcher struct {
	listeners map[string][]reflect.Value
}

func NewEventDispatcher() EventDispatcher {
	return &eventDispatcher{listeners: make(map[string][]reflect.Value)}
}

// 添加事件监听
func (p *eventDispatcher) AddEventListener(eventType string, f interface{}) error {
	tp := reflect.TypeOf(f)
	if tp.Kind() != reflect.Func {
		return errors.New("listener is not a func")
	}
	reflectVal := reflect.ValueOf(f)

	// // 校验方法是否是导出函数
	// if !isExportedOrBuiltinType(tp) {
	// 	return errors.New("listener is not a export func")
	// }

	if _, ok := p.listeners[eventType]; !ok {
		p.listeners[eventType] = make([]reflect.Value, 0)
	}

	// 去重处理
	for _, v := range p.listeners[eventType] {
		if v.Pointer() == reflectVal.Pointer() {
			return errors.New("listener already exist")
		}
	}

	// 注册
	p.listeners[eventType] = append(p.listeners[eventType], reflectVal)
	return nil
}

// 移除事件监听
func (p *eventDispatcher) RemoveEventListener(eventType string) {
	delete(p.listeners, eventType)
}

// 派发事件
func (p *eventDispatcher) DispatchEvent(eventType string, params ...interface{}) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Error("event DispatchEvent panic eventType:%s dump err:%v stack:%s", eventType, err, string(debug.Stack()))
			return
		}
	}()

	list, ok := p.listeners[eventType]
	if !ok {
		return false
	}

	if len(list) <= 0 {
		return false
	}

	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	for _, v := range list {
		v.Call(in)
	}

	return true
}

// func isExported(name string) bool {
// 	r, _ := utf8.DecodeRuneInString(name)
// 	return unicode.IsUpper(r)
// }

// func isExportedOrBuiltinType(t reflect.Type) bool {
// 	for t.Kind() == reflect.Ptr {
// 		t = t.Elem()
// 	}
// 	return isExported(t.Name()) || t.PkgPath() == ""
// }
