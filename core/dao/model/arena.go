package model

import (
	pbarena "xy3-proto/arena"
	battle "xy3-proto/battle"
)

////////////////////////////////////////////////////////////////////////////////////
//机器人

// // Robot 机器人结构体
// type Robot struct {
// 	PlayerInfo  *CacheRole         // 机器人基础信息
// 	LineupInfo1 *battle.LineupInfo // 展示阵容1
// 	LineupInfo2 *battle.LineupInfo // 展示阵容2
// 	LineupInfo3 *battle.LineupInfo // 展示阵容3
// 	CampParam1  *battle.CampParam  // 战斗阵容1信息
// 	CampParam2  *battle.CampParam  // 战斗阵容2信息
// 	CampParam3  *battle.CampParam  // 战斗阵容3信息
// 	PowerID     int32              // 战力id
// }

// type RobotLineupCache struct {
// 	RobotMap map[int64]*Robot // 机器人阵容信息
// }

////////////////////////////////////////////////////////////////////////////////////
//斗法

// PvpFightingReport 战报列表
type PvpFightingReport struct {
	FightReport []*pbarena.FightReport // 战报ID列表
}

// PvpFightingRobot 机器人斗法信息
type PvpFightingRobot struct {
	RoleID   int64 // 角色id
	Rank     int32 // 排行
	PvpLevel int32 // 段位
}

////////////////////////////////////////////////////////////////////////////////////
//天梯斗法

// TTPvpPlayerInfo 天梯斗法玩家数据
type TTPvpPlayerInfo struct {
	ID             int64           // 玩家id
	PvpLevel       int32           // 段位
	ServerID       int64           // 所属serverid
	ChallengeCount int32           // 本赛季累计挑战次数
	WinCount       int32           // 本赛季累计胜利次数
	BestPvpLevel   int32           // 本赛季历史最高段位
	BestScore      int64           // 本赛季历史最高积分
	WinningStreak  int32           // 连胜场次
	WhiteList      map[int64]int32 // 白名单
}

// TTPvpPlayerDailyTaskInfo 天梯斗法玩家每日数据
type TTPvpPlayerDailyTaskInfo struct {
	ID             int64   // 玩家id
	ServerID       int64   // 所属serverid
	ChallengeCount int32   // 每日挑战次数
	WinCount       int32   // 每日胜利次数
	AwardRecord    []int32 // 已领取任务id列表
}

// TTPvpPlayerLevelTaskInfo 天梯斗法玩家历史段位数据
type TTPvpPlayerLevelTaskInfo struct {
	ID          int64   // 玩家id
	ServerID    int64   // 所属serverid
	AwardRecord []int32 // 已领取任务id列表
	TaskLimit   []int32 // 已达到上限的任务
}

// TTPvpReport 战报列表
type TTPvpReport struct {
	RoleID      int64                  // 角色id
	FightReport []*pbarena.FightReport // 战报列表
}

// 天梯斗法冠绝诸天排行信息
type TTPvpHistoryRank struct {
	RoleID     int64              // 角色id
	Rank       int64              // 排行
	Score      int64              // 积分
	PvpLevel   int32              // 段位
	LineupInfo *battle.LineupInfo // 阵容信息
}

// 天梯斗法冠绝诸天信息
type TTPvpHistory struct {
	SeasonID    int32               // 赛季id
	StartTime   int64               // 开始时间
	EndTime     int64               // 结束时间
	HistoryRank []*TTPvpHistoryRank // 排行列表
}

// //////////////////////////////////////////////////////////////////////////////////
// 诸天斗法
// ZTPvpPlayerInfo 诸天斗法玩家数据
type ZTPvpPlayerInfo struct {
	ID             int64           // 玩家id
	PvpLevel       int32           // 段位
	ServerID       int64           // 所属serverid
	ChallengeCount int32           // 本赛季累计挑战次数
	WinCount       int32           // 本赛季累计胜利次数
	BestPvpLevel   int32           // 本赛季历史最高段位
	BestScore      int64           // 本赛季历史最高积分
	WinningStreak  int32           // 连胜场次
	WhiteList      map[int64]int32 // 白名单
}

// ZTPvpPlayerDailyTaskInfo 诸天斗法玩家每日数据
type ZTPvpPlayerDailyTaskInfo struct {
	ID             int64   // 玩家id
	ServerID       int64   // 所属serverid
	ChallengeCount int32   // 每日挑战次数
	WinCount       int32   // 每日胜利次数
	AwardRecord    []int32 // 已领取任务id列表
}

// ZTPvpPlayerLevelTaskInfo 诸天斗法玩家历史段位数据
type ZTPvpPlayerLevelTaskInfo struct {
	ID          int64   // 玩家id
	ServerID    int64   // 所属serverid
	AwardRecord []int32 // 已领取任务id列表
	TaskLimit   []int32 // 已达到上限的任务
}

// ZTPvpReport 战报列表
type ZTPvpReport struct {
	RoleID      int64                  // 角色id
	FightReport []*pbarena.FightReport // 战报列表
}
