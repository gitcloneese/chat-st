package util

import (
	"time"
	"xy3-proto/pkg/log"
)

const (
	TimeParseMDY = "01/02/2006 15:04:05"
	TimeParseYMD = "2006/01/02 15:04:05"
)

func CurrentSeconds() int64 {
	return time.Now().Unix()
}

func CurrentMilliseconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func SecondsToTimeStr(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(TimeParseYMD)
}

func SecondsToTime(timestamp int64) (time.Time, error) {
	return StrToTime(SecondsToTimeStr(timestamp))
}

func GetTimeStamp() int64 {
	return time.Now().Unix()
}

func GetTimeStampMs() int64 {
	return time.Now().UnixNano() / (1000 * 1000)
}

func DiffTimeDay(early int64, later int64) bool {
	if early > later {
		early, later = later, early
	}

	t1 := time.Unix(early, 0)
	t2 := time.Unix(later, 0)

	if t1.Year() != t2.Year() || t1.Month() != t2.Month() || t1.Day() != t2.Day() {
		return true
	}

	return false
}

func GetDate() string {
	return FormatDate(GetTimeStamp())
}

func GetDay() string {
	return FormatDay(GetTimeStamp())
}

func FormatDay(timestamp int64) string {
	timeTemplate := "2006-01-02"
	return time.Unix(timestamp, 0).Format(timeTemplate)
}

func FormatDate(timestamp int64) string {
	timeTemplate := "2006-01-02 15:04:05"
	return time.Unix(timestamp, 0).Format(timeTemplate)
}

// MaxTime 服务器停服时间 2050-01-01 00:00:00
const ServerMaxTime = int64(2521814400)

// EncryptTime 时间戳裁剪，保留9位
func EncryptTime(timeStamp int64, desc bool) int64 {
	if desc {
		return ServerMaxTime - timeStamp // 当前时间戳越小，返回的时间越大
	} else {
		return 1000000000 - (ServerMaxTime - timeStamp) // 当前时间戳越小，返回的时间越小
	}
}

// DecryptTime 裁剪后时间戳还原
func DecryptTime(timeStamp int64, desc bool) int64 {
	if desc {
		return ServerMaxTime - timeStamp
	} else {
		return ServerMaxTime - (1000000000 - timeStamp)
	}
}

//--------------------------------------------------------

// GetTick 毫秒级时间戳
func GetTick() int64 {
	return time.Now().UnixNano() / 1e6
}

// 获取当前时间
func GetCurSec() int64 {
	return time.Now().Unix()
}

// 时间转字符串
func TimeToStr(Date time.Time) string {
	return Date.Format(TimeParseMDY)
}

// 字符串转时间
func StrToTime(Date string) (time.Time, error) {
	curTime := time.Now()
	return time.ParseInLocation(TimeParseMDY, Date, curTime.Location())
}

func ParseRedisTimeStr(timeStr string) (time.Time, error) {
	layout := "2006-01-02T15:04:05.9999999Z07:00"
	return time.Parse(layout, timeStr)
}

// 计算与现在的天数差
func CalcCurDaySub(Date string) int {
	daySub := 0
	if len(Date) == 0 {
		return daySub
	}
	curTime := time.Now()
	getTime, err := time.ParseInLocation(TimeParseMDY, Date, curTime.Location())
	if err != nil {
		return daySub
	}

	//跨年处理
	if getTime.Year() != curTime.Year() {
		curTime, err = time.ParseInLocation("2006-01-02", curTime.Format("2006-01-02"), curTime.Location())
		if err != nil {
			log.Error("CalcCurDaySub ParseInLocation Error:%v", err)
		}
		getTime, err = time.ParseInLocation("2006-01-02", getTime.Format("2006-01-02"), curTime.Location())
		if err != nil {
			log.Error("CalcCurDaySub ParseInLocation Error:%v", err)
		}

		daySub = int(curTime.Sub(getTime).Hours() / 24)
	} else {
		daySub = curTime.YearDay() - getTime.YearDay()
	}

	return daySub
}

// 判断是否是今天
func IsToday(Date string) bool {
	if len(Date) == 0 {
		return false
	}
	return CalcCurDaySub(Date) == 0
}

// 判断是否是今天
func IsTodayT(date time.Time) bool {
	now := time.Now()
	if now.Year() == date.Year() && now.Month() == date.Month() && now.Day() == date.Day() {
		return true
	}
	return false
}

// 判断时间是当年的第几周 ，用于处理周期事件
// 从周一开始算作新的一周
/*
存在一些跨年问题
        1 跨年那一周 会返回不同的值， 可能是 星期一返回54， 星期二跨年了 返回 1
        2 另外代码里单双周的判断跨年也会出错
        3 目前这个函数的 周期是 星期天-下个星期六是一个周期
*/
func WeekByDateDeprecated(t time.Time) int {
	yearDay := t.YearDay()

	yearFirstDay := t.AddDate(0, 0, -yearDay+1)
	firstDayInWeek := int(yearFirstDay.Weekday())

	//今年第一周有几天
	firstWeekDays := 1
	if firstDayInWeek != 0 {
		firstWeekDays = 7 - firstDayInWeek + 1
	}
	var week int
	if yearDay <= firstWeekDays {
		week = 1
	} else {
		week = (yearDay-firstWeekDays)/7 + 2
	}
	return week
}

// WeekByDate 判断时间是距离中国时区（东八区）2022年12月26日 0:0:0的第几周
// 2022-09-26 0:0:0(utc+8) = 1671984000
func WeekByDate(t time.Time) int {
	var (
		now            = t.Unix()
		baseTime int64 = 1671984000
	)
	if baseTime > now {
		return 0
	}
	return int((now-baseTime)/(3600*24*7) + 1)
}

// GetWeekDay 将0~6从星期天开始的转为1~7的周日期
func GetWeekDay(t time.Time) int32 {
	weekDay := int32(t.Weekday())
	if weekDay == 0 {
		return 7
	} else {
		return weekDay
	}
}
