package dao

import (
	"x-server/core/dao/arena"
	"x-server/core/dao/chat"
	"x-server/core/dao/friend"
	"x-server/core/dao/guild"
	"x-server/core/dao/login"
	"x-server/core/dao/push"
	"x-server/core/dao/scene"
	"x-server/core/dao/state"

	v8 "github.com/go-redis/redis/v8"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(New)

var (
	defaultDao Dao
)

type dao struct {
	client *v8.Client
	*arena.Arena
	*chat.Chat
	*friend.Friend
	*guild.Guild
	*login.Login
	*push.Push
	*scene.Scene
	*state.State
}

func New(r *v8.Client) (Dao, func(), error) {
	d := &dao{
		client: r,
		Arena:  arena.New(r),
		Chat:   chat.New(r),
		Friend: friend.New(r),
		Guild:  guild.New(r),
		Login:  login.New(r),
		Push:   push.New(r),
		Scene:  scene.New(r),
		State:  state.New(r),
	}
	defaultDao = d
	return d, d.Close, nil
}

func (d *dao) Close() {

}

// Default 获取接口对象
func Default() Dao {
	return defaultDao
}
