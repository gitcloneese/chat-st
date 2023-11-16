package model

import (
	"fmt"
	env "xy3-proto/pkg/conf/env"
)

const (
	RedisServiceListKey2 = "coordinator:server:%s:%s:%s:hash" // 拉取某个服务的所有实例列表,

	RedisServiceAppIdListKey = "coordinator:appid:%s:%s:set" // 拉取某个服务的所有实例列表,

	//
	RedisUserServerListKey = "coordinator:player:server:%d:hash" // 用户分配的服务器列表

	RedisServiceUserKey = "coordinator:players:%s:%s:%s:set" // 服务器关联的用户列表

	RedisKeyPlayerLoginLockKey = "lock:login:player:%s" // sdk uuid

	RedisKeyPlayerStateLockKey = "lock:state:player:%d"

	DefaultRedLockTTL = 3000

	RedisScenePlayerNum = "playerNum:scene" // 场景服人数

	// Namespaces 游戏服id 集合
	namespaces = "namespaces"

	playerNumNamespaceKey = "playerNum:%v" // player:xy3-1 -> {"scene-0":1,"scene-1":1}

)

var AllocServerList = []string{"scene"}

type ServerNumStruct struct {
	ServerId  int
	PlayerNum int
}

const AllocTimeOut = 600
const CPULoadRateLimit = 75

//func GetServerHostsKey(serviceName string) string {
//	return fmt.Sprintf(RedisServiceListKey, env.Namespace, serviceName)
//}

// 某个服务的实例列表
func GetServerAppIdKey(serviceName string) string {
	return fmt.Sprintf(RedisServiceAppIdListKey, env.Namespace, serviceName)
}

// 某个实例的状态
func GetServerHostsKey2(serviceName, appId string) string {
	return fmt.Sprintf(RedisServiceListKey2, env.Namespace, serviceName, appId)
}

// 某个服务实例的玩家
func GetServiceUserKey(serviceName, appId string) string {
	return fmt.Sprintf(RedisServiceUserKey, env.Namespace, serviceName, appId)
}

// 用户关联的服务器列表
func GetUserServerListKey(playerID int64) string {
	return fmt.Sprintf(RedisUserServerListKey, playerID)
}

// 用户的状态
func GetPlayerStateKey(playerID int64) string {
	return fmt.Sprintf(RedisUserStateKey, playerID)
}

func GetPlayerStateLockKey(playerId int64) string {
	return fmt.Sprintf(RedisKeyPlayerStateLockKey, playerId)
}

func GetPlayerLoginLockKey(sdkAccountId string) string {
	return fmt.Sprintf(RedisKeyPlayerLoginLockKey, sdkAccountId)
}

func GetScenePlayerNum() string {
	return RedisScenePlayerNum
}

// NamespacesKey
// 游戏服集合
func NamespacesKey() string {
	return namespaces
}

// PlayerNumNamespaceKey
// 获取服务器玩家数量key
func PlayerNumNamespaceKey(namespace string) string {
	return fmt.Sprintf(playerNumNamespaceKey, namespace)
}
