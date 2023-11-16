package state

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	coremodel "x-server/core/model"
	"x-server/core/pkg/util"
	"xy3-proto/pkg/conf/env"
	pbscene "xy3-proto/scene"

	v8 "github.com/go-redis/redis/v8"

	"x-server/core/dao/model"
	pbaccount "xy3-proto/account"
	pb "xy3-proto/coordinator"
	"xy3-proto/pkg/log"
)

// StatefulServiceInfo redis hash field
const (
	StsAppID         = "sts_app_id"          //string
	StsServiceName   = "sts_service_name"    //string
	StsVersion       = "sts_version"         //int64
	StsLastHeartBeat = "sts_last_heartbeat"  //int64
	StsPlayerNum     = "sts_player_quantity" //int64
	StsCPULoadRate   = "sts_cpu_load_rate"   //int64
	StsMemLoadRate   = "sts_mem_load_rate"   //int64
)

// AddServiceAppIdList 服务器app id set
func (s *State) AddServiceAppIdList(ctx context.Context, name, appId string) (err error) {
	return s.client.SAdd(context.TODO(), model.GetServerAppIdKey(name), appId).Err()
}

func (s *State) GetServiceAppIdList(ctx context.Context, name string) (appIdList []string, err error) {
	return s.client.SMembers(context.TODO(), model.GetServerAppIdKey(name)).Result()
}

func (s *State) DelFromServiceAppIdList(ctx context.Context, name, appId string) (err error) {
	return s.client.SRem(context.TODO(), model.GetServerAppIdKey(name), appId).Err()
}

func (s *State) AddServiceList2(ctx context.Context, info *pb.StatefulServiceInfo) (err error) {
	//
	var args []interface{}
	if info.Name != "" {
		args = append(args, StsServiceName, info.Name)
	}
	if info.AppID != "" {
		args = append(args, StsAppID, info.AppID)
	}
	if info.Version > 0 {
		args = append(args, StsVersion, info.Version)
	}
	if info.LastHeartBeat > 0 {
		args = append(args, StsLastHeartBeat, info.LastHeartBeat)
	}
	if info.CPULoadRate > 0 {
		args = append(args, StsCPULoadRate, info.CPULoadRate)
	}
	if info.MemLoadRate > 0 {
		args = append(args, StsMemLoadRate, info.MemLoadRate)
	}
	if info.PlayerNum > 0 {
		args = append(args, StsPlayerNum, info.PlayerNum)
	}

	pp := s.client.Pipeline()
	pp.HSet(context.TODO(), model.GetServerHostsKey2(info.Name, info.AppID), args...)
	pp.SAdd(context.TODO(), model.GetServerAppIdKey(info.Name), info.AppID)
	_, err = pp.Exec(context.TODO())
	if err != nil {
		log.Error("AddServiceList2 pipeline exec service app id:%s, err:%v", info.AppID, err)
		return err
	}
	return err
}

func (s *State) GetServiceInfo(ctx context.Context, name, appId string) (info *pb.StatefulServiceInfo, err error) {
	rs, err := s.client.HGetAll(context.TODO(), model.GetServerHostsKey2(name, appId)).Result()
	if err != nil {
		return
	}

	info = &pb.StatefulServiceInfo{}
	for k, v := range rs {
		switch k {
		case StsAppID:
			info.AppID = v
		case StsServiceName:
			info.Name = v
		case StsVersion:
			info.Version, _ = strconv.ParseInt(v, 10, 64)
		case StsLastHeartBeat:
			info.LastHeartBeat, _ = strconv.ParseInt(v, 10, 64)
		case StsPlayerNum:
			info.PlayerNum, _ = strconv.ParseInt(v, 10, 64)
		case StsCPULoadRate:
			info.CPULoadRate, _ = strconv.ParseInt(v, 10, 64)
		case StsMemLoadRate:
			info.MemLoadRate, _ = strconv.ParseInt(v, 10, 64)
		}
	}
	return
}

