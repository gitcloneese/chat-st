package friend

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	util2 "x-server/core/pkg/util"

	"x-server/core/dao/model"
	"xy3-proto/pkg/log"

	v8 "github.com/go-redis/redis/v8"
)

// CacheLeaseHero 缓存角色租借英雄
func (f *Friend) CacheLeaseHero(id int64, leaseHero *model.LeaseHero) (err error) {
	key := fmt.Sprintf(model.RedisKey_LeaseHero_Hero, id)
	_, err = f.client.SAdd(context.TODO(), key, leaseHero.Hero.HeroID).Result()
	if err != nil {
		return err
	}
	key = fmt.Sprintf(model.RedisKey_LeaseHero_Hero_Hash, id, leaseHero.Hero.HeroID)
	data := make(map[string]interface{})
	heroBytes, err := json.Marshal(leaseHero.Hero)
	if err != nil {
		return err
	}
	data["OwnerID"] = id
	data["Hero"] = heroBytes
	exBytes, err := json.Marshal(leaseHero.ExclusiveWeapon)
	if err != nil {
		return err
	}
	data["ExclusiveWeapon"] = exBytes
	_, err = f.client.HSet(context.TODO(), key, data).Result()
	if err != nil {
		return err
	}
	return err
}

// 更新战斗次数
func (f *Friend) UpdateLeaseFightCount(id int64, fightCount map[int32]int32) (err error) {
	data := make(map[string]interface{})
	for k, v := range fightCount {
		data[util2.Int32ToStr(k)] = v
	}
	_, err = f.client.HSet(context.TODO(), fmt.Sprintf(model.RedisKey_LeaseHero_FightCount, id), data).Result()
	return err
}

func (f *Friend) CacheLeaseRequest(uid int64, request *model.LeaseRequest) (err error) {
	key := fmt.Sprintf(model.RedisKey_LeaseHero_Request, uid)
	f.client.SAdd(context.TODO(), key, request.ID)

	data := make(map[string]interface{})
	data["ID"] = request.ID
	data["RoleID"] = request.RoleID
	data["HeroID"] = request.HeroID
	data["Time"] = request.Time
	data["State"] = request.State
	key = fmt.Sprintf(model.RedisKey_LeaseHero_Request_Hash, request.ID)
	f.client.HSet(context.TODO(), key, data)
	return err
}

func (f *Friend) UpdateLeaseRequest(request *model.LeaseRequest) (err error) {
	data := make(map[string]interface{})
	data["ID"] = request.ID
	data["RoleID"] = request.RoleID
	data["HeroID"] = request.HeroID
	data["Time"] = request.Time
	data["State"] = request.State
	key := fmt.Sprintf(model.RedisKey_LeaseHero_Request_Hash, request.ID)
	f.client.HSet(context.TODO(), key, data)
	return err
}

func (f *Friend) GetLeaseRequest(reqIds []int64) (requestList []*model.LeaseRequest, err error) {
	pipe := f.client.Pipeline()
	for _, reqId := range reqIds {
		key := fmt.Sprintf(model.RedisKey_LeaseHero_Request_Hash, reqId)
		pipe.HGetAll(context.TODO(), key)
	}
	result, err := pipe.Exec(context.TODO())
	if err != nil {
		return requestList, err
	}
	requestList = make([]*model.LeaseRequest, 0)
	for _, res := range result {
		m, err := res.(*v8.StringStringMapCmd).Result()
		if err != nil {
			continue
		}
		request := &model.LeaseRequest{
			ID:     util2.StrToInt64(m["ID"]),
			RoleID: util2.StrToInt64(m["RoleID"]),
			HeroID: util2.StrToInt32(m["HeroID"]),
			Time:   util2.StrToInt64(m["Time"]),
			State:  util2.StrToInt32(m["State"]),
		}
		requestList = append(requestList, request)
	}
	return requestList, err
}

