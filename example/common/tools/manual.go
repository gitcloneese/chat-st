package tools

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"sync"
)

// 任务手册
// manual/getList
// manual/getGrandTotalTaskList
const (
	getList            = "/xy3-%v/manual/getList"
	grandTotalTaskList = "/xy3-%v/manual/getGrandTotalTaskList"
)

var (
	getListPath               string
	getGrandTotalTaskListPath string
)

func initManualPath() {
	getListPath = fmt.Sprintf(getList, ServerId)
	getGrandTotalTaskListPath = fmt.Sprintf(grandTotalTaskList, ServerId)
}

func RunManualListReq() {
	RunWithLogTick("manualGetListReq", FunByKey(1), fmt.Sprintf("%v%v", AccountAddr, getListPath))
}

func RunManualGrandTotalListReq() {
	RunWithLogTick("manualGrandTotalListReq", FunByKey(2), fmt.Sprintf("%v%v", AccountAddr, getGrandTotalTaskListPath))
}

// 压一个玩家的所有接口， 每个接口执行10000次

// FunByKey
// 任务手册
func FunByKey(key int) func() {
	if len(GameLoginResp) < 1 {
		panic("GameLoginResp is empty")
	}
	f := funByKey(key)
	var token string
	if TestOne {
		for _, v := range GameLoginResp {
			token = v.PlayerToken
			break
		}
	}
	return func() {
		wg := &sync.WaitGroup{}
		p, _ := ants.NewPoolWithFunc(C, func(i interface{}) {
			defer wg.Done()
			req := i.(string)
			err := f(req)
			if err != nil {
				Error("%v failed, err:%v", GetFunctionName(f), err)
			}
		})
		defer p.Release()
		if TestOne {
			for i := 0; i < N; i++ {
				wg.Add(1)
				err := p.Invoke(token)
				if err != nil {
					print("%v failed Invoke failed, err:%v\n", GetFunctionName(f), err)
				}
			}
		} else {
			for _, v := range GameLoginResp {
				wg.Add(1)
				err := p.Invoke(v.PlayerToken)
				if err != nil {
					print("%v failed Invoke failed, err:%v\n", GetFunctionName(f), err)
				}
			}
		}
		wg.Wait()
	}
}

func funByKey(key int) func(string) error {
	switch key {
	case 1:
		return reqManualList
	case 2:
		return reqGrandTotalList
	}
	return reqManualList
}

func reqManualList(token string) error {
	path := fmt.Sprintf("%v%v", AccountAddr, getListPath)
	_, err := HttpGet(path, nil, map[string]string{
		"Authorization": bearToken(token),
	})
	return err
}

func reqGrandTotalList(token string) error {
	path := fmt.Sprintf("%v%v", AccountAddr, getGrandTotalTaskListPath)
	_, err := HttpGet(path, nil, map[string]string{
		"Authorization": bearToken(token),
	})
	return err
}
