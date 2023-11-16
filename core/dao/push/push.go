package push

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	coremodel "x-server/core/model"
	"xy3-proto/pkg/log"
	pbscene "xy3-proto/scene"

	v8 "github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"
)

// 下面这三个表紧密关联调度服(coordinator:constants.go)的表
const (
	RedisPushSceneList        string = "coordinator:appid:%v:scene:set" //set:{scene-x}
	RedisKey_Player_Line      string = "player_line:%v"                 //玩家所在分线，比如：scene-0， 过期时间8min，玩家每5min写入一次
	RedisPushPlayer           string = "player:%v"
	DefaultPlayerLineServerId string = "scene-0" // 默认场景服id, 玩家不在线时, 由默认场景服务处理
)

type Push struct {
	client *v8.Client
}

func New(r *v8.Client) *Push {
	return &Push{
		client: r,
	}
}

// PushScene 推送通知：推送给某个服务器ID下的所有scene分线. 例如: dao.PushScene(1, pbscene.PushMessageID_PushMailNotice, &pbactivity.Empty{})
func (p *Push) PushScene(serverid int32, pushid int32, msg proto.Message) {
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Error("push scene marshal err! pushid:%v", pushid)
		return
	}

	p.PushSceneData(serverid, pushid, data)
}

func (p *Push) PushSceneData(serverid int32, pushid int32, data []byte) {
	log.Info("push scene data! serverid:%v pushid:%v", serverid, pushid)

	scenes := p.querySceneServerList(serverid)
	log.Debug("PushSceneData scenes list:%+v", scenes)

	for _, name := range scenes {
		if err := pbscene.BroadcastCrossMsg(context.Background(), fmt.Sprintf(coremodel.NamespaceV, serverid), name, []int64{}, pushid, data); err != nil {
			log.Error("push scene:%v err! pushid:%v err:%v", name, pushid, err)
		}
	}
}

func (p *Push) querySceneServerList(serverId int32) (list []string) {
	list = []string{}

	key := fmt.Sprintf(RedisPushSceneList, fmt.Sprintf(coremodel.NamespaceV, serverId))
	result, err := p.client.SMembers(context.Background(), key).Result()
	if err != nil {
		log.Error("querySceneServerList err! err:%v", err)
		return list
	}

	for _, k := range result {
		if strings.HasPrefix(k, "scene") {
			list = append(list, k)
		}
	}
	return list
}

// 推送通知：ids为空表示推送本服通知，否则推送给具体的指定玩家. 例如：dao.PushUser([]int64{190860, 127648}, pbscene.PushMessageID_PushMailNotice, &pbactivity.Empty{})
func (p *Push) PushUser(ctx context.Context, ids []int64, pushid int32, msg proto.Message) {
	log.Info("PushUser! pushid:%v", pushid)
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Errorc(ctx, "push err! msg marshal err! pushid:%v", pushid)
		return
	}

	p.PushUserData(ctx, ids, pushid, data)
}

func (p *Push) PushUserData(traceCtx context.Context, ids []int64, pushid int32, data []byte) {
	ctx := context.Background()
	maddr := p.queryUserServerInfo(ctx, ids)
	log.Debug("push user data! player id:%v addr:%v pushid:%v", ids, maddr, pushid)

	for serverid, mm := range maddr {
		if serverid <= 0 { // invalid server id , possible is a robot,  UPDATE: 24OCT23, robots are cross servers so serverid init as 0,
			continue
		}
		for name, v := range mm {
			if err := pbscene.BroadcastCrossMsg(ctx, fmt.Sprintf(coremodel.NamespaceV, serverid), name, v, pushid, data); err != nil {
				log.Errorc(traceCtx, "push user:%v err! pushid:%v err:%v", name, pushid, err)
			}
		}
	}
}

// maddr: <所在服索引, <所在场景分线, 玩家ID列表>>
//
//nolint:funlen
func (p *Push) queryUserServerInfo(ctx context.Context, ids []int64) (maddr map[int32]map[string][]int64) {
	maddr = make(map[int32]map[string][]int64)
	pipe := p.client.Pipeline()
	m := make(map[int64][]*v8.StringCmd)
	for _, userId := range ids {
		//取得玩家所在服ID
		m[userId] = append(m[userId], pipe.HGet(ctx, fmt.Sprintf(RedisPushPlayer, userId), "ServerID"))
		//取得对应的场景服分线
		m[userId] = append(m[userId], pipe.Get(ctx, fmt.Sprintf(RedisKey_Player_Line, userId)))
	}

	if _, err := pipe.Exec(ctx); err != nil {
		log.Warn("queryUserServerInfo Exec err:%v", err)
	}

	for uid, commands := range m {
		server, err1 := commands[0].Result()
		if err1 != nil {
			log.Error("queryUserServerInfo server Result error:%v", err1)
			continue
		}
		serverid0, err2 := strconv.ParseInt(server, 10, 64)
		if err2 != nil {
			log.Error("queryUserServerInfo server ParseInt err:%v", err2)
			continue
		}
		serverId := int32(serverid0)
		if _, ok := maddr[serverId]; !ok {
			maddr[serverId] = make(map[string][]int64)
		}

		appid, err3 := commands[1].Result()
		if err3 != nil && !errors.Is(err3, v8.Nil) {
			log.Error("queryUserServerInfo app result err:%v", err3)
			continue
		}
		// 玩家不在线时， 由默认场景夫处理scene-0
		if errors.Is(err3, v8.Nil) {
			maddr[serverId][DefaultPlayerLineServerId] = append(maddr[serverId][DefaultPlayerLineServerId], uid)
			continue
		}

		if appid != "" {
			maddr[serverId][appid] = append(maddr[serverId][appid], uid)
		}
	}

	// 以上代码 重写下面代码， 多个redis req 合并为一个 pipeline

	//}
	//
	//for _, userid := range ids {
	//	var mm map[string][]int64
	//	//取得玩家所在服ID
	//	key := fmt.Sprintf(RedisPushPlayer, userid)
	//	result, err := p.client.HGet(ctx, key, "ServerID").Result()
	//	if err != nil {
	//		log.Warn("queryUserServerInfo HGet err:%v player id:%d", err, userid)
	//		continue
	//	}
	//	serverid, err := strconv.ParseInt(result, 10, 64)
	//	if err != nil {
	//		log.Error("queryUserServerInfo ParseInt err:%v player id:%d", err, userid)
	//		continue
	//	}
	//	_, ok := maddr[int32(serverid)]
	//	if !ok {
	//		maddr[int32(serverid)] = make(map[string][]int64)
	//	}
	//	mm = maddr[int32(serverid)]
	//
	//	//取得对应的场景服分线
	//	key = fmt.Sprintf(RedisKey_Player_Line, userid)
	//	appid, err := p.client.Get(ctx, key).Result()
	//	if err != nil {
	//		log.Debug("queryUserServerInfo player line err! userid:%v err:%v", userid, err)
	//		err = nil
	//		continue
	//	}
	//	list, ok := mm[appid]
	//	if !ok {
	//		mm[appid] = []int64{userid}
	//	} else {
	//		mm[appid] = append(list, userid)
	//	}
	//}
	return maddr
}
