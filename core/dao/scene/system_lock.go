package scene

import (
	"context"
	"fmt"

	"x-server/core/dao/model"
)

// 同步微服务所需的系统解锁信息, id:SystemUnlock.xlsx的解锁ID
func (s *Scene) CacheRoleSystemUnlock(userid int64, id int32) {
	key := fmt.Sprintf(model.RedisKey_SystemUnlock, userid)
	s.client.SAdd(context.TODO(), key, id)
}
