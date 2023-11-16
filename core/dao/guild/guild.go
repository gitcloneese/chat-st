package guild

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"x-server/core/dao/model"
	"x-server/core/pkg/util"
	coreutil "x-server/core/pkg/util"
	"xy3-proto/pkg/conf/env"
	"xy3-proto/pkg/log"

	v8 "github.com/go-redis/redis/v8"
)

type Guild struct {
	client *v8.Client
}

func New(r *v8.Client) *Guild {
	return &Guild{
		client: r,
	}
}

func KeyGuildInfo(guildid int64) string {
	return fmt.Sprintf(model.RedisGuildInfo, env.Namespace, guildid)
}

func KeyGuildMembers(guildid int64) string {
	return fmt.Sprintf(model.RedisGuildMembers, env.Namespace, guildid)
}

func KeyGuildUserInfo(userid int64) string {
	return fmt.Sprintf(model.RedisGuildUserInfo, env.Namespace, userid)
}

func KeyGuildHunts(guildid int64) string {
	return fmt.Sprintf(model.RedisGuildHunts, env.Namespace, guildid)
}

func KeyMemberRedPackets(userid int64) string {
	return fmt.Sprintf(model.RedisMemberRedPackets, env.Namespace, userid)
}

func KeyRedPacket(uuid int64) string {
	return fmt.Sprintf(model.RedisRedPacketInfo, env.Namespace, uuid)
}

func (g *Guild) HasGuildName(name string) bool {
	ok, _ := g.client.SIsMember(context.TODO(), fmt.Sprintf(model.RedisGuildNames, env.Namespace), name).Result()
	return ok
}

// IsExistGuild 是否存在公会
func (g *Guild) IsExistGuild(guildid int64) bool {
	key := KeyGuildInfo(guildid)
	res := g.client.Exists(context.TODO(), key)
	return res.Val() == 1
}

// GetGuildInfo 取得redis中的公会信息
func (g *Guild) GetGuildInfo(guildid int64) (info *model.GuildInfo) {
	key := KeyGuildInfo(guildid)
	result, err := g.client.HGetAll(context.TODO(), key).Result()
	if err != nil {
		log.Error("GetGuildInfo redis-err! err:%v", err)
		return nil
	}
	if len(result) == 0 {
		return nil
	}
	info = &model.GuildInfo{
		ID:           guildid,
		Name:         result["Name"],
		Notice:       result["Notice"],
		Icon:         util.StrToInt32(result["Icon"]),
		HallLv:       util.StrToInt32(result["HallLv"]),
		TowerLv:      util.StrToInt32(result["TowerLv"]),
		PavilionLv:   util.StrToInt32(result["PavilionLv"]),
		ShopLv:       util.StrToInt32(result["ShopLv"]),
		Fund:         util.StrToInt32(result["Fund"]),
		FightLimit:   util.StrToInt64(result["FightLimit"]),
		Chairman:     util.StrToInt64(result["Chairman"]),
		CoolName:     util.StrToInt64(result["CoolName"]),
		CoolNotice:   util.StrToInt64(result["CoolNotice"]),
		CoolMail:     util.StrToInt64(result["CoolMail"]),
		CoolEnlist:   util.StrToInt64(result["CoolEnlist"]),
		ImpeachST:    util.StrToInt64(result["ImpeachST"]),
		ImpeachState: util.StrToInt32(result["ImpeachState"]),
		LiveDay:      util.StrToInt32(result["LiveDay"]),
		LiveWeek:     util.StrToInt32(result["LiveWeek"]),
		Inactive:     util.StrToBool(result["Inactive"]),
	}
	return info
}

func (g *Guild) UpdateGuildInfo(info *model.GuildInfo) (err error) {
	key := KeyGuildInfo(info.ID)
	data := map[string]interface{}{
		"ID":            info.ID,
		"Name":          info.Name,
		"Notice":        info.Notice,
		"Icon":          info.Icon,
		"HallLv":        info.HallLv,
		"TowerLv":       info.TowerLv,
		"PavilionLv":    info.PavilionLv,
		"ShopLv":        info.ShopLv,
		"Fund":          info.Fund,
		"FightLimit":    info.FightLimit,
		"Chairman":      info.Chairman,
		"CoolName":      info.CoolName,
		"CoolNotice":    info.CoolNotice,
		"CoolMail":      info.CoolMail,
		"CoolEnlist":    info.CoolEnlist,
		"ImpeachST":     info.ImpeachST,
		"ImpeachState":  info.ImpeachState,
		"LiveDay":       info.LiveDay,
		"LiveWeek":      info.LiveWeek,
		"MSoul":         info.MSoul,
		"LastFlushTm":   info.LastFlushTm,
		"AutoAgreeJoin": info.AutoAgreeJoin,
		"TBoxNum":       info.TBoxNum,
		"Inactive":      info.Inactive,
	}
	_, err = g.client.HSet(context.TODO(), key, data).Result()
	return err
}