func (s *State) GetServiceList2(ctx context.Context, name string) (srvInfos []*pb.StatefulServiceInfo, err error) {
	appIds, err := s.GetServiceAppIdList(context.TODO(), name)
	if err != nil {
		log.Error("GetServiceList2 GetServiceAppIdList err:%v", err)
		return nil, err
	}

	for k := range appIds {
		var info *pb.StatefulServiceInfo
		var err1 error
		info, err1 = s.GetServiceInfo(context.TODO(), name, appIds[k])
		log.Info("GetServiceList2 GetServiceInfo info:%v err:%v", info, err1)
		if err1 != nil {
			log.Error("GetServiceList2 GetServiceInfo err:%v", err1)
			continue
		}
		srvInfos = append(srvInfos, info)
	}
	log.Info("GetServiceList2 %v, err:%v", srvInfos, err)
	return
}

func (s *State) UpdateServiceHeartBeat2(ctx context.Context, req *pb.HeartBeatReq) (err error) {
	info := &pb.StatefulServiceInfo{
		Name:    req.Name,
		AppID:   req.AppID,
		Version: req.Version,
	}

	// 更新心跳
	if req.CPULoadRate > 0 {
		info.CPULoadRate = req.CPULoadRate
	}
	if req.MemLoadRate > 0 {
		info.MemLoadRate = req.MemLoadRate
		info.LastHeartBeat = util.GetTimeStamp()
	}

	err = s.AddServiceList2(context.TODO(), info)
	if err != nil {
		log.Error("[UpdateServiceHeartBeat2] AddServiceList2 hset err:%+v, [value %s]", err, info.String())
		return err
	}
	return nil
}

func (s *State) DeleteServiceState2(ctx context.Context, name string, appID string) (err error) {
	pp := s.client.Pipeline()
	pp.SRem(context.TODO(), model.GetServerAppIdKey(name), appID)
	pp.Del(context.TODO(), model.GetServerHostsKey2(name, appID))
	_, err = pp.Exec(context.TODO())
	if err != nil {
		log.Error("DeleteServiceState2 pipeline Exec app id:%s err:%v", appID, err)
		return
	}
	return
}

// UpdateServicePlayerNum  更新在线人数
// args:: ctx, xy3-1, scene, scene-0, 1
func (s *State) UpdateServicePlayerNum(ctx context.Context, nameSpace, serviceName, appId string, change int64) error {
	// 更新在线用户数量
	s.updateNamespacePlayerLineNum(change)

	// This updates the number of players in coordinator:server:scene:{app-id}
	// E.g: If 1 player enters scene-0, and 1 player enters scene-1,
	// then the number of players in coordinator:server:scene:scene-0 is 1, and the number of players in coordinator:server:scene:scene-1 is 1
	rest1, err1 := s.client.HIncrBy(ctx, model.GetServerHostsKey2(serviceName, appId), StsPlayerNum, change).Result()
	if err1 != nil {
		log.Error("[UpdateServicePlayerNum] HIncrBy GetServerHostsKey2 HIncrBy nameSpace:%s, serviceName:%s, appId:%s, change:%d, err:%v", nameSpace, serviceName, appId, change, err1)
	}
	if rest1 < 0 {
		s.client.HSet(ctx, model.GetServerHostsKey2(serviceName, appId), StsPlayerNum, 0)
	}
	return nil
}

// AddNamespaceSet
// 更新namespace集合 场景服上线时更新
func (s *State) AddNamespaceSet() {
	err := s.client.SAdd(context.TODO(), model.NamespacesKey(), env.Namespace).Err()
	if err != nil {
		log.Error("AddNamespaceSet error:%v", err)
	}
}

// DelNamespaceSet
// 更新namespace集合 场景服下线时更新 TODO 暂时用不到
func (s *State) DelNamespaceSet() {
	err := s.client.SRem(context.TODO(), model.NamespacesKey(), env.Namespace).Err()
	if err != nil {
		log.Error("AddNamespaceSet error:%v", err)
	}
}

// GetAllNamespaces
// 获取所有的namespace
func (s *State) GetAllNamespaces() ([]string, error) {
	res, err := s.client.SMembers(context.TODO(), model.NamespacesKey()).Result()
	if err != nil {
		log.Error("GetAllNamespaces error:%v", err)
		return nil, err
	}
	return res, nil
}

// ClearNamespacePlayerLineNum
// 场景分线重新启动 清空场景分线玩家数量
func (s *State) ClearNamespacePlayerLineNum() {
	err := s.client.HSet(context.TODO(), model.PlayerNumNamespaceKey(env.Namespace), env.AppID, 0).Err()
	if err != nil {
		log.Error("ClearNamespacePlayerLineNum err:%v", err)
	}
}

