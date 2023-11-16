package dao

import (
	"context"
	"time"
	"x-server/core/dao/model"
	coremodel "x-server/core/model"
	pbaccount "xy3-proto/account"
	battle "xy3-proto/battle"
	pb "xy3-proto/coordinator"
	pbrank "xy3-proto/rank"
	pbscene "xy3-proto/scene"
	pbworld "xy3-proto/world"

	"google.golang.org/protobuf/proto"
)

type Dao interface {
	//////////////////////Arena相关接口
	// 记录斗法每日排行
	SaveDailyRank(serverID int64, pvpType int32, day string, dailyRankMap map[int64]int32) (err error)
	// 获取玩家斗法每日排行
	GetDailyRankByUid(serverID int64, pvpType int32, day string, uid int64) (rank int32, err error)
	// 删除斗法每日排行
	DelDailyRank(serverID int64, pvpType int32, day string) (err error)
	// 记录斗法成就排行
	SaveAchievementRank(serverID int64, pvpType int32, achievementRankMap map[int64]int32) (err error)
	// 获取斗法成就排行
	GetAchievementRank(serverID int64, pvpType int32) (achievementRankMap map[int64]int32, err error)
	// 获取斗法单人成就排行
	GetAchievementRankByUid(serverID int64, pvpType int32, uid int64) (rank int32, err error)
	IsPvpInit(serverID int64) (b bool, err error)
	GetPvpTime(serverID int64) (time int64, err error)
	SavePvpTime(serverID int64, time int64)
	IsNeedEveryRefresh(serverID int64) (b bool, err error)
	GetPvpFighting(serverID int64, roleID int64) (pvpFighting *model.PvpFightingReport, err error)
	SavePvpFighting(serverID int64, roleID int64, pvpFighting *model.PvpFightingReport) (err error)
	// 斗法挑战锁定双方排行
	PvpFightingLocked(serverID int64, pvpLevel int32, selfRank int32, targetRank int32) (islock bool, err error)
	// 斗法挑战完毕解锁双方排行
	PvpFightingUnlock(serverID int64, pvpLevel int32, selfRank int32, targetRank int32) (err error)
	// AppendPvpRank 追加
	AppendPvpRank(serverID int64, pvpLevel int64, list []*pbrank.RankData) (appendedRank int64, err error)
	// ReplaceRank 替换双方排行
	ReplacePvpRank(serverID int64, pvpLevel int64, uniqueID int64, targetUniqueID int64, isTargetRobot bool) (err error)
	// 玩家排行晋升
	UpPvpLevel(serverID, pvpLevel, playerID int64) (err error)
	// GetPvpRankInfo 查询唯一id对应的排行信息
	GetPvpRankInfo(serverID int64, pvpLevel int64, uniqueID int64) (info *pbrank.RankData, err error)
	//获得排行列表
	GetPvpRankList(serverID int64, pvpLevel int64, start, stop int64) (list []*pbrank.RankData, err error)
	// 获取单个玩家基本信息
	GetTTPvpPlayerInfo(seasonID int32, id int64) (info *model.TTPvpPlayerInfo, err error)
	// 获取多个玩家基本信息
	GetMutliTTPvpPlayerInfo(seasonId int32, ids []int64) (infos map[int64]*model.TTPvpPlayerInfo, err error)
	// 更新玩家基本信息
	UpdateTTPvpPlayerInfo(seasonID int32, info *model.TTPvpPlayerInfo) (err error)
	// 更新多个玩家基本信息
	UpdateTTPvpMutliPlayerInfo(seasonID int32, infos map[int64]*model.TTPvpPlayerInfo) (err error)
	// 获取玩家录像列表
	GetTTPvpRecord(serverID int64, id int64) (info *model.TTPvpReport, err error)
	// 保存玩家录像信息
	SaveTTPvpRecord(serverID int64, id int64, info *model.TTPvpReport) (err error)
	// 获取巅峰对决战报
	GetTTPvpBossRecord(serverID int64) (replayIDs []string, err error)
	// 保存巅峰对决战报
	SaveTTPvpBossRecord(serverID int64, replayID string) (err error)
	// GetTTPvpDailyTaskInfo 查询玩家每日任务进度数据
	GetTTPvpDailyTaskInfo(serverID int64, season int32, seasonDay int32, playerID int64) (info *model.TTPvpPlayerDailyTaskInfo, err error)
	// AddTTPvpDailyWinCount 累加玩家每日挑战次数
	AddTTPvpDailyChallengeCount(serverID int64, season, seasonDay int32, playerID int64) (ok bool, err error)
	// AddTTPvpDailyWinCount 累加玩家每日胜利次数
	AddTTPvpDailyWinCount(serverID int64, season, seasonDay int32, playerID int64) (ok bool, err error)
	// GetTTPvpDailyTaskRewardRecord 查询玩家每日任务已领取奖励记录
	GetTTPvpDailyTaskRewardRecord(serverID int64, season, seasonDay int32, playerID int64) (list []int32, err error)
	// AddTTPvpDailyTaskRewardRecord  更新每日任务领奖记录
	AddTTPvpDailyTaskRewardRecord(serverID int64, season, seasonDay int32, playerID int64, taskID int32) (ok bool, err error)
	// ExistTTPvpDailyTaskRewardRecord 查询玩家是否已领取指定每日任务奖励
	ExistTTPvpDailyTaskRewardRecord(serverID int64, season, seasonDay int32, playerID int64, taskID int32) (ok bool, err error)
	// GetTTPvpLevelTaskRewardRecord 查询玩家段位任务已领取奖励记录
	GetTTPvpLevelTaskRewardRecord(serverID int64, season int32, playerID int64) (list []int32, err error)
	// AddTTPvpDailyTaskRewardRecord  更新段位已领奖记录
	AddTTPvpLevelTaskRewardRecord(serverID int64, season int32, playerID int64, taskID int32) (ok bool, err error)
	// ExistTTPvpLevelTaskRewardRecord 查询玩家是否已领取指定段位任务奖励
	ExistTTPvpLevelTaskRewardRecord(serverID int64, season int32, playerID int64, taskID int32) (ok bool, err error)
	// GetTTPvpLevelTaskRewardStat 查询段位任务累计领取奖励次数
	GetTTPvpLevelTaskRewardStat(serverID int64, season int32) (infos map[int32]int64, err error)
	// GetMulityTTPvpLevelTaskRewardRecord 获取一批玩家的领奖记录
	GetMulityTTPvpLevelTaskRewardRecord(serverID int64, season int32, playerIDs []int64) (recordMap map[int64][]int32, err error)
	// AddTTPvpLevelTaskRewardCount 累加本服段位任务已领奖数量统计
	AddTTPvpLevelTaskRewardCount(serverID int64, season int32, taskID int32) (ok bool, err error)
	// 获取冠绝诸天信息
	GetTTPvpHistoryRank(serverID int64) (infos []*model.TTPvpHistory, err error)
	// 天梯斗法获取段位限制
	GetTTPvpLevelLimit(serverID int64, seasonID int32) map[int32][]int64
	// 保存天梯斗法段位限制
	SaveTTPvpLevelLimit(serverID int64, seasonID int32, limitMap map[int32][]int64) (err error)
	// AddTTPvpRobotToRank 将机器人批量插入排行榜
	AddTTPvpRobotToRank(zoneID int64, seasonID int32, robots map[int64]*model.TTPvpPlayerInfo, robotScore map[int64]int64) (err error)
	// AppendTTPVPRank 增加诸天pvp排行
	AppendTTPVPRank(zoneID int64, seasonID int32, list []*pbrank.RankData) (err error)
	// UpdateTTPvpRank 更新诸天斗法玩家分数
	UpdateTTPvpScore(zoneID int64, seasonID int32, pid1, svrid1 int64, score1 float64, pid2, svrid2 int64, score2 float64) (pid1NewScore, pid2NewScore float64, err error)
	// GetTTPvpRankInfo 查询唯一id对应的排行信息
	GetTTPvpRankInfo(zoneID int64, serverID int64, seasonID int32, uniqueID int64) (info *pbrank.RankData, err error)
	// GetTTPvpZoneRank 获取诸天斗法战区排行榜
	GetTTPvpRankList(zoneID int64, serverID int64, seasonID int32, start, stop int64) (list []*pbrank.RankData, err error)
	// GetTTPvpScoreList 获取诸天斗法积分排行榜
	GetTTPvpScoreList(zoneID int64, serverID int64, seasonID int32, startScore, stopScore int64) (list []*pbrank.RankData, err error)
	// 保存 TTPVP 之前的分数
	SaveTTPvpPreviousScore(zoneID int64, serverID int64, seasonID int32, uniqueID int64, previousScore int64) (err error)
	// 获取 TTPVP 之前的分数
	GetTTPvpPreviousScore(zoneID int64, serverID int64, seasonID int32, uniqueID int64) (previousScore int64, err error)
	// 获取单个玩家基本信息
	GetZTPvpPlayerInfo(seasonID int32, id int64) (info *model.ZTPvpPlayerInfo, err error)
	// 获取多个玩家基本信息
	GetMutliZTPvpPlayerInfo(seasonId int32, ids []int64) (infos map[int64]*model.ZTPvpPlayerInfo, err error)
	// 获取多个玩家基本信息
	UpdateZTPvpPlayerInfo(seasonID int32, info *model.ZTPvpPlayerInfo) (err error)
	// 更新多个玩家基本信息
	UpdateZTPvpMutliPlayerInfo(seasonID int32, infos map[int64]*model.ZTPvpPlayerInfo) (err error)
	// 获取玩家录像列表
	GetZTPvpRecord(serverID int64, id int64) (info *model.ZTPvpReport, err error)
	// 保存玩家录像信息
	SaveZTPvpRecord(serverID int64, id int64, info *model.ZTPvpReport) (err error)
	// 获取巅峰对决战报
	GetZTPvpBossRecord(serverID int64) (replayIDs []string, err error)
	// 保存巅峰对决战报
	SaveZTPvpBossRecord(serverID int64, replayID string) (err error)
	// GetZTPvpDailyTaskInfo 查询玩家每日任务进度数据
	GetZTPvpDailyTaskInfo(serverID int64, season int32, seasonDay int32, playerID int64) (info *model.ZTPvpPlayerDailyTaskInfo, err error)
	// AddZTPvpDailyWinCount 累加玩家每日挑战次数
	AddZTPvpDailyChallengeCount(serverID int64, season, seasonDay int32, playerID int64) (ok bool, err error)
	// AddZTPvpDailyWinCount 累加玩家每日胜利次数
	AddZTPvpDailyWinCount(serverID int64, season, seasonDay int32, playerID int64) (ok bool, err error)
	// GetZTPvpDailyTaskRewardRecord 查询玩家每日任务已领取奖励记录
	GetZTPvpDailyTaskRewardRecord(serverID int64, season, seasonDay int32, playerID int64) (list []int32, err error)
	// AddZTPvpDailyTaskRewardRecord  更新每日任务领奖记录
	AddZTPvpDailyTaskRewardRecord(serverID int64, season, seasonDay int32, playerID int64, taskID int32) (ok bool, err error)
	// ExistZTPvpDailyTaskRewardRecord 查询玩家是否已领取指定每日任务奖励
	ExistZTPvpDailyTaskRewardRecord(serverID int64, season, seasonDay int32, playerID int64, taskID int32) (ok bool, err error)
	// GetZTPvpLevelTaskRewardRecord 查询玩家段位任务已领取奖励记录
	GetZTPvpLevelTaskRewardRecord(serverID int64, season int32, playerID int64) (list []int32, err error)
	// GetMulityZTPvpLevelTaskRewardRecord 获取一批玩家的领奖记录
	GetMulityZTPvpLevelTaskRewardRecord(serverID int64, season int32, playerIDs []int64) (recordMap map[int64][]int32, err error)
	// AddZTPvpDailyTaskRewardRecord  更新段位已领奖记录
	AddZTPvpLevelTaskRewardRecord(serverID int64, season int32, playerID int64, taskID int32) (ok bool, err error)
	// ExistZTPvpLevelTaskRewardRecord 查询玩家是否已领取指定段位任务奖励
	ExistZTPvpLevelTaskRewardRecord(serverID int64, season int32, playerID int64, taskID int32) (ok bool, err error)
	// GetZTPvpLevelTaskRewardStat 查询段位任务累计领取奖励次数
	GetZTPvpLevelTaskRewardStat(serverID int64, season int32) (infos map[int32]int64, err error)
	// AddZTPvpLevelTaskRewardStat 累加本服段位任务已领奖数量统计
	AddZTPvpLevelTaskRewardStat(serverID int64, season int32, taskID int32) (ok bool, err error)
	// 天梯诸天获取段位限制
	GetZTPvpLevelLimit(serverID int64, seasonID int32) map[int32][]int64
	// 保存诸天斗法段位限制
	SaveZTPvpLevelLimit(serverID int64, seasonID int32, limitMap map[int32][]int64) (err error)
	// AppendZTPVPRank 增加诸天pvp排行
	AppendZTPVPRank(zoneID int64, seasonID int32, list []*pbrank.RankData) (err error)
	// UpdateZTPvpRank 更新诸天斗法玩家分数
	UpdateZTPvpScore(zoneID int64, seasonID int32, pid1, svrid1 int64, score1 float64, pid2, svrid2 int64, score2 float64) (pid1NewScore, pid2NewScore float64, err error)
	// GetZTPvpRankInfo 查询唯一id对应的排行信息
	GetZTPvpRankInfo(zoneID int64, serverID int64, seasonID int32, uniqueID int64) (info *pbrank.RankData, err error)
	// GetZTPvpScoreList 获取诸天斗法积分排行榜
	GetZTPvpScoreList(zoneID int64, serverID int64, seasonID int32, startScore, stopScore int64) (list []*pbrank.RankData, err error)
	// GetZTPvpZoneRank 获取诸天斗法战区排行榜
	GetZTPvpRankList(zoneID int64, serverID int64, seasonID int32, start, stop int64) (list []*pbrank.RankData, err error)
	// AddZTPvpRobotToRank 将机器人批量插入排行榜
	AddZTPvpRobotToRank(zoneID int64, seasonID int32, robots map[int64]*model.ZTPvpPlayerInfo, robotScore map[int64]int64) (err error)
	// SaveZTPvpPreviousScore 保存 ZTPvp 之前的分数
	SaveZTPvpPreviousScore(zoneID int64, serverID int64, seasonID int32, uniqueID int64, previousScore int64) (err error)
	// GetZTPvPPreviousScore 获取 ZTPvp 之前的分数
	GetZTPvpPreviousScore(zoneID int64, serverID int64, seasonID int32, uniqueID int64) (previousScore int64, err error)

	//////////////////////////Friend相关接口
	// 好友跨天重置
	FriendReset(now int64)
	// 好友点信息存储
	CacheFriendPointInfo(roleID int64, friendInfo *model.FriendPoint) (err error)
	// 好友点信息读取
	GetFriendPointInfo(roleID int64) (friendInfo *model.FriendPoint, err error)
	// 存储好友申请id列表
	CacheRequestIDs(roleID int64, requestIDs []int64) (err error)
	// 是否是好友申请
	IsRequest(roleID, requestID int64) (isRequest bool, err error)
	// 获取所有好友申请id
	GetAllRequestID(roleID int64) (ids []int64, err error)
	// 删除好友申请
	DelRequest(roleID, id int64) (err error)
	// 删除好友申请列表
	DelRequests(roleID int64, ids []int64) (err error)
	// 存储好友申请详情
	CacheRequest(roleID int64, request *model.FriendRequest) (err error)
	// 获取所有好友申请详情
	GetAllRequest(roleId int64, requestIds []int64) (request []*model.FriendRequest, err error)
	// 存储好友id
	CacheFriendID(roleID, friendID int64) (err error)
	// 是否是好友
	IsFriend(roleID, friendID int64) (isFriend bool, err error)
	// 获取所有好友id
	GetAllFriendID(roleID int64) (ids []int64, err error)
	// 存储好友信息
	CacheFriend(roleID int64, info *model.FriendInfo) (err error)
	// 存储好友信息列表
	CacheFriends(roleID int64, infos map[int64]*model.FriendInfo) (err error)
	// 获取好友信息
	GetFriendInfo(roleID, friendID int64) (info *model.FriendInfo, err error)
	// 获取好友信息列表
	GetFriendInfos(roleID int64, ids []int64) (infos map[int64]*model.FriendInfo, err error)
	// 删除好友
	DelFriend(roleID, friendID int64) (err error)
	// 已经赠送或者领取过好友点的好友记录
	CacheFriendPoint(roleID, friendID int64) (err error)
	// 存储黑名单
	CacheBlackListID(roleID, blackListID int64) (err error)
	// 是否是黑名单玩家
	IsBlackList(roleID, blackListID int64) (isBlackList bool, err error)
	// 获取所有黑名单玩家
	GetAllBlackListID(roleID int64) (ids []int64, err error)
	// 移除黑名单
	DelBlackList(roleID, blackListID int64) (err error)
	// 是否是今天删除过的好友
	IsDel(roleId, friendId int64) (isDel bool, err error)
	// 存储已经推荐的好友
	CacheRecommend(roldId, friendId int64) (err error)
	// 根据等级获取一批玩家id
	GetPlayerByLevels(uid int64, levels []int32) (roleIds []int64, err error)
	// CacheLeaseHero 缓存角色租借英雄
	CacheLeaseHero(id int64, leaseHero *model.LeaseHero) (err error)
	// 更新战斗次数
	UpdateLeaseFightCount(id int64, fightCount map[int32]int32) (err error)
	CacheLeaseRequest(uid int64, request *model.LeaseRequest) (err error)
	UpdateLeaseRequest(request *model.LeaseRequest) (err error)
	GetLeaseRequest(reqIds []int64) (requestList []*model.LeaseRequest, err error)
	GetAllLeaseRequest(uid int64) (requestMap map[int64]*model.LeaseRequest, err error)
	GetLeaseHero(uid int64) (heros []*model.LeaseHero, err error)
	UpdateLeaseHero(id int64, hero *model.LeaseHero) (err error)
	GetLeaseHeroList(ids []int64) (heros []*model.LeaseHero, err error)
	GetLeaseHeroListByHeroID(ids []int64, heroID int32) (heros []*model.LeaseHero, err error)
	GetSelfLeaseHeroList(id int64) (heros []*model.LeaseHero, err error)
	CacheSelfLeaseHero(id int64, hero *model.LeaseHero) (err error)
	RemoveSelfHero(id int64) (err error)
	DelSelfLeasehero(id int64, heroId int32) (err error)
	DelLeaseRequestList(id int64, reqIds []int64) (err error)
	GetLeaseHeroTask(id int64, taskIds []int32) (leaseTask []*model.LeaseTask, err error)
	UpdateTask(leaseType int32, uid int64, leaseHero *model.LeaseHero, task map[int32][]int32)
	// 保存任务
	SaveTask(uid int64, task *model.LeaseTask) (err error)
	// 获取战斗次数
	GetFightCount(id int64) (fightCount map[int32]int32, err error)
	// 更新战斗次数
	UpdateFightCount(id int64, fightCount map[int32]int32) (err error)
	// 重置数据
	LeaseHeroReset(now int64)
	// 重置租借数据
	ResetLeaseHero(uid int64) (err error)
	UpdateHistoryLease(uid int64, count int32) (err error)
	GetHistoryLease(uid int64) (count int32, err error)

	//////////////////Guild相关接口
	HasGuildName(name string) bool
	// IsExistGuild 是否存在公会
	IsExistGuild(guildid int64) bool
	//取得redis中的公会信息
	GetGuildInfo(guildid int64) *model.GuildInfo
	UpdateGuildInfo(info *model.GuildInfo) (err error)
	// 取得redis中的个人关于公会的信息
	GetGuildUser(userid int64) (info *model.GuildUser)
	UpdateGuildUser(info *model.GuildUser) (err error)
	ListIDs(guildId int64, nameSpace string) (ids []int64)
	HasGuildHighHuntTI(guildid int64, ti int64) bool
	GetRedPacketNum(userid int64) int32
	GetRedPacket(uuid int64) (info *model.CacheRedPacket)
	GetAuctionGoods(uuid int64) (info *model.CacheAuctionGoods)
	GetAuctionBider(uuid int64) (info *model.CacheAuctionBider)
	GetAuctionBidback(userid int64) (int64, error)
	AuctionLastBiddingFailedTime(userid int64) int64
	UpdateAuctionBidback(userid int64, val int64) error
	GetAuctionLastUUID(userid int64) int64
	GetDCAgentInfo(userid int64) (info *model.CacheDCAgentInfo)
	UpdateDCAgentPreRank(userid int64, rank int32) error
	GetDCAgentPreRank(userid int64) int32
	IsInDungeon(activityid int32, userid int64) bool

	/////////////////////Login相关接口
	SetAccountToken(ctx context.Context, accountId string, platformId, channelId int) (accountToken string, err error)
	VerifyAccountToken(ctx context.Context, accountId string, accountToken string) (platformId, channelId int64, err error)
	AllocPlayerID() int64
	AddUserID(sdkuuid string, userid int64) (err error)
	AddLoginDeviceInfo(playerId int64, data map[string]interface{}) error
	GetLoginInfo(playerId int64) (loginInfo *coremodel.TLogFields, err error)
	// 登录生成token
	// token =  base64({userid:1001}).md5(userid+ts+basestr)
	// 同一个accountId 每次都生成的是相同的token
	GetToken(ctx context.Context, playerId int64) (accessToken, refreshToken string, err error)
	//刷新Token
	RefreshToken(ctx context.Context, refreshToken string) (accessToken, refreshTokenNew string, err error)

	///////////////////////Push相关接口
	// PushScene 推送通知：推送给某个服务器ID下的所有scene分线. 例如: dao.PushScene(1, pbscene.PushMessageID_PushMailNotice, &pbactivity.Empty{})
	PushScene(serverid int32, pushid int32, msg proto.Message)
	// 推送场景数据
	PushSceneData(serverid int32, pushid int32, data []byte)
	// 推送通知：ids为空表示推送本服通知，否则推送给具体的指定玩家. 例如：dao.PushUser([]int64{190860, 127648}, pbscene.PushMessageID_PushMailNotice, &pbactivity.Empty{})
	PushUser(ctx context.Context, ids []int64, pushid int32, msg proto.Message)
	// 推送用户数据
	PushUserData(ctx context.Context, ids []int64, pushid int32, data []byte)

	//////////////////////Scene相关接口
	GetCacheLineup(userid int64, group int32) (list []*battle.LineupInfo, powers []int64)
	GetCacheCampParam(userid int64, group int32) (list []*battle.CampParam)
	CacheRoleLineup(id int64, m pbscene.GroupType, lineup []*battle.LineupInfo, battle []*battle.CampParam) (err error)
	IsExistRole(userid int64) bool
	GetCacheRole(userid int64) (info *model.CacheRole)
	SetCacheRoleOs(userid int64, os int32) error
	GetCacheRoleOs(userid int64) (int32, error)
	SetCacheRoleUnionId(userId int64, unionId string) error
	GetCacheRoleUnionId(userId int64) (string, error)
	GetCacheRoleLevel(userid int64) int
	// 获取多个玩家基本信息
	GetMutliPlayerInfo(ids []int64) (infos map[int64]*model.CacheRole)
	// CacheRoleInfo 缓存角色数据
	CacheRoleInfo(info *model.CacheRole) (err error)
	CacheRoleServer(playerId, server int64) (err error)
	SaveRegTime(userId, time int64) error
	// CacheMutliRoleInfo 存储一批玩家信息(给机器人用)
	CacheMutliRoleInfo(infos []*model.CacheRole) (err error)
	UpdateLevel(oldLevel, newLevel int32, roleID int64) (err error)

	// update names pool TODO: might need rename to proper name to avoid confusing, after robot gen logic change merged
	UpdateName(name string, serverid int32)
	// delete name from names pool// TODO: might need rename to proper name to avoid confusing, after robot gen logic change merged
	DeleteName(name string, serverid int32)

	AddRoleID(serverid int32, userid int64)
	IsHaveName(name string, serverid int32) bool
	GetScenePrime() (info *model.CacheScenePrime)
	UpdateScenePrime(info *model.CacheScenePrime)
	UpdateUserHeartBeat(playerID uint64)
	DelAccountUserID(sdkuuid string, userid int64) (err error)
	UpdatePlayerLine(userid int64, appid string)
	DeletePlayerLine(userid int64)
	PlayerLine(userid int64) (string, error)
	PlayerLineId(userid int64) int
	SetUserDisable(userid int64, disableTime int64, tips string) (err error)
	GetUserDisable(userid int64) int64
	SetForbidChat(userid int64, channel int32, unixtm int64, tips string) (err error)
	GetForbidChat(userid int64, channel int32) (unixtm int64, tips string)
	// 同步微服务所需的系统解锁信息, id:SystemUnlock.xlsx的解锁ID
	CacheRoleSystemUnlock(userid int64, id int32)
	GetSystemScore(userid int64, t pbscene.CompareType) int64
	SetSystemScore(userid int64, t pbscene.CompareType, val int64)
	GetCacheEW(userid int64, ewid int32) (info *model.CacheExclusiveWeapon)
	UpdateCacheEW(userid int64, info *model.CacheExclusiveWeapon) bool
	GetDHClock(userid int64) int64
	SetDHClock(userid int64, val int64)
	GetCacheLingbao(userid int64, id int32) (info *model.CacheLingbao)
	UpdateCacheLingbao(userid int64, info *model.CacheLingbao)
	GetCachePenglai(userid int64) (info *model.CachePenglai)
	UpdateCachePenglai(userid int64, info *model.CachePenglai)
	GetCacheBannerInfo(userid int64) (info *model.CacheBannerInfoAll)
	UpdateCacheBannerInfo(userid int64, bannerIndex int32, info *model.CacheBannerInfo)
	GetCacheConstellation(userid int64) (info *model.CacheConstellation)
	UpdateCacheConstellation(userid int64, info *model.CacheConstellation)
	VerifyToken(playerID int64, stoken string) (err error)

	///////////////////////State相关接口
	HeartBeat(ctx context.Context, req *pb.HeartBeatReq) (resp *pb.CommonResp, err error)
	GetServiceList2(ctx context.Context, name string) (srvInfos []*pb.StatefulServiceInfo, err error)
	DeleteServiceState2(ctx context.Context, name string, appID string) (err error)
	UpdateServicePlayerNum(ctx context.Context, namespace, serviceName, appId string, change int64) error
	GetScenePlayerNum(ctx context.Context, namespace string) (num int, err error)
	ResetScenePlayerNum(ctx context.Context, namespace string) (err error)
	LoadBalance(stslist []*pb.StatefulServiceInfo) (target *pb.StatefulServiceInfo)
	PlayerOnline(ctx context.Context, req *pb.PlayerOnlineReq) (resp *pb.CommonResp, err error)
	PlayerOffline(ctx context.Context, req *pb.PlayerOfflineReq) (resp *pb.CommonResp, err error)
	PlayerOfflineHandler(ctx context.Context, playerID int64) (err error)
	GetPlayerState(ctx context.Context, playerId int64) (state int64, err error)
	GetPlayerAllServiceAddress(ctx context.Context, playerID int64) (addrMap map[string]string, err error)
	UpdatePlayerState(ctx context.Context, playerID int64, addressMap map[string]string, state int64) (err error)

	//////////////////////////////多个模块有耦合相关接口
	GetRoleInfos(uids []int64) map[int64]*pbworld.RoleExInfo
	GetRoleInfo(uid int64) *pbworld.RoleExInfo
	GetLineupFromOnlineOrRedis(userid int64, group int32) (lineups []*battle.LineupInfo, powers []int64, camps []*battle.CampParam)
	GetMainLineupFromOnlineOrRedis(userid int64) (lineups []*battle.LineupInfo, powers []int64, camps []*battle.CampParam)
	QueryPlayer(param string) (playerIDs []int64, err error)
	// 模糊查询存在的玩家名称 可能不准 使用redis scan 实现
	MatchPlayerNames(name string) []string
	// 模糊查询玩家名称 对应的玩家id 可能不准 使用redis scan 实现
	MatchPlayerNameIds(name string) []int64
	// 玩家id到账号id查询 redis
	Role2Account(userId []int64) map[int64]string

	// 处理在线玩家数量
	AddNamespaceSet()
	DelNamespaceSet()
	ClearNamespacePlayerLineNum()
	BatchGetScenePlayerNum(context.Context) (serverPlayerNumList []*pbaccount.ServerPlayerNum, err error)
	BatchGetScenePlayerNumV1(_ context.Context) ([]*pbaccount.ServerPlayerNum, error)
	DelNamespacePlayerLine() error // 去掉分线记录

	// add a new func to handle the rename, remove old name form name pool, update name to name pool, handle player_name:xxx key,
	// remove old key and add new key(add new item(serverid) to set),
	// remove old name from names pool and add new name to names pool
	// if is robot, skip some steps, eg: updatename(pool) & player_name:xx key related changes
	ChangeName(playerID int64, oldName, newName string) error

	// GetRes 接口频率限制
	// 这是一个公共方法，所以谨慎命名uniqueKey
	GetRes(uniqueKey string, expire ...time.Duration) (bool, error)
}
