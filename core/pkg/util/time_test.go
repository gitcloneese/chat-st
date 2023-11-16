package util

import (
	"fmt"
	"testing"
	"time"
)

func TestIsToday(t *testing.T) {
	fmt.Println(IsToday(time.Now().Local().String()))
}

func TestCalcCurDaySub(t *testing.T) {
	now := time.Now()
	day := time.Date(now.Year(), now.Month(), now.Day()-2, 0, 0, 0, 0, now.Location())
	fmt.Println(CalcCurDaySub(TimeToStr(day)))

}

func TestSecondsToTime(t *testing.T) {
	unix := time.Now().Unix()
	fmt.Println(SecondsToTime(unix - 3600))
}

func TestGetWeekDay(t *testing.T) {
	now := time.Now()
	now = now.AddDate(0, 0, 6)
	fmt.Println(GetWeekDay(now))
}
