package model

const (
	RedisScenePrime            string = "%v:sceneprime" //本服通用数据
	RedisKey_Player                   = "player:%d"     // 角色信息
	RedisKey_Player_Level             = "player_level:%d"
	RedisKey_Player_Name              = "player_name:%v"
	RedisKey_Player_Line              = "player_line:%v"                         //玩家所在分线，比如：scene-0， 过期时间8min，玩家每5min写入一次
	RedisKey_Player_Disable           = "player_disable:%v"                      // 玩家封号列表
	RedisKey_Player_ForbidChat        = "player_forbidchat:%v"                   // 玩家禁言
	RedisNameTable             string = "names:%v"                               //名字集合：stdbase64(name)->{服ID集合}
	RedisRoleTable             string = "roles:%v"                               //真实角色ID集合：服ID->{playerid}
	RedisUserStateKey                 = "coordinator:player:heartbeat:%d:string" //  用户的心跳时间
	RedisKey_Account                  = "account:%v"                             //account:unionid set:<playerid>
	RedisUserServerAlloc       string = "coordinator:player:server:%d:hash"      //<scene, line>
)

const (
	// 系统解锁
	RedisKey_SystemUnlock = "systemunlock:%v"
)

const (
	// 展示阵容
	RedisKey_Lineup = "lineup:%v:%v"
	// 战斗阵容
	RedisKey_Battle = "battle:%v:%v"
)

const (
	RedisKey_SysScore        = "sysscore:%v:%v"   //系统评分(用于战力对比) <playerid:sysid, score>
	RedisKey_TouhouClock     = "dhclock:%v"       //东皇钟 <playerid, order>
	RedisKey_LingBao         = "lingbao:%v:%v"    //灵宝 <playerid:id, lingbao{}>
	RedisKey_PengLai         = "penglai:%v"       //蓬莱 <playerid, jsonobj>
	RedisKey_BannerInUse     = "bannerInUse:%v"   //战旗 <playerid, jsonobj>
	RedisKey_Constellation   = "constellation:%v" //星宿 <playerid, jsonobj>
	RedisKey_ExclusiveWeapon = "ew:%v:%v"         //专属 <playerid:ewid, ew{}>
)

const (
	RedisGuildNames               string = "%v:guildnames" //公会名字(Set)
	RedisGuildUserInfo            string = "%v:guildmem:%v:info"
	RedisGuildMembers             string = "%v:guild:%v:members"
	RedisGuildInfo                string = "%v:guild:%v:info"
	RedisGuildHunts               string = "%v:guild:%v:hunts"        //公会开启的高阶狩猎boss(开启时间戳<<16+index)(Set)
	RedisMemberRedPackets         string = "%v:guildmem:%v:ungive"    //成员的未发的红包集合
	RedisRedPacketInfo            string = "%v:redpacket:%v"          //红包池
	RedisAuctionGoods             string = "%v:auction:goods:%v:%v"   //当天拍卖品
	RedisAuctionBider             string = "%v:auction:bider:%v:%v"   //当天竞拍投标
	RedisAuctionBidback           string = "%v:auction:bidback:%v"    //竞拍失败被退回的资金
	RedisAuctionPlayerLastUUID    string = "%v:auction:lastuuid:%v"   //玩家最近一次竞价成功拍卖品UUID, 被别人超越不用覆盖
	RedisDevilConquerAgentInfo    string = "%v:devil:agent:%v"        //进入地牢ID的玩家数据 Hash
	RedisDungeonPlayers           string = "%v:dungeon:%v"            //进入地牢ID的玩家Set
	RedisDevilConquerAgentPreRank string = "%v:devil:prerank:%v"      //仙盟诛邪前次排名，用于成就检查
	RedisAuctionBiddingFailed     string = "%v:auction:failedtime:%v" //玩家最近一次竞拍失败的时间戳
)

/////////////////////////////////////////////////////////////////////////////////////////////////////
// 好友

