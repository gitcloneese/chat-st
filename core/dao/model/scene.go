package model

import (
	"fmt"
	"strconv"
	"strings"
)

type CacheRoleFieldEnum string

func (s CacheRoleFieldEnum) String() string {
	return string(s)
}

const (
	CacheRoleFieldID         CacheRoleFieldEnum = "ID"
	CacheRoleFieldNick       CacheRoleFieldEnum = "Nick"
	CacheRoleFieldSex        CacheRoleFieldEnum = "Sex"
	CacheRoleFieldLevel      CacheRoleFieldEnum = "Level"
	CacheRoleFieldExp        CacheRoleFieldEnum = "Exp"
	CacheRoleFieldPower      CacheRoleFieldEnum = "Power"
	CacheRoleFieldHeadID     CacheRoleFieldEnum = "HeadID"
	CacheRoleFieldFrameID    CacheRoleFieldEnum = "FrameID"
	CacheRoleFieldDrawID     CacheRoleFieldEnum = "DrawID"
	CacheRoleFieldTitle      CacheRoleFieldEnum = "Title"
	CacheRoleFieldLoginTime  CacheRoleFieldEnum = "LoginTime"
	CacheRoleFieldLogoutTime CacheRoleFieldEnum = "LogoutTime"
	CacheRoleFieldServerID   CacheRoleFieldEnum = "ServerID"
	CacheRoleFieldLoginDays  CacheRoleFieldEnum = "LoginDays"
	CacheRoleFieldDailyTime  CacheRoleFieldEnum = "DailyTime"
	CacheRoleFieldIsRobot    CacheRoleFieldEnum = "IsRobot"
)

type CacheRole struct {
	ID         int64  // 角色id
	Nick       string // 角色名
	Sex        int32  // 性别
	Level      int32  // 等级
	Exp        int64  // 经验
	Power      int64  // 战力
	HeadID     int32  // 头像id
	FrameID    int32  // 边框id
	DrawID     int32  // 立绘
	Title      int32  // 头衔
	LoginTime  int64  // 登陆时间
	LogoutTime int64  // 登出时间
	ServerID   int64  // 服务器id
	LoginDays  int32  // 累计登陆天数
	DailyTime  int64  // 每日在线时长(秒)
	ClientIp   string // 客户端ip
	RegTime    int64  // 注册时间
	LockTime   int64  // 封号截至时间

	IsRobot int32  `bson:"isrobot"` // 是否机器人
	Os      int32  // 每次登录时设置Os
	UnionId string // 账户id
}

type CacheScenePrime struct {
	NowMaxLevel int32 //本服最高等级
}

type UserInfo struct {
	UnionID  string `bson:"unionid"`
	PlayerID int64  `bson:"playerid"`
}

type CacheExclusiveWeapon struct {
	ID      int32
	Star    int32
	HoleNum int32
	SuitID  int32
	Rune    []int32
	VBar    [][2]int32
}

func (p *CacheExclusiveWeapon) RunesToString() string {
	if p.Rune == nil {
		return ""
	}
	str := ""
	for _, v := range p.Rune {
		str += fmt.Sprintf("%v,", v)
	}
	return str
}

func (p *CacheExclusiveWeapon) RunesFromString(str string) {
	p.Rune = []int32{}

	vstr := strings.Split(str, ",")
	for _, one := range vstr {
		v, _ := strconv.ParseInt(one, 10, 64)
		p.Rune = append(p.Rune, int32(v))
	}
}

func (p *CacheExclusiveWeapon) VBarToString() string {
	if p.VBar == nil {
		return ""
	}
	str := ""
	for _, v := range p.VBar {
		str += fmt.Sprintf("%v,%v#", v[0], v[1])
	}
	return str
}

func (p *CacheExclusiveWeapon) VBarFromString(str string) {
	p.VBar = [][2]int32{}

	vstr := strings.Split(str, "#")
	for _, one := range vstr {
		v2 := strings.Split(one, ",")
		if len(v2) == 2 {
			k, _ := strconv.ParseInt(v2[0], 10, 64)
			n, _ := strconv.ParseInt(v2[1], 10, 64)
			p.VBar = append(p.VBar, [2]int32{int32(k), int32(n)})
		}
	}
}

type CacheLingbao struct {
	ID      int32
	Advance int32
	Star    int32
	VBar    []int32
}

func (p *CacheLingbao) VBarsToString() string {
	if p.VBar == nil {
		return ""
	}
	str := ""
	for _, v := range p.VBar {
		str += fmt.Sprintf("%v,", v)
	}
	return str
}

func (p *CacheLingbao) VBarsFromString(str string) {
	p.VBar = []int32{}

	vstr := strings.Split(str, ",")
	for _, one := range vstr {
		v, _ := strconv.ParseInt(one, 10, 64)
		p.VBar = append(p.VBar, int32(v))
	}
}

type CachePenglaiHero struct {
	ID      int32 `json:"ID"`
	Star    int32 `json:"Star"`
	Quality int32 `json:"Quality"`
	Awaken  int32 `json:"Awaken"`
}
type CachePenglaiOne struct {
	Order   int32               `json:"Order"`
	Heros   []*CachePenglaiHero `json:"Heros"`
	Unlocks []int32             `json:"Unlocks"`
}
type CachePenglai struct {
	List []*CachePenglaiOne `json:"List"`
}

type CacheBannerInfoAll struct {
	BannersInUse [4]*CacheBannerInfo `json:"BannersInUse"`
}

type CacheBannerInfo struct {
	Index             int32           `json:"Index"`
	ID                int32           `json:"ID"`
	Quality           int32           `json:"Quality"`
	Star              int32           `json:"Star"`
	Subs              []int32         `json:"Subs"`
	Runes             [5]int32        `json:"Runes"`
	AttributeUpgrades map[int32]int64 `json:"AttributeUpgrades"`
}

type CacheConstellation struct {
	Points map[int32]int32 `json:"Points"`
	Skills map[int32]int32 `json:"Skills"`
}