func (f *Friend) GetAllLeaseRequest(uid int64) (requestMap map[int64]*model.LeaseRequest, err error) {
	key := fmt.Sprintf(model.RedisKey_LeaseHero_Request, uid)
	result, err := f.client.SMembers(context.TODO(), key).Result()
	if err != nil {
		return requestMap, err
	}
	pipe := f.client.Pipeline()
	for _, id := range result {
		key := fmt.Sprintf(model.RedisKey_LeaseHero_Request_Hash, util2.StrToInt64(id))
		pipe.HGetAll(context.TODO(), key)
	}
	pipeResult, err := pipe.Exec(context.TODO())
	if err != nil {
		return requestMap, err
	}
	requestMap = make(map[int64]*model.LeaseRequest)
	for _, res := range pipeResult {
		m, err := res.(*v8.StringStringMapCmd).Result()
		if err != nil {
			continue
		}
		request := &model.LeaseRequest{
			ID:     util2.StrToInt64(m["ID"]),
			RoleID: util2.StrToInt64(m["RoleID"]),
			HeroID: util2.StrToInt32(m["HeroID"]),
			Time:   util2.StrToInt64(m["Time"]),
			State:  util2.StrToInt32(m["State"]),
		}
		requestMap[request.ID] = request
	}
	return requestMap, err
}

func (f *Friend) GetLeaseHero(uid int64) (heros []*model.LeaseHero, err error) {
	key := fmt.Sprintf(model.RedisKey_LeaseHero_Hero, uid)
	result, err := f.client.SMembers(context.TODO(), key).Result()
	if err != nil {
		return heros, err
	}

	pipe := f.client.Pipeline()
	for _, res := range result {
		key := fmt.Sprintf(model.RedisKey_LeaseHero_Hero_Hash, uid, res)
		pipe.HGetAll(context.TODO(), key)
	}
	pipeResult, err := pipe.Exec(context.TODO())
	if err != nil {
		return heros, err
	}
	heros = make([]*model.LeaseHero, 0)
	for _, res := range pipeResult {
		m, err := res.(*v8.StringStringMapCmd).Result()
		if err != nil {
			continue
		}
		leaseHero := &model.LeaseHero{
			Hero:            &model.LeaseHeroObj{},
			ExclusiveWeapon: &model.LeaseExclusiveWeaponObj{},
			OwnerID:         util2.StrToInt64(m["OwnerID"]),
			LeaseID:         util2.StrToInt64(m["LeaseID"]),
			LeaseTime:       util2.StrToInt64(m["LeaseTime"]),
			CD:              util2.StrToInt64(m["CD"]),
		}

		if _, ok := m["Hero"]; ok {
			err = json.Unmarshal([]byte(m["Hero"]), leaseHero.Hero)
			if err != nil {
				continue
			}
		}

		if _, ok := m["ExclusiveWeapon"]; ok {
			err = json.Unmarshal([]byte(m["ExclusiveWeapon"]), leaseHero.ExclusiveWeapon)
			if err != nil {
				continue
			}
		}
		heros = append(heros, leaseHero)
	}
	return heros, err
}

func (f *Friend) UpdateLeaseHero(id int64, hero *model.LeaseHero) (err error) {
	data := make(map[string]interface{})
	heroBytes, err := json.Marshal(hero.Hero)
	if err != nil {
		return err
	}
	data["Hero"] = heroBytes
	exBytes, err := json.Marshal(hero.ExclusiveWeapon)
	if err != nil {
		return err
	}
	data["ExclusiveWeapon"] = exBytes
	data["OwnerID"] = hero.OwnerID
	data["LeaseID"] = hero.LeaseID
	data["LeaseTime"] = hero.LeaseTime
	data["CD"] = hero.CD

	key := fmt.Sprintf(model.RedisKey_LeaseHero_Hero_Hash, id, hero.Hero.HeroID)
	_, err = f.client.HSet(context.TODO(), key, data).Result()
	if err != nil {
		return err
	}

	return err
}