const (
	RedisKey_Friend_Request          = "friend:friend:request:%v"             // 好友请求:角色id
	RedisKey_Friend_Request_Hash     = "friend:friend:request_hash:%v:%v"     // 好友请求hash:角色id:好友id
	RedisKey_Friend_Friend           = "friend:friend:friend:%v"              // 好友列表:角色id
	RedisKey_Friend_Friend_Hash      = "friend:friend:friend_hash:%v:%v"      // 好友列表hash:角色id:好友id
	RedisKey_Friend_CrossFriend      = "friend:friend:crossfriend:%v"         // 跨服好友:角色id
	RedisKey_friend_CrossFriend_Hash = "friend:friend:crossfriend_hash:%v:%v" // 好友列表hash:角色id:好友id
	RedisKey_Friend_BlackList        = "friend:friend:blacklist:%v"           // 黑名单:角色id
	RedisKey_Friend_Point            = "friend:friend:point"                  // 今天赠送了好友点的玩家
	RedisKey_Friend_Point_Hash       = "friend:friend:point_hash:%v"          // 好友点赠送hash:角色id
	RedisKey_Friend_Del              = "friend:friend:del:%v"                 // 已经赠送或者领取过好友点的好友记录
	RedisKey_Friend_Recommend        = "friend:friend:recommend:%v"           // 今天已经推荐过的好友
)

// 租借

const (
	// 可供租借的英雄
	RedisKey_LeaseHero_Hero           = "friend:leasehero:hero:%v"             // 可出借英雄列表:角色id
	RedisKey_LeaseHero_Hero_Hash      = "friend:leasehero:hero:hash:%v:%v"     // 可出借英雄信息:角色id:英雄id
	RedisKey_LeaseHero_FightCount     = "friend:leasehero:fightcount:%v"       // 租借英雄战斗次数:角色id
	RedisKey_LeaseHero_Request        = "friend:leasehero:request:%v"          // 租借申请列表:角色id
	RedisKey_LeaseHero_Request_Hash   = "friend:leasehero:request:hash:%v"     // 租借申请列表详情:角色id
	RedisKey_Self_LeaseHero_Hero      = "friend:leasehero:selfhero:%v"         // 已经借入英雄:角色id
	RedisKey_Self_LeaseHero_Hero_Hash = "friend:leasehero:selfhero:hash:%v:%v" // 已经借入英雄详情:角色id
	RedisKey_LeaseHero_Task           = "friend:leasehero:task:%v:%v"          // 租借英雄任务:角色id:任务id
	RedisKey_LeaseHero_HistoryLease   = "friend:leasehero:historylease:%v"     // 租借历史记录:角色id

	// 借出x次
	TaskType1 = 1
	// 借出x次x觉醒
	TaskType2 = 2
	// 借入x次
	TaskType3 = 3
	// 借入x次x觉醒
	TaskType4 = 4
)

/////////////////////////////////////////////////////////////////////////////////////////////////////
// 竞技场

const (
	RefreshEveryDay = "arena:timer:refresh_flag:%v"
)

// const (
// 	RedisKey_Robot_Init  = "arena:robot:init:%v"   // 机器人是否已经初始化标识
// 	RedisKey_Robot       = "arena:robot:robots:%v" // 机器人key:id
// 	RedisKey_Robot_Power = "arena:robot:power:%v"  // 机器人战力分段:战力id
// )

const (
	PvpFighting                     int32 = 1                              //斗法
	TTPvp                           int32 = 2                              //天梯斗法
	ZTPvp                           int32 = 3                              // 诸天斗法
	RedisKey_Arena_Daily_Rank             = "arena:daily:rank:%v:%v:%v"    // 斗法/诸天斗法/天梯斗法 每日任务排行记录:服务器id:斗法类型:时间
	RedisKey_Arena_Achievement_Rank       = "arena:achievement:rank:%v:%v" // 斗法/诸天斗法/天梯斗法 成就排行记录:服务器id:斗法类型
)

const (
	RedisKey_PvpFighting_Record = "arena:pvpfighting:record:%v:%v" // 斗法战报:服务器id:角色id
	RedisKey_PvpFighting_Init   = "arena:pvpfighting:init:%v"      // 斗法是否初始化:服务器id
	RedisKey_PvpFighting_Time   = "arena:pvpfighting:time:%v"      // 斗法第一次有玩家进入时间:服务器id
	RedisKey_PvpFighting_Lock   = "arena:pvpfighting:lock:%v:%v"   // 斗法排行交换锁:服务器id:段位

	RedisKey_PvpFighting_Rankings = "arena:ranklist:pvpfighting:rankings:%v:%v" // 排名映射到玩家ID
)

