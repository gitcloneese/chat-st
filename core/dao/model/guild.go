package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	ImpeachNone  int32 = 0
	ImpeachCan   int32 = 1
	ImpeachBeing int32 = 2
)

const (
	JobNormal   int32 = 0
	JobElite    int32 = 1
	JobViceman  int32 = 2
	JobChairman int32 = 3
)

// GuildInfo 场景服读取公会服写入redis中的公会信息
type GuildInfo struct {
	ID            int64  //公会唯一ID
	Name          string //公会名
	Icon          int32  //会徽ID
	HallLv        int32  //大殿等级
	TowerLv       int32  //炼宝阁等级
	PavilionLv    int32  //修身阁等级
	ShopLv        int32  //商店等级
	Fund          int32  //资金
	FightLimit    int64  //入会战力限制
	Notice        string //公告
	Chairman      int64  //会长
	CoolName      int64  //改名冷却
	CoolNotice    int64  //改公告冷却
	CoolMail      int64  //发邮件冷却
	CoolEnlist    int64  //世界频道招募冷却
	ImpeachST     int64  //弹劾开始时间
	ImpeachState  int32  //弹劾状态
	LiveDay       int32  //仙盟日活跃
	LiveWeek      int32  //仙盟周活跃
	MSoul         int32  //怪物精魄
	LastFlushTm   int64  //上次刷新的时间，用于公会某些数据的刷新，采取延迟刷新机制
	AutoAgreeJoin int32  //自动同意入会申请，1-自动同意
	TBoxNum       int32  //公会专属宝箱数量
	Inactive      bool   //如果所有成员都处于不活跃状态，公会将被设置为不活跃
}

func (p *GuildInfo) IsCoolName() bool {
	return p.CoolName <= time.Now().Unix()
}

func (p *GuildInfo) IsCoolNotice() bool {
	return p.CoolNotice <= time.Now().Unix()
}

func (p *GuildInfo) IsCoolMail() bool {
	return p.CoolMail <= time.Now().Unix()
}

func (p *GuildInfo) IsCoolEnlist() bool {
	return p.CoolEnlist <= time.Now().Unix()
}

func (p *GuildInfo) IsImpeach() bool {
	return (p.ImpeachState == ImpeachBeing) || (p.ImpeachState == ImpeachCan)
}

// GuildUser 场景服读取公会服写入redis中的个人公会信息
type GuildUser struct {
	ID                  int64 //角色ID
	GuildID             int64 //所在公会ID，0-表示无公会
	Job                 int32 //职位
	WeekLive            int32 //周活跃度
	DayLive             int32 //日活跃度
	TotalLive           int32 //总活跃度
	TaskCnt             int32 //任务总接取次数
	TaskNum             int32 //任务已接取次数(+完成的)
	CoolKick            int64 //被踢再次加入冷却
	CoolExit            int64 //主动退出再次加入冷却
	LastFlushTm         int64 //上次刷新的时间，用于公会微服务玩家的数据更新，采取的是延迟到玩家登陆或者跨天刷新机制，而不是服务器主动更新
	RedPGetedNum        int32 //当天领取的红包数
	RedPGetedRes        int32 //当天领取的红包总金额
	RedPGiveNum         int32 //当天已发的红包数
	RefuseAllInviteFlag int32 //拒绝所有邀请
	DissolveFlag        int32 //公会解散标识，1：表示公会已解散，并且玩家未登陆状态，用于客户端提示Tips用
	LiveDayRewardFlag   int32 //领取公会日活跃奖励标识：位标记。第i位0表示未领取，1表示领取
	AuctionBuyNum       int32 //统计拍卖行购买次数
	AuctionCost         int64 //统计拍卖行花费
}

func (p *GuildUser) IsChairman() bool {
	return p.Job == JobChairman
}

func (p *GuildUser) IsInGuild() bool {
	return p.GuildID != 0
}

func (p *GuildUser) IsChairAndVice() bool {
	return p.Job == JobChairman || p.Job == JobViceman
}

type CacheRedPacket struct {
	UUID int64
	ID   int32
}

//-------------------------------------------------------------------------------------------------------------------------------

// 拍卖品
type CacheAuctionGoods struct {
	UUID       int64
	ID         int32
	ActivityID int32
	ShareRate  int32
	List       []int64 //竞价UUID列表
}

func (p *CacheAuctionGoods) ListToString() string {
	if len(p.List) == 0 {
		return ""
	}
	str := ""
	for _, v := range p.List {
		str += fmt.Sprintf("%v,", v)
	}
	return str
}

func (p *CacheAuctionGoods) ListFromString(str string) {
	p.List = []int64{}
	vss := strings.Split(str, ",")

	for _, v := range vss {
		if uuid, err := strconv.ParseInt(v, 10, 64); err == nil {
			p.List = append(p.List, uuid)
		}
	}
}

// 一次竞价者
type CacheAuctionBider struct {
	UUID      int64
	UserID    int64
	GoodsUUID int64 //拍卖品UUID
	Bidprice  int32 //本次出价
	UnixTime  int64 //出价时间
}

// 魔将讨伐玩家数据
type CacheDCAgentInfo struct {
	UserID       int64
	Score        int64 //分数
	CoolBoss     int64 //挑战Boss冷却时间
	CoolRob      int64 //主动抢夺冷却时间
	CoolBeRobbed int64 //被抢冷却时间
}