func (f *Friend) GetLeaseHeroList(ids []int64) (heros []*model.LeaseHero, err error) {
	pipe := f.client.Pipeline()
	for _, id := range ids {
		key := fmt.Sprintf(model.RedisKey_LeaseHero_Hero, id)
		pipe.SMembers(context.TODO(), key)
	}
	heroIdMap := make(map[int64][]int32)
	result, err := pipe.Exec(context.TODO())
	if err != nil {
		log.Error("GetLeaseHeroList Exec Error: %v", err)
	}
	for index, res := range result {
		m, err := res.(*v8.StringSliceCmd).Result()
		if err != nil {
			continue
		}
		for _, heroID := range m {
			heroIdMap[ids[index]] = append(heroIdMap[ids[index]], util2.StrToInt32(heroID))
		}
	}
	pipe = f.client.Pipeline()
	for key, heroids := range heroIdMap {
		for _, heroid := range heroids {
			key := fmt.Sprintf(model.RedisKey_LeaseHero_Hero_Hash, key, heroid)
			pipe.HGetAll(context.TODO(), key)
		}
	}
	heros = make([]*model.LeaseHero, 0)
	result, err = pipe.Exec(context.TODO())
	for _, res := range result {
		m, err := res.(*v8.StringStringMapCmd).Result()
		if err != nil {
			continue
		}
		leaseHero := &model.LeaseHero{
			Hero:            &model.LeaseHeroObj{},
			ExclusiveWeapon: &model.LeaseExclusiveWeaponObj{},
			OwnerID:         util2.StrToInt64(m["OwnerID"]),
			LeaseID:         util2.StrToInt64(m["LeaseID"]),
			LeaseTime:       util2.StrToInt64(m["LeaseTime"]),
			CD:              util2.StrToInt64(m["CD"]),
		}

		if _, ok := m["Hero"]; ok {
			err = json.Unmarshal([]byte(m["Hero"]), leaseHero.Hero)
			if err != nil {
				continue
			}
		}

		if _, ok := m["ExclusiveWeapon"]; ok {
			err = json.Unmarshal([]byte(m["ExclusiveWeapon"]), leaseHero.ExclusiveWeapon)
			if err != nil {
				continue
			}
		}

		heros = append(heros, leaseHero)
	}
	return heros, err
}

func (f *Friend) GetLeaseHeroListByHeroID(ids []int64, heroID int32) (heros []*model.LeaseHero, err error) {
	pipe := f.client.Pipeline()
	for _, id := range ids {
		key := fmt.Sprintf(model.RedisKey_LeaseHero_Hero, id)
		pipe.SIsMember(context.TODO(), key, heroID)
	}
	result, err := pipe.Exec(context.TODO())
	if err != nil {
		return heros, err
	}
	friendIds := []int64{}
	for index, res := range result {
		m, err := res.(*v8.BoolCmd).Result()
		if err != nil {
			continue
		}
		if m {
			friendIds = append(friendIds, ids[index])
		}
	}
	pipe = f.client.Pipeline()
	for _, friendId := range friendIds {
		key := fmt.Sprintf(model.RedisKey_LeaseHero_Hero_Hash, friendId, heroID)
		pipe.HGetAll(context.TODO(), key)
	}
	heros = make([]*model.LeaseHero, 0)
	result, err = pipe.Exec(context.TODO())
	if err != nil {
		return heros, err
	}
	for _, res := range result {
		m, err := res.(*v8.StringStringMapCmd).Result()
		if err != nil {
			continue
		}
		leaseHero := &model.LeaseHero{
			Hero:            &model.LeaseHeroObj{},
			ExclusiveWeapon: &model.LeaseExclusiveWeaponObj{},
			OwnerID:         util2.StrToInt64(m["OwnerID"]),
			LeaseID:         util2.StrToInt64(m["LeaseID"]),
			LeaseTime:       util2.StrToInt64(m["LeaseTime"]),
			CD:              util2.StrToInt64(m["CD"]),
		}

		if _, ok := m["Hero"]; ok {
			err = json.Unmarshal([]byte(m["Hero"]), leaseHero.Hero)
			if err != nil {
				continue
			}
		}

		if _, ok := m["ExclusiveWeapon"]; ok {
			err = json.Unmarshal([]byte(m["ExclusiveWeapon"]), leaseHero.ExclusiveWeapon)
			if err != nil {
				continue
			}
		}
		heros = append(heros, leaseHero)
	}
	return heros, err
}