const (
	RedisKey_TTPvp_Player = "arena:ttpvp:player:%v:%v" // 天梯斗法玩家信息:赛季:玩家id

	RedisKey_TTPvp_DailyTask        = "arena:ttpvp:dailytask:%v:%v:%v:%v"       // 天梯斗法每日任务:服务器id:赛季:天数:玩家id
	RedisKey_TTPvp_DailyTask_Reward = "arena:ttpvp:dailytaskreward:%v:%v:%v:%v" // 天梯斗法每日任务领取情况:服务器id:赛季:天数:玩家id
	RedisKey_TTPvp_LevelTask_Reward = "arena:ttpvp:leveltask:%v:%v:%v"          // 天梯斗法段位任务:服务器id:赛季:玩家id
	RedisKey_TTPvp_LevelTask_Count  = "arena:ttpvp:leveltask:%v:%v"             // 天梯斗法段位任务领取计数:服务器id:赛季id

	RedisKey_TTPvp_Zone_Rank = "arena:ranklist:ttpvp:zone:%d_%d"  // 天梯斗法战区排行:战区id:赛季id
	RedisKey_TTPvp_Rank      = "arena:ranklist:ttpvp:local:%d_%d" // 天梯斗法服务器排行:服务器id:赛季id

	RedisKey_TTPvp_Zone_PreviousScore = "arena:ttpvp:previous_score:zone:zone_%d:season_%d"    // 天梯斗法战区之前的分数:战区id:赛季id
	RedisKey_TTPvp_PreviousScore      = "arena:ttpvp:previous_score:local:server_%d:season_%d" // 天梯斗法上一次分数:服务器id:赛季id

	RedisKey_TTPvp_Record       = "arena:ttpvp:record:%v:%v"     // 天梯斗法战报:服务器id:玩家id
	RedisKey_TTPvp_Boss_Record  = "arena:ttpvp:boss_record:%v"   // 天梯斗法巅峰对决战报:服务器id
	RedisKey_TTPvp_History_Rank = "arena:ttpvp:history_rank:%v"  // 天梯斗法冠绝诸天:服务器id
	RedisKey_TTPvp_Init         = "arena:ttpvp:Init:%v"          // 天梯斗法是否初始化标识:服务器id
	RedisKey_TTPvp_Level_Limit  = "arena:ttpvp:levellimit:%v:%v" // 天梯斗法段位限制情况:服务器id:赛季id
)

const (
	RedisKey_ZTPvp_Player = "arena:ztpvp:player:%v:%v" // 诸天斗法玩家信息:赛季:玩家id

	RedisKey_ZTPvp_DailyTask        = "arena:ztpvp:dailytask:%v:%v:%v:%v"       // 诸天斗法每日任务:服务器id:赛季:天数:玩家id
	RedisKey_ZTPvp_DailyTask_Reward = "arena:ztpvp:dailytaskreward:%v:%v:%v:%v" // 诸天斗法每日任务领取情况:服务器id:赛季:天数:玩家id
	RedisKey_ZTPvp_LevelTask_Reward = "arena:ztpvp:leveltask:%v:%v:%v"          // 诸天斗法段位任务:服务器id:赛季:玩家id
	RedisKey_ZTPvp_LevelTask_Count  = "arena:ztpvp:leveltask:%v:%v"             // 诸天斗法段位任务领取计数:服务器id:赛季id

	RedisKey_ZTPvp_Zone_Rank = "arena:ranklist:ztpvp:zone:%d_%d"  // 诸天斗法战区排行:战区id:赛季id
	RedisKey_ZTPvp_Rank      = "arena:ranklist:ztpvp:local:%d_%d" // 诸天斗法服务器排行:服务器id:赛季id

	RedisKey_ZTPvp_Zone_PreviousScore = "arena:ztpvp:previous_score:zone:zone_%d:season_%d"    // 诸天斗法战区之前的分数:战区id:赛季id
	RedisKey_ZTPvp_PreviousScore      = "arena:ztpvp:previous_score:local:server_%d:season_%d" // 诸天斗法上一次分数:服务器id:赛季id

	RedisKey_ZTPvp_Record      = "arena:ztpvp:record:%v:%v"     // 诸天斗法战报:服务器id:玩家id
	RedisKey_ZTPvp_Boss_Record = "arena:ztpvp:boss_record:%v"   // 诸天斗法巅峰对决战报:服务器id
	RedisKey_ZTPvp_Init        = "arena:ztpvp:Init:%v"          // 诸天斗法是否初始化标识:服务器id
	RedisKey_ZTPvp_Level_Limit = "arena:ztpvp:levellimit:%v:%v" // 诸天斗法段位限制情况:服务器id:赛季id
)