// GetGuildUser 取得redis中的个人关于公会的信息
func (g *Guild) GetGuildUser(userid int64) (info *model.GuildUser) {
	key := KeyGuildUserInfo(userid)
	result, err := g.client.HGetAll(context.TODO(), key).Result()
	if err != nil {
		log.Error("GetGuildUser redis-err! err:%v", err)
		return nil
	}
	if len(result) == 0 {
		return nil
	}
	info = &model.GuildUser{
		ID:                  userid,
		GuildID:             util.StrToInt64(result["GuildID"]),
		Job:                 util.StrToInt32(result["Job"]),
		DayLive:             util.StrToInt32(result["DayLive"]),
		WeekLive:            util.StrToInt32(result["WeekLive"]),
		TotalLive:           util.StrToInt32(result["TotalLive"]),
		CoolKick:            util.StrToInt64(result["CoolKick"]),
		CoolExit:            util.StrToInt64(result["CoolExit"]),
		RedPGetedNum:        util.StrToInt32(result["RedPGetedNum"]),
		RedPGetedRes:        util.StrToInt32(result["RedPGetedRes"]),
		RedPGiveNum:         util.StrToInt32(result["RedPGiveNum"]),
		RefuseAllInviteFlag: util.StrToInt32(result["RefuseAllInviteFlag"]),
		DissolveFlag:        util.StrToInt32(result["DissolveFlag"]),
		LiveDayRewardFlag:   util.StrToInt32(result["LiveDayRewardFlag"]),
		AuctionBuyNum:       util.StrToInt32(result["AuctionBuyNum"]),
		AuctionCost:         util.StrToInt64(result["AuctionCost"]),
	}
	return
}

func (g *Guild) UpdateGuildUser(info *model.GuildUser) (err error) {
	key := KeyGuildUserInfo(info.ID)
	data := map[string]interface{}{
		"ID":                  info.ID,
		"GuildID":             info.GuildID,
		"Job":                 info.Job,
		"TotalLive":           info.TotalLive,
		"WeekLive":            info.WeekLive,
		"DayLive":             info.DayLive,
		"TaskCnt":             info.TaskCnt,
		"TaskNum":             info.TaskNum,
		"CoolKick":            info.CoolKick,
		"CoolExit":            info.CoolExit,
		"LastFlushTm":         info.LastFlushTm,
		"RedPGetedNum":        info.RedPGetedNum,
		"RedPGetedRes":        info.RedPGetedRes,
		"RedPGiveNum":         info.RedPGiveNum,
		"RefuseAllInviteFlag": info.RefuseAllInviteFlag,
		"DissolveFlag":        info.DissolveFlag,
		"LiveDayRewardFlag":   info.LiveDayRewardFlag,
		"AuctionBuyNum":       info.AuctionBuyNum,
		"AuctionCost":         info.AuctionCost,
	}
	_, err = g.client.HSet(context.TODO(), key, data).Result()
	return
}

func GetGuildMemberNum(redis *v8.Client, guildid int64) int32 {
	key := KeyGuildMembers(guildid)
	num, _ := redis.SCard(context.TODO(), key).Result()
	return int32(num)
}

func (g *Guild) ListIDs(guildId int64, nameSpace string) (ids []int64) {
	ids = []int64{}
	key := fmt.Sprintf(model.RedisGuildMembers, nameSpace, guildId)
	result, err := g.client.SMembers(context.TODO(), key).Result()
	if err != nil {
		log.Error("GetGuildMembers err:[%v]", err)
		return
	}
	for _, v := range result {
		ids = append(ids, util.StrToInt64(v))
	}
	return
}

func (g *Guild) HasGuildHighHuntTI(guildid int64, ti int64) bool {
	key := KeyGuildHunts(guildid)
	ok, _ := g.client.SIsMember(context.TODO(), key, ti).Result()
	return ok
}

func (g *Guild) GetRedPacketNum(userid int64) int32 {
	key := KeyMemberRedPackets(userid)
	num, _ := g.client.SCard(context.TODO(), key).Result()
	return int32(num)
}

func (g *Guild) GetRedPacket(uuid int64) (info *model.CacheRedPacket) {
	key := KeyRedPacket(uuid)
	result, err := g.client.HGetAll(context.TODO(), key).Result()
	if err != nil {
		return nil
	}
	info = &model.CacheRedPacket{
		UUID: uuid,
		ID:   util.StrToInt32(result["ID"]),
	}
	return
}

func KeyAuctionGoods(uuid int64) string {
	return fmt.Sprintf(model.RedisAuctionGoods, env.Namespace, coreutil.GetDay(), uuid)
}

func (g *Guild) GetAuctionGoods(uuid int64) (info *model.CacheAuctionGoods) {
	key := KeyAuctionGoods(uuid)

	result, err := g.client.HGetAll(context.TODO(), key).Result()
	if err != nil || len(result) == 0 {
		log.Error("GetAuctionGoods err! err:%v", err)
		return nil
	}

	info = &model.CacheAuctionGoods{
		UUID:       uuid,
		ID:         util.StrToInt32(result["ID"]),
		ActivityID: util.StrToInt32(result["ActivityID"]),
		ShareRate:  util.StrToInt32(result["ShareRate"]),
	}
	info.ListFromString(result["List"])
	return
}