func (f *Friend) GetSelfLeaseHeroList(id int64) (heros []*model.LeaseHero, err error) {
	key := fmt.Sprintf(model.RedisKey_Self_LeaseHero_Hero, id)
	result, err := f.client.SMembers(context.TODO(), key).Result()
	if err != nil {
		return heros, err
	}
	pipe := f.client.Pipeline()
	for _, res := range result {
		key = fmt.Sprintf(model.RedisKey_Self_LeaseHero_Hero_Hash, id, res)
		pipe.HGetAll(context.TODO(), key)
	}
	heros = make([]*model.LeaseHero, 0)
	pipeResult, err := pipe.Exec(context.TODO())
	if err != nil {
		return heros, err
	}
	for _, res := range pipeResult {
		m, err := res.(*v8.StringStringMapCmd).Result()
		if err != nil {
			continue
		}
		leaseHero := &model.LeaseHero{
			Hero:            &model.LeaseHeroObj{},
			ExclusiveWeapon: &model.LeaseExclusiveWeaponObj{},
			OwnerID:         util2.StrToInt64(m["OwnerID"]),
			LeaseID:         util2.StrToInt64(m["LeaseID"]),
			LeaseTime:       util2.StrToInt64(m["LeaseTime"]),
			CD:              util2.StrToInt64(m["CD"]),
		}

		if _, ok := m["Hero"]; ok {
			err = json.Unmarshal([]byte(m["Hero"]), leaseHero.Hero)
			if err != nil {
				continue
			}
		}

		if _, ok := m["ExclusiveWeapon"]; ok {
			err = json.Unmarshal([]byte(m["ExclusiveWeapon"]), leaseHero.ExclusiveWeapon)
			if err != nil {
				continue
			}
		}
		heros = append(heros, leaseHero)
	}
	return heros, err
}

func (f *Friend) CacheSelfLeaseHero(id int64, hero *model.LeaseHero) (err error) {
	key := fmt.Sprintf(model.RedisKey_Self_LeaseHero_Hero, id)
	_, err = f.client.SAdd(context.TODO(), key, hero.Hero.HeroID).Result()
	if err != nil {
		return
	}
	data := make(map[string]interface{})
	heroBytes, err := json.Marshal(hero.Hero)
	if err != nil {
		return
	}
	data["Hero"] = heroBytes
	exBytes, err := json.Marshal(hero.ExclusiveWeapon)
	if err != nil {
		return
	}
	data["ExclusiveWeapon"] = exBytes
	data["OwnerID"] = hero.OwnerID
	data["LeaseID"] = hero.LeaseID
	data["LeaseTime"] = hero.LeaseTime
	data["CD"] = hero.CD

	key = fmt.Sprintf(model.RedisKey_Self_LeaseHero_Hero_Hash, id, hero.Hero.HeroID)
	_, err = f.client.HSet(context.TODO(), key, data).Result()
	if err != nil {
		return
	}
	return
}

func (f *Friend) RemoveSelfHero(id int64) (err error) {
	key := fmt.Sprintf(model.RedisKey_Self_LeaseHero_Hero, id)
	result, err := f.client.SMembers(context.TODO(), key).Result()
	if err != nil {
		return
	}
	pipe := f.client.Pipeline()
	for _, res := range result {
		key = fmt.Sprintf(model.RedisKey_Self_LeaseHero_Hero_Hash, id, util2.StrToInt64(res))
		pipe.Del(context.TODO(), key)
	}
	key = fmt.Sprintf(model.RedisKey_Self_LeaseHero_Hero, id)
	pipe.Del(context.TODO(), key)
	_, err = pipe.Exec(context.TODO())
	return
}

func (f *Friend) DelSelfLeasehero(id int64, heroId int32) (err error) {
	key := fmt.Sprintf(model.RedisKey_Self_LeaseHero_Hero, id)
	_, err = f.client.SRem(context.TODO(), key, heroId).Result()
	if err != nil {
		return
	}

	key = fmt.Sprintf(model.RedisKey_Self_LeaseHero_Hero_Hash, id, heroId)
	_, err = f.client.Del(context.TODO(), key).Result()
	return
}

func (f *Friend) DelLeaseRequestList(id int64, reqIds []int64) (err error) {
	pipe := f.client.Pipeline()
	for _, reqId := range reqIds {
		key := fmt.Sprintf(model.RedisKey_LeaseHero_Request, id)
		_, err = pipe.SRem(context.TODO(), key, reqId).Result()
		if err != nil {
			continue
		}
		_, err = pipe.Del(context.TODO(), fmt.Sprintf(model.RedisKey_LeaseHero_Request_Hash, reqId)).Result()
		if err != nil {
			continue
		}
	}
	_, err = pipe.Exec(context.TODO())
	return
}

