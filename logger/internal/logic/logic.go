package logic

import (
	"context"
	"reflect"
	"x-server/logger/internal/dao"
	pb "xy3-proto/logger"
	"xy3-proto/pkg/log"

	"github.com/google/wire"
)

var (
	ProviderSet = wire.NewSet(New)
)

// Logic
// 逻辑处理
type Logic struct {
	dao      dao.Dao
	handlers map[pb.ELogCategory]string
}

func New(d dao.Dao) (*Logic, error) {
	logic := &Logic{
		dao:      d,
		handlers: map[pb.ELogCategory]string{},
	}
	logic.initRoute()
	logic.createTable()
	return logic, nil
}

// 添加路由
func (p *Logic) addRoute(logCategory pb.ELogCategory, name string) {
	if _, has := p.handlers[logCategory]; has {
		log.Error("addRoute %v already exists!", logCategory)
		return
	}
	p.handlers[logCategory] = name
}

func (p *Logic) Handler(msg *pb.LogMsgs) error {
	if msg == nil {
		return nil
	}
	for _, v := range msg.Messages {
		if name, has := p.handlers[v.Category]; has {
			object := reflect.ValueOf(p)
			fc := object.MethodByName(name)
			if fc.IsNil() {
				log.Error("Category %v Handle %v Not Find!", v.Category, name)
			} else {
				args := []reflect.Value{
					reflect.ValueOf(context.Background()),
					reflect.ValueOf(v),
				}
				err := fc.Call(args)
				if err != nil || len(err) > 0 {
					log.Error("Category %v Handle %v Error err:%v", v.Category, name, err)
				}
			}
		} else {
			log.Error("LogCategory %v Not find handler!", v.Category)
		}
	}
	return nil
}
