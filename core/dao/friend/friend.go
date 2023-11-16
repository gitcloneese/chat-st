package friend

import (
	"context"
	"fmt"
	"strconv"
	"x-server/core/pkg/util"
	"xy3-proto/pkg/log"

	"x-server/core/dao/model"

	v8 "github.com/go-redis/redis/v8"
)

type Friend struct {
	client *v8.Client
}

func New(r *v8.Client) *Friend {
	return &Friend{
		client: r,
	}
}

// 好友跨天重置
func (f *Friend) FriendReset(now int64) {
	ids, _ := f.client.SMembers(context.TODO(), model.RedisKey_Friend_Point).Result()
	pipe := f.client.Pipeline()
	for _, strId := range ids {
		id := util.StrToInt64(strId)
		pipe.Del(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_Point_Hash, id))
		pipe.Del(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_Del, id))
		pipe.Del(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_Recommend, id))
	}
	pipe.Del(context.TODO(), model.RedisKey_Friend_Point)
	_, err := pipe.Exec(context.TODO())
	if err != nil {
		log.Error("FriendReset Exec Error: %v", err)
	}
}

// 好友点信息存储
func (f *Friend) CacheFriendPointInfo(roleID int64, friendInfo *model.FriendPoint) (err error) {
	f.client.SAdd(context.TODO(), model.RedisKey_Friend_Point, roleID)

	data := make(map[string]interface{})
	data["RoleID"] = friendInfo.RoleID
	data["GetPointCount"] = friendInfo.GetPointCount
	data["GivePointCount"] = friendInfo.GivePointCount
	_, err = f.client.HSet(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_Point_Hash, roleID), data).Result()
	return err
}

// 好友点信息读取
func (f *Friend) GetFriendPointInfo(roleID int64) (friendInfo *model.FriendPoint, err error) {
	result, err := f.client.HGetAll(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_Point_Hash, roleID)).Result()
	if err != nil || len(result) == 0 {
		return friendInfo, err
	}
	friendInfo = &model.FriendPoint{
		RoleID:         util.StrToInt64(result["RoleID"]),
		GetPointCount:  util.StrToInt32(result["GetPointCount"]),
		GivePointCount: util.StrToInt32(result["GivePointCount"]),
	}
	return friendInfo, err
}

// 存储好友申请id列表
func (f *Friend) CacheRequestIDs(roleID int64, requestIDs []int64) (err error) {
	pipe := f.client.Pipeline()
	for _, id := range requestIDs {
		_, err = pipe.SAdd(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_Request, roleID), id).Result()
		if err != nil {
			log.Error("CacheRequestIDs Result Error: %v", err)
		}
	}
	_, err = pipe.Exec(context.TODO())
	return err
}

// 是否是好友申请
func (f *Friend) IsRequest(roleID, requestID int64) (isRequest bool, err error) {
	isRequest, err = f.client.SIsMember(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_Request, roleID), requestID).Result()
	if err != nil {
		isRequest = false
	}
	return isRequest, err
}

// 获取所有好友申请id
func (f *Friend) GetAllRequestID(roleID int64) (ids []int64, err error) {
	result, err := f.client.SMembersMap(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_Request, roleID)).Result()
	if err != nil {
		return
	}
	ids = make([]int64, 0)
	for friendID := range result {
		ids = append(ids, util.StrToInt64(friendID))
	}
	return ids, err
}

// 删除好友申请
func (f *Friend) DelRequest(roleID, id int64) (err error) {
	_, err = f.client.SRem(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_Request, roleID), id).Result()
	return err
}

// 删除好友申请列表
func (f *Friend) DelRequests(roleID int64, ids []int64) (err error) {
	pipe := f.client.Pipeline()
	for _, id := range ids {
		pipe.SRem(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_Request, roleID), id)
	}
	_, err = pipe.Exec(context.TODO())
	return err
}

// 存储好友申请详情
func (f *Friend) CacheRequest(roleID int64, request *model.FriendRequest) (err error) {
	key := fmt.Sprintf(model.RedisKey_Friend_Request_Hash, roleID, request.RoleID)
	data := make(map[string]interface{})
	data["RoleID"] = request.RoleID
	data["Time"] = request.Time
	f.client.HSet(context.TODO(), key, data)
	return err
}