func (f *Friend) GetLeaseHeroTask(id int64, taskIds []int32) (leaseTask []*model.LeaseTask, err error) {
	pipe := f.client.Pipeline()
	for _, taskId := range taskIds {
		key := fmt.Sprintf(model.RedisKey_LeaseHero_Task, id, taskId)
		pipe.HGetAll(context.TODO(), key)
	}
	result, err := pipe.Exec(context.TODO())
	if err != nil {
		return
	}
	leaseTask = make([]*model.LeaseTask, 0)
	for _, res := range result {
		m, err := res.(*v8.StringStringMapCmd).Result()
		if err != nil || len(m) == 0 {
			continue
		}
		task := &model.LeaseTask{
			TaskID: util2.StrToInt32(m["TaskID"]),
			Num:    util2.StrToInt32(m["Num"]),
			State:  util2.StrToInt32(m["State"]),
		}
		leaseTask = append(leaseTask, task)
	}
	return
}

func (f *Friend) updateLeaseHeroTask(uid int64, taskID int32, condition []int32) (err error) {
	tasks, err := f.GetLeaseHeroTask(uid, []int32{taskID})
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		key := fmt.Sprintf(model.RedisKey_LeaseHero_Task, uid, taskID)
		data := make(map[string]interface{})
		data["TaskID"] = taskID
		data["Num"] = 1
		data["State"] = 0
		f.client.HSet(context.TODO(), key, data)
	} else {
		for _, task := range tasks {
			if task.State != 0 {
				continue
			}
			task.Num++
			if task.Num >= condition[1] {
				task.State = 1
			}
			key := fmt.Sprintf(model.RedisKey_LeaseHero_Task, uid, task.TaskID)
			data := make(map[string]interface{})
			data["TaskID"] = task.TaskID
			data["Num"] = task.Num
			data["State"] = task.State
			f.client.HSet(context.TODO(), key, data)
		}
	}
	return err
}

func (f *Friend) UpdateTask(leaseType int32, uid int64, leaseHero *model.LeaseHero, task map[int32][]int32) {
	for taskID, condition := range task {
		if leaseType == 1 { // 借出
			if condition[0] == model.TaskType1 {
				err := f.updateLeaseHeroTask(uid, taskID, condition)
				if err != nil {
					log.Error("UpdateTask updateLeaseHeroTask Error: %v", err)
				}
			}
			if condition[0] == model.TaskType2 && len(condition) >= 2 {
				if leaseHero.Hero.Awaken >= condition[1] {
					err := f.updateLeaseHeroTask(uid, taskID, condition)
					if err != nil {
						log.Error("UpdateTask updateLeaseHeroTask Error: %v", err)
					}
				}
			}
		} else {
			if condition[0] == model.TaskType3 {
				err := f.updateLeaseHeroTask(uid, taskID, condition)
				if err != nil {
					log.Error("UpdateTask updateLeaseHeroTask Error: %v", err)
				}
			}
			if condition[0] == model.TaskType4 && len(condition) >= 2 {
				if leaseHero.Hero.Awaken >= condition[1] {
					err := f.updateLeaseHeroTask(uid, taskID, condition)
					if err != nil {
						log.Error("UpdateTask updateLeaseHeroTask Error: %v", err)
					}
				}
			}
		}
	}
}

// 保存任务
func (f *Friend) SaveTask(uid int64, task *model.LeaseTask) (err error) {
	key := fmt.Sprintf(model.RedisKey_LeaseHero_Task, uid, task.TaskID)
	data := make(map[string]interface{})
	data["TaskID"] = task.TaskID
	data["Num"] = task.Num
	data["State"] = task.State
	f.client.HSet(context.TODO(), key, data)
	return
}

// 获取战斗次数
func (f *Friend) GetFightCount(id int64) (fightCount map[int32]int32, err error) {
	result, err := f.client.HGetAll(context.TODO(), fmt.Sprintf(model.RedisKey_LeaseHero_FightCount, id)).Result()
	if err != nil {
		return
	}
	fightCount = make(map[int32]int32)
	for k, v := range result {
		fightCount[util2.StrToInt32(k)] = util2.StrToInt32(v)
	}
	return
}