// 更新 游戏服 各个场景玩家数量 delta : 1 / -1
func (s *State) updateNamespacePlayerLineNum(delta int64) {
	err := s.client.HIncrBy(context.TODO(), model.PlayerNumNamespaceKey(env.Namespace), env.AppID, delta).Err()
	if err != nil {
		log.Error("updateNamespacePlayerLineNum err:%v", err)
	}
}

// DelNamespacePlayerLine
// 去掉场景服分线记录,分线服务器下线时调用
func (s *State) DelNamespacePlayerLine() error {
	log.Info("DelNamespacePlayerLine ns:%v, scene:%v", env.Namespace, env.AppID)
	err := s.client.HDel(context.TODO(), model.PlayerNumNamespaceKey(env.Namespace), env.AppID).Err()
	if err != nil {
		log.Error("DelNamespacePlayerLine ns:%v, scene:%v err:%v", env.Namespace, env.AppID, err)
		return err
	}
	return nil
}

// BatchGetScenePlayerNum
// 获取所有服务器在线人数
func (s *State) BatchGetScenePlayerNum(_ context.Context) (serverPlayerNumList []*pbaccount.ServerPlayerNum, err error) {
	ctx := context.Background()
	ns, err0 := s.GetAllNamespaces()
	if err0 != nil {
		return nil, err0
	}
	pipe := s.client.Pipeline()
	nsM := make(map[string]*v8.StringStringMapCmd)

	for _, v := range ns {
		nsM[v] = pipe.HGetAll(ctx, model.PlayerNumNamespaceKey(v))
	}
	if _, err = pipe.Exec(ctx); err != nil {
		log.Error("BatchGetScenePlayerNum Exec err:%v", err)
	}

	for k, v := range nsM {
		if arr := strings.Split(k, "-"); len(arr) > 1 {
			if serverId, err1 := strconv.Atoi(arr[1]); err1 == nil {
				var sum int64
				res, err1 := v.Result()
				if err1 != nil {
					continue
				}
				// 没有场景分线时，服务有问题
				if len(res) == 0 {
					serverPlayerNumList = append(serverPlayerNumList, &pbaccount.ServerPlayerNum{
						ServerID:  int32(serverId),
						PlayerNum: -1,
					})
					continue
				}
				for kk, vv := range res {
					num, _ := strconv.ParseInt(vv, 10, 64)
					if num > 0 {
						sum += num
					} else {
						// 检查服务器是否可用， 或者维护
						if s.checkSceneAvailableOrMaintain(k, kk) == -1 {
							sum = -1
							break
						}
					}
				}
				serverPlayerNumList = append(serverPlayerNumList, &pbaccount.ServerPlayerNum{
					ServerID:  int32(serverId),
					PlayerNum: int32(sum),
				})
			}
		}
	}
	return serverPlayerNumList, nil
}

// BatchGetScenePlayerNumV1
// fix 效率问题
// 获取所有服务器在线人数
//
//nolint:all
func (s *State) BatchGetScenePlayerNumV1(_ context.Context) ([]*pbaccount.ServerPlayerNum, error) {
	ns, err0 := s.GetAllNamespaces()
	if err0 != nil {
		log.Error("BatchGetScenePlayerNumV1 GetAllNamespaces err:%v", err0)
		return nil, err0
	}
	serverPlayerNumList := make([]*pbaccount.ServerPlayerNum, 0, len(ns))
	if len(ns) == 0 {
		return serverPlayerNumList, nil
	}
	servers := make([]int32, 0, len(ns))
	for _, v := range ns {
		if arr := strings.Split(v, "-"); len(arr) > 1 {
			serverId, err1 := strconv.ParseInt(arr[1], 10, 64)
			if err1 != nil {
				log.Warn("BatchGetScenePlayerNumV1: ns:%v parse serverId err: %v", v, err1)
				continue
			}
			servers = append(servers, int32(serverId))
		}
	}

	serverLen := len(servers)
	if serverLen == 0 {
		return serverPlayerNumList, nil
	}
	ch := make(chan *pbaccount.ServerPlayerNum, len(servers))
	defer close(ch)

	for _, v := range servers {
		go func(c chan *pbaccount.ServerPlayerNum, serverId int32) {
			s.getServerPlayerNum(c, serverId)
		}(ch, v)
	}

	var stop bool
	for !stop {
		select {
		case v, ok := <-ch:
			if ok {
				serverPlayerNumList = append(serverPlayerNumList, v)
			}
		}
		if len(serverPlayerNumList) == serverLen {
			stop = true
		}
	}
	return serverPlayerNumList, nil
}