// 获取所有好友申请详情
func (f *Friend) GetAllRequest(roleId int64, requestIds []int64) (request []*model.FriendRequest, err error) {
	pipe := f.client.Pipeline()
	for _, requestId := range requestIds {
		key := fmt.Sprintf(model.RedisKey_Friend_Request_Hash, roleId, requestId)
		pipe.HGetAll(context.TODO(), key)
	}
	result, err := pipe.Exec(context.TODO())
	if err != nil {
		return request, err
	}
	request = make([]*model.FriendRequest, 0)
	for _, res := range result {
		m, err := res.(*v8.StringStringMapCmd).Result()
		if err != nil {
			continue
		}
		friendRequest := &model.FriendRequest{
			RoleID: util.StrToInt64(m["RoleID"]),
			Time:   util.StrToInt64(m["Time"]),
		}

		request = append(request, friendRequest)
	}
	return request, err
}

// 存储好友id
func (f *Friend) CacheFriendID(roleID, friendID int64) (err error) {
	_, err = f.client.SAdd(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_Friend, roleID), friendID).Result()
	return err
}

// 是否是好友
func (f *Friend) IsFriend(roleID, friendID int64) (isFriend bool, err error) {
	isFriend, err = f.client.SIsMember(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_Friend, roleID), friendID).Result()
	if err != nil {
		isFriend = false
	}
	return isFriend, err
}

// 获取所有好友id
func (f *Friend) GetAllFriendID(roleID int64) (ids []int64, err error) {
	result, err := f.client.SMembersMap(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_Friend, roleID)).Result()
	if err != nil {
		return ids, err
	}
	ids = make([]int64, 0)
	for friendID := range result {
		ids = append(ids, util.StrToInt64(friendID))
	}
	return ids, err
}

// 存储好友信息
func (f *Friend) CacheFriend(roleID int64, info *model.FriendInfo) (err error) {
	dataMap := make(map[string]interface{})
	dataMap["RoleID"] = info.RoleID
	dataMap["ServerID"] = info.ServerID
	dataMap["GiveFriendPoint"] = info.GiveFriendPoint
	dataMap["GetFriendPoint"] = info.GetFriendPoint
	dataMap["Level"] = info.Level
	dataMap["Exp"] = info.Exp
	dataMap["OperationTime"] = info.OperationTime
	_, err = f.client.HSet(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_Friend_Hash, roleID, info.RoleID), dataMap).Result()

	return err
}

// 存储好友信息列表
func (f *Friend) CacheFriends(roleID int64, infos map[int64]*model.FriendInfo) (err error) {
	pipe := f.client.Pipeline()
	for _, info := range infos {
		dataMap := make(map[string]interface{})
		dataMap["RoleID"] = info.RoleID
		dataMap["ServerID"] = info.ServerID
		dataMap["GiveFriendPoint"] = info.GiveFriendPoint
		dataMap["GetFriendPoint"] = info.GetFriendPoint
		dataMap["Level"] = info.Level
		dataMap["Exp"] = info.Exp
		pipe.HSet(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_Friend_Hash, roleID, info.RoleID), dataMap)
	}
	_, err = pipe.Exec(context.TODO())
	return err
}

// 获取好友信息
func (f *Friend) GetFriendInfo(roleID, friendID int64) (info *model.FriendInfo, err error) {
	result, err := f.client.HGetAll(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_Friend_Hash, roleID, friendID)).Result()
	if err != nil {
		return info, err
	}
	givePoint, err := strconv.ParseBool(result["GiveFriendPoint"])
	info = &model.FriendInfo{
		RoleID:          util.StrToInt64(result["RoleID"]),
		ServerID:        util.StrToInt64(result["ServerID"]),
		GiveFriendPoint: givePoint,
		GetFriendPoint:  util.StrToInt32(result["GetFriendPoint"]),
		Level:           util.StrToInt32(result["Level"]),
		Exp:             util.StrToInt32(result["Exp"]),
		OperationTime:   util.StrToInt64(result["OperationTime"]),
	}
	return info, err
}