// 更新战斗次数
func (f *Friend) UpdateFightCount(id int64, fightCount map[int32]int32) (err error) {
	f.client.HSet(context.TODO(), fmt.Sprintf(model.RedisKey_LeaseHero_FightCount, id), fightCount)
	return
}

// 重置数据
func (f *Friend) LeaseHeroReset(now int64) {
	weekday := util2.GetWeekDay(time.Now())
	if weekday != 1 {
		return
	}
}

// 重置租借数据
//
//nolint:funlen
func (f *Friend) ResetLeaseHero(uid int64) (err error) {
	// 删除自己对别人的请求
	selfLeaseHeros, err := f.GetSelfLeaseHeroList(uid)
	if err != nil {
		return err
	}
	for _, hero := range selfLeaseHeros {
		selfRequestMap, err := f.GetAllLeaseRequest(hero.OwnerID)
		if err != nil {
			log.Error("ResetLeaseHero GetAllLeaseRequest err:[%v]", err)
			continue
		}
		delReq := []int64{}
		for _, request := range selfRequestMap {
			if request.RoleID == uid {
				delReq = append(delReq, request.ID)
			}
		}
		err = f.DelLeaseRequestList(hero.OwnerID, delReq)
		if err != nil {
			log.Error("ResetLeaseHero DelLeaseRequestList Error: %v", err)
		}
	}
	// 删除别人对自己的请求
	selfRequestMap, err := f.GetAllLeaseRequest(uid)
	if err != nil {
		return err
	}
	delReq := []int64{}
	for _, req := range selfRequestMap {
		delReq = append(delReq, req.ID)
	}
	err = f.DelLeaseRequestList(uid, delReq)
	if err != nil {
		return err
	}

	selfHeroList, err := f.GetSelfLeaseHeroList(uid)
	if err != nil {
		return err
	}
	for _, selfHero := range selfHeroList {
		// 删除对方列表请求
		requestMap, err := f.GetAllLeaseRequest(selfHero.OwnerID)
		if err != nil {
			log.Error("ResetLeaseHero GetAllLeaseRequest Error: %v", err)
		}
		for _, req := range requestMap {
			if req.RoleID == uid && req.HeroID == selfHero.Hero.HeroID {
				err = f.DelLeaseRequestList(req.RoleID, []int64{req.ID})
				if err != nil {
					log.Error("ResetLeaseHero DelLeaseRequestList Error: %v", err)
				}
			}
		}
		// 通过租借的英雄重置对方的英雄
		targetHeros, err := f.GetLeaseHero(selfHero.OwnerID)
		if err != nil {
			continue
		}
		for _, targetHero := range targetHeros {
			if selfHero.Hero.HeroID == targetHero.Hero.HeroID {
				targetHero.LeaseID = 0
				targetHero.LeaseTime = 0
				targetHero.CD = 0
				err = f.UpdateLeaseHero(targetHero.OwnerID, targetHero)
				if err != nil {
					log.Error("ResetLeaseHero UpdateLeaseHero Error: %v", err)
				}
			}
		}
	}
	// 删除自己租借的数据
	err = f.RemoveSelfHero(uid)
	if err != nil {
		return err
	}
	// 删除自己租借的历史记录
	_, err = f.client.HDel(context.TODO(), fmt.Sprintf(model.RedisKey_LeaseHero_HistoryLease, uid)).Result()
	return err
}

func (f *Friend) UpdateHistoryLease(uid int64, count int32) (err error) {
	key := fmt.Sprintf(model.RedisKey_LeaseHero_HistoryLease, uid)
	_, err = f.client.HSet(context.TODO(), key, "count", count).Result()
	return
}

func (f *Friend) GetHistoryLease(uid int64) (count int32, err error) {
	key := fmt.Sprintf(model.RedisKey_LeaseHero_HistoryLease, uid)
	result, err := f.client.HGet(context.TODO(), key, "count").Result()
	if err != nil {
		if errors.Is(err, v8.Nil) {
			return 0, nil
		}
		return
	}
	count = util2.StrToInt32(result)
	return
}
