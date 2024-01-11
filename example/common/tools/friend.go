package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"sync"
	pbfriend "xy3-proto/friend"
)

const (
	// 好友请求
	friendRequestPath = "/xy3-cross/friend/FriendRequest"
	// 好友请求列表
	friendRequestListPath = "/xy3-cross/friend/FriendRequestList"
	// 好有搜索
	friendSearchPath = "/xy3-cross/friend/FriendSeach"
)
const (
	keyFriendRequest = iota
	keyFriendRequestList
	keyFriendSearch
)

func RunFriendRequestReq() {
	RunWithLogTick("FriendRequestReq", friendRequestByKey(keyFriendRequest), fmt.Sprintf("%v%v", AccountAddr, friendRequestPath))
}

func RunFriendRequestListReq() {
	RunWithLogTick("FriendRequestListReq", friendRequestByKey(keyFriendRequestList), fmt.Sprintf("%v%v", AccountAddr, friendRequestListPath))
}

func RunFriendSearchReq() {
	RunWithLogTick("FriendSearchReq", friendRequestByKey(keyFriendSearch), fmt.Sprintf("%v%v", AccountAddr, friendSearchPath))
}

// 好友请求
func friendRequestByKey(key int) func() {
	if len(GameLoginResp) < 1 {
		panic("GameLoginResp is empty")
	}
	var f func(string) error
	switch key {
	case keyFriendRequest:
		f = friendRequestFun
	case keyFriendRequestList:
		f = friendRequestListFun
	case keyFriendSearch:
		f = friendRequestSearchFun
	}
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
			wg.Wait()
		}
	}
}

// 好友请求
func friendRequestFun(token string) error {
	reqB, err := json.Marshal(pbfriend.FriendRequestReq{
		FriendId: []int64{11451102, 11781103},
	})
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%v%v", AccountAddr, friendRequestPath)
	headers := map[string]string{
		"Authorization": bearToken(token),
	}
	_, err = HttpPost(path, bytes.NewReader(reqB), headers)
	if err != nil {
		return err
	}
	return nil
}

// 好友请求
func friendRequestListFun(token string) error {
	reqB, err := json.Marshal(pbfriend.FriendRequestListReq{})
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%v%v", AccountAddr, friendRequestListPath)
	headers := map[string]string{
		"Authorization": bearToken(token),
	}
	_, err = HttpPost(path, bytes.NewReader(reqB), headers)
	if err != nil {
		return err
	}
	return nil
}

// 好友搜索推荐
func friendRequestSearchFun(token string) error {
	var (
		reqB []byte
		err  error
	)
	if Data != "" {
		reqB = []byte(Data)
	} else {
		reqB, err = json.Marshal(pbfriend.FriendSeachReq{
			SeachParam: "南海",
		})
		if err != nil {
			return err
		}
	}
	path := fmt.Sprintf("%v%v", AccountAddr, friendSearchPath)
	headers := map[string]string{
		"Authorization": bearToken(token),
	}
	_, err = HttpPost(path, bytes.NewReader(reqB), headers)
	if err != nil {
		return err
	}
	return nil
}