// 获取好友信息列表
func (f *Friend) GetFriendInfos(roleID int64, ids []int64) (infos map[int64]*model.FriendInfo, err error) {
	pipe := f.client.Pipeline()
	for _, id := range ids {
		pipe.HGetAll(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_Friend_Hash, roleID, id))
	}
	results, err := pipe.Exec(context.TODO())
	if err != nil {
		return infos, err
	}
	infos = make(map[int64]*model.FriendInfo)
	for _, result := range results {
		m, err := result.(*v8.StringStringMapCmd).Result()
		if err != nil {
			return nil, err
		}
		givePoint, err := strconv.ParseBool(m["GiveFriendPoint"])
		if err != nil {
			log.Error("GetFriendInfos ParseBool Error: %v", err)
		}
		info := &model.FriendInfo{
			RoleID:          util.StrToInt64(m["RoleID"]),
			ServerID:        util.StrToInt64(m["ServerID"]),
			GiveFriendPoint: givePoint,
			GetFriendPoint:  util.StrToInt32(m["GetFriendPoint"]),
			Level:           util.StrToInt32(m["Level"]),
			Exp:             util.StrToInt32(m["Exp"]),
			OperationTime:   util.StrToInt64(m["OperationTime"]),
		}
		infos[info.RoleID] = info
	}
	return infos, err
}

// 删除好友
func (f *Friend) DelFriend(roleID, friendID int64) (err error) {
	pipe := f.client.Pipeline()
	pipe.SRem(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_Friend, roleID), friendID)
	pipe.Del(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_Friend_Hash, roleID, friendID))
	_, err = pipe.Exec(context.TODO())
	return err
}

// 已经赠送或者领取过好友点的好友记录
func (f *Friend) CacheFriendPoint(roleID, friendID int64) (err error) {
	pipe := f.client.Pipeline()
	pipe.SAdd(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_Del, roleID), friendID)
	pipe.SAdd(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_Del, friendID), roleID)
	_, err = pipe.Exec(context.TODO())
	return err
}

// 存储黑名单
func (f *Friend) CacheBlackListID(roleID, blackListID int64) (err error) {
	_, err = f.client.SAdd(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_BlackList, roleID), blackListID).Result()
	return err
}

// 是否是黑名单玩家
func (f *Friend) IsBlackList(roleID, blackListID int64) (isBlackList bool, err error) {
	isBlackList, err = f.client.SIsMember(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_BlackList, roleID), blackListID).Result()
	if err != nil {
		isBlackList = false
	}
	return isBlackList, err
}

// 获取所有黑名单玩家
func (f *Friend) GetAllBlackListID(roleID int64) (ids []int64, err error) {
	result, err := f.client.SMembersMap(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_BlackList, roleID)).Result()
	if err != nil {
		return
	}
	ids = make([]int64, 0)
	for friendID := range result {
		ids = append(ids, util.StrToInt64(friendID))
	}
	return ids, err
}

// 移除黑名单
func (f *Friend) DelBlackList(roleID, blackListID int64) (err error) {
	_, err = f.client.SRem(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_BlackList, roleID), blackListID).Result()
	return err
}

// 是否是今天删除过的好友
func (f *Friend) IsDel(roleId, friendId int64) (isDel bool, err error) {
	isDel, err = f.client.SIsMember(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_Del, roleId), friendId).Result()
	if err != nil {
		isDel = false
		return isDel, err
	}
	return isDel, err
}

// 存储已经推荐的好友
func (f *Friend) CacheRecommend(roldId, friendId int64) (err error) {
	key := fmt.Sprintf(model.RedisKey_Friend_Recommend, roldId)
	f.client.SAdd(context.TODO(), key, friendId)
	return err
}

// 是否今天已经推荐过
func IsRecommend(redis *v8.Client, roldId, friendId int64) (isRecommend bool, err error) {
	isRecommend, err = redis.SIsMember(context.TODO(), fmt.Sprintf(model.RedisKey_Friend_Recommend, roldId), friendId).Result()
	if err != nil {
		isRecommend = false
		return isRecommend, err
	}
	return isRecommend, err
}

// 根据等级获取一批玩家id
func (f *Friend) GetPlayerByLevels(uid int64, levels []int32) (roleIds []int64, err error) {
	pipe := f.client.Pipeline()
	for _, level := range levels {
		key := fmt.Sprintf(model.RedisKey_Player_Level, level)
		pipe.SMembers(context.TODO(), key)
	}
	result, err := pipe.Exec(context.TODO())
	roleIds = make([]int64, 0)
	for _, res := range result {
		m, err := res.(*v8.StringSliceCmd).Result()
		if err != nil {
			continue
		}
		for _, id := range m {
			friendID := util.StrToInt64(id)
			isFriend, _ := f.IsFriend(uid, friendID)
			isBlackList, _ := f.IsBlackList(uid, friendID)
			if friendID == uid || isFriend || isBlackList {
				continue
			}
			roleIds = append(roleIds, friendID)
		}
	}
	return roleIds, err
}
