package dao

import (
	"context"
	"fmt"
	"time"
	"xy3-proto/pkg/log"
	pbworld "xy3-proto/world"
)

func (d *dao) GetRoleInfo(uid int64) *pbworld.RoleExInfo {
	roleInfo := d.Scene.GetCacheRole(uid)
	if roleInfo == nil {
		log.Error("GetRoleInfo uid:%v", uid)
		return nil
	}
	pbInfo := &pbworld.RoleExInfo{
		RoleId:   roleInfo.ID,
		RoleName: roleInfo.Nick,
		Level:    roleInfo.Level,
		HeadID:   roleInfo.HeadID,
		FrameID:  roleInfo.FrameID,
		DrawID:   roleInfo.DrawID,
		Title:    roleInfo.Title,
		Sex:      roleInfo.Sex,
		ServerID: roleInfo.ServerID,
		Power:    roleInfo.Power,
	}

	if roleInfo.IsRobot == 0 {
		guilduser := d.Guild.GetGuildUser(uid)
		if guilduser != nil && guilduser.GuildID != 0 {
			guildinfo := d.Guild.GetGuildInfo(guilduser.GuildID)
			if guildinfo != nil {
				pbInfo.GuildID = guilduser.GuildID
				pbInfo.GuildName = guildinfo.Name
			}
		}
	}
	return pbInfo
}

func (d *dao) GetRoleInfos(uids []int64) map[int64]*pbworld.RoleExInfo {
	infos := d.Scene.GetMutliPlayerInfo(uids)
	roleInfoMap := map[int64]*pbworld.RoleExInfo{}
	for uid, info := range infos {
		pbInfo := &pbworld.RoleExInfo{
			RoleId:   info.ID,
			RoleName: info.Nick,
			Level:    info.Level,
			HeadID:   info.HeadID,
			FrameID:  info.FrameID,
			DrawID:   info.DrawID,
			Title:    info.Title,
			Sex:      info.Sex,
			ServerID: info.ServerID,
			Power:    info.Power,
		}
		if info.IsRobot == 0 {
			guilduser := d.Guild.GetGuildUser(uid)
			if guilduser != nil && guilduser.GuildID != 0 {
				guildinfo := d.Guild.GetGuildInfo(guilduser.GuildID)
				if guildinfo != nil {
					pbInfo.GuildID = guilduser.GuildID
					pbInfo.GuildName = guildinfo.Name
				}
			}
		}
		roleInfoMap[info.ID] = pbInfo
	}
	return roleInfoMap
}

// GetRes
// 提供频率限制
func (d *dao) GetRes(uniqueKey string, expire ...time.Duration) (bool, error) {
	// 默认超时时间为10毫秒
	t := time.Millisecond * 10
	if len(expire) > 0 {
		t = expire[0]
	}
	return d.client.SetNX(context.Background(), fmt.Sprintf("getRes:%v", uniqueKey), 1, t).Result()
}