// 获取服务器玩家数量
// -1  服务器不可用
func (s *State) getServerPlayerNum(ch chan *pbaccount.ServerPlayerNum, serverId int32) {
	ctx := context.Background()
	res := &pbaccount.ServerPlayerNum{
		ServerID: serverId,
	}
	defer func() {
		ch <- res
	}()
	// lines: {"scene-0":"1", "scene-1": "0"}
	lines, err := s.client.HGetAll(ctx, model.PlayerNumNamespaceKey(fmt.Sprintf(coremodel.NamespaceV, serverId))).Result()
	if err != nil {
		log.Warn("getServerPlayerNum HGetAll serverId:%v error:%v", serverId, err)
	}
	//  scene-0, 0
	for lineId, num0 := range lines {
		num, _ := strconv.ParseInt(num0, 10, 64)
		if num > 0 {
			res.PlayerNum += int32(num)
			continue
		}
		// 玩家数量为0时, 检查line是否可用
		if state := s.checkSceneAvailableOrMaintain(fmt.Sprintf(coremodel.NamespaceV, serverId), lineId); state == -1 {
			res.PlayerNum = -1
			break
		}
	}
}

// 检查服务器场景服是否可用， -1 服务器不可用， -2 维护
func (s *State) checkSceneAvailableOrMaintain(ns string, lineId string) int {
	// TODO  加维护接口
	// 检查服务器是否可用， 或者维护
	client := pbscene.DefaultCrossClient(ns, lineId)
	if client == nil {
		log.Error("checkSceneAvailableOrMaintain NewCrossClient ns:%v lineId:%v", ns, lineId)
		return -1
	}
	ctx, clear := context.WithTimeout(context.Background(), time.Second)
	defer clear()
	_, err1 := client.Ping(ctx, &pbscene.Empty{})
	if err1 != nil {
		log.Error("checkSceneAvailableOrMaintain Ping ns:%v linId:%v err:%v", ns, lineId, err1)
		return -1
	}
	return 0
}

// GetScenePlayerNum fix
// 获取某个服务器在线人数
func (s *State) GetScenePlayerNum(ctx context.Context, namespace string) (num int, err error) {
	res, err1 := s.client.HGetAll(ctx, model.PlayerNumNamespaceKey(namespace)).Result()
	if err1 != nil {
		log.Error("GetScenePlayerNum key:%v field:%v err:%v", model.PlayerNumNamespaceKey(namespace), namespace, err1)
		return 0, err1
	}
	for _, v := range res {
		if playerNum, err := strconv.Atoi(v); err == nil {
			num += playerNum
		}
	}
	return
}

func (s *State) ResetScenePlayerNum(ctx context.Context, namespace string) (err error) {
	_, err = s.client.HSet(ctx, model.GetScenePlayerNum(), namespace, 0).Result()
	if err != nil {
		log.Error("ResetScenePlayerNum key:%v field:%v err:%v", model.GetScenePlayerNum, namespace, err)
	}
	return nil
}

func (s *State) LoadBalance(stslist []*pb.StatefulServiceInfo) *pb.StatefulServiceInfo {
	if len(stslist) == 0 {
		log.Info("[LoadBalance] len(list) == 0")
		return nil
	}
	// 去除无效的信息
	for k, v := range stslist {
		if v.Name == "" || v.AppID == "" {
			stslist = append(stslist[:k], stslist[k+1:]...)
		}
	}
	if len(stslist) == 0 {
		return stslist[0]
	}

	sort.SliceStable(stslist, func(i, j int) bool {
		return stslist[i].CPULoadRate < stslist[j].CPULoadRate
	})

	target := stslist[0]
	for _, info := range stslist {
		if info.CPULoadRate >= model.CPULoadRateLimit {
			continue
		}
		if target.Version > info.Version {
			continue
		} else if target.Version < info.Version {
			target = info
		} else if target.PlayerNum > info.PlayerNum {
			target = info
		}
	}
	log.Info("[LoadBalance] 2 target: %v", target)
	return target
}
