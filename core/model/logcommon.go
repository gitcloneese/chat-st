package model

import (
	"encoding/json"

	"xy3-proto/pkg/log"
)

const (
	Operate_Nil     int32 = 0 //未知操作
	Operate_New     int32 = 1 //新增操作
	Operate_Add     int32 = 2 //增加操作
	Operate_Del     int32 = 3 //删除操作
	Operate_Change  int32 = 4 //修改操作
	Operate_Overdue int32 = 5 //过期操作
)

type LogItem struct {
	Type int32 `json:"type"`
	ID   int32 `json:"id"`
	Num  int64 `json:"num"`
}

type LogEmbed struct {
	ID        uint  `gorm:"primary_key;auto_increment"`
	CreatedAt int64 `gorm:"Time"`
	OS        int32 `gorm:"OS"`
}

func JsonListLogItem(list []*LogItem) string {
	buffer, err := json.Marshal(list)
	if err != nil {
		log.Error("JsonListLogItem err:%v list:%v", err, list)
		return "[]"
	}
	return string(buffer)
}