func KeyAuctionBider(uuid int64) string {
	return fmt.Sprintf(model.RedisAuctionBider, env.Namespace, coreutil.GetDay(), uuid)
}

func (g *Guild) GetAuctionBider(uuid int64) (info *model.CacheAuctionBider) {
	key := KeyAuctionBider(uuid)

	result, err := g.client.HGetAll(context.TODO(), key).Result()
	if err != nil || len(result) == 0 {
		log.Error("GetAuctionBider err! err:%v", err)
		return nil
	}

	info = &model.CacheAuctionBider{
		UUID:      uuid,
		UserID:    util.StrToInt64(result["UserID"]),
		GoodsUUID: util.StrToInt64(result["GoodsUUID"]),
		Bidprice:  util.StrToInt32(result["Bidprice"]),
		UnixTime:  util.StrToInt64(result["UnixTime"]),
	}
	return
}

func KeyAuctionBidback(userid int64) string {
	return fmt.Sprintf(model.RedisAuctionBidback, env.Namespace, userid)
}

func (g *Guild) GetAuctionBidback(userid int64) (int64, error) {
	key := KeyAuctionBidback(userid)

	if str, err := g.client.Get(context.TODO(), key).Result(); err != nil {
		if errors.Is(err, v8.Nil) {
			return 0, nil
		} else {
			return 0, err
		}
	} else {
		val, _ := strconv.ParseInt(str, 10, 64)
		return val, nil
	}
}

func KeyAuctionLastBiddingFailedTime(userid int64) string {
	key := fmt.Sprintf(model.RedisAuctionBiddingFailed, env.Namespace, userid)
	return key
}

func (g *Guild) AuctionLastBiddingFailedTime(userid int64) int64 {
	key := KeyAuctionLastBiddingFailedTime(userid)
	if str, err := g.client.Get(context.TODO(), key).Result(); err != nil {
		return 0
	} else {
		val, _ := strconv.ParseInt(str, 10, 64)
		return val
	}
}

func (g *Guild) UpdateAuctionBidback(userid int64, val int64) error {
	key := KeyAuctionBidback(userid)
	_, err := g.client.Set(context.TODO(), key, fmt.Sprintf("%v", val), 0).Result()
	return err
}

func KeyAuctionLastUUID(userid int64) string {
	return fmt.Sprintf(model.RedisAuctionPlayerLastUUID, env.Namespace, userid)
}
func (g *Guild) GetAuctionLastUUID(userid int64) int64 {
	key := KeyAuctionLastUUID(userid)
	if str, err := g.client.Get(context.TODO(), key).Result(); err != nil {
		return 0
	} else {
		uuid, _ := strconv.ParseInt(str, 10, 64)
		return uuid
	}
}

func KeyDevilConquerAgentInfo(userid int64) string {
	return fmt.Sprintf(model.RedisDevilConquerAgentInfo, env.Namespace, userid)
}

func (g *Guild) GetDCAgentInfo(userid int64) (info *model.CacheDCAgentInfo) {
	key := KeyDevilConquerAgentInfo(userid)

	result, err := g.client.HGetAll(context.TODO(), key).Result()
	if err != nil {
		log.Error("GetCacheDCAgentInfo err! err:%v", err)
		return nil
	}
	if len(result) == 0 {
		return nil
	}

	info = &model.CacheDCAgentInfo{
		UserID:       userid,
		Score:        util.StrToInt64(result["Score"]),
		CoolBoss:     util.StrToInt64(result["CoolBoss"]),
		CoolRob:      util.StrToInt64(result["CoolRob"]),
		CoolBeRobbed: util.StrToInt64(result["CoolBeRobbed"]),
	}
	return
}

func KeyDevilConquerAgentPreRank(userid int64) string {
	return fmt.Sprintf(model.RedisDevilConquerAgentPreRank, env.Namespace, userid)
}

func (g *Guild) UpdateDCAgentPreRank(userid int64, rank int32) error {
	key := KeyDevilConquerAgentPreRank(userid)
	_, err := g.client.Set(context.TODO(), key, fmt.Sprintf("%v", rank), 0).Result()
	return err
}

func (g *Guild) GetDCAgentPreRank(userid int64) int32 {
	key := KeyDevilConquerAgentPreRank(userid)

	if str, err := g.client.Get(context.TODO(), key).Result(); err != nil {
		return 0
	} else {
		val, _ := strconv.ParseInt(str, 10, 64)
		return int32(val)
	}
}

func KeyDungeon(activityid int32) string {
	return fmt.Sprintf(model.RedisDungeonPlayers, env.Namespace, activityid)
}

func (g *Guild) IsInDungeon(activityid int32, userid int64) bool {
	key := KeyDungeon(activityid)
	ok, _ := g.client.SIsMember(context.TODO(), key, userid).Result()
	return ok
}
