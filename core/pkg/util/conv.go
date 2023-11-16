package util

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"xy3-proto/pkg/log"
)

var (
	escapeCharMap map[string]string
)

func init() {
	escapeCharMap = map[string]string{
		"\\":  "",
		"\a":  " ",
		"\b":  " ",
		"\f":  " ",
		"\n":  " ",
		"\r":  " ",
		"\t":  " ",
		"\v":  " ",
		"\\?": "?",
		"\\0": "0",
	}
}

//---------------------------------------------------------------------------------------------

// string转到int
func StrToInt(src string) int {
	if src == "" {
		return 0
	}
	//log.Printf("StrToInt:%s.", src)
	dst, err := strconv.Atoi(src)
	if err != nil {
		log.Warn("str to int error(%v).", err)
		return 0
	}
	return dst
}

// string转到int64
func StrToInt64(src string) int64 {
	//	log.Print("src:%s.", src)
	if src == "" {
		return 0
	}
	dst, err := strconv.ParseInt(src, 10, 64)
	if err != nil {
		log.Warn("str to int64 error(%v).", err)
		return 0
	}
	return dst
}

// string转到int32
func StrToInt32(src string) int32 {
	if src == "" {
		return 0
	}
	dst, err := strconv.Atoi(src)
	if err != nil {
		log.Warn("str to int32 error(%v).", err)
		return 0
	}
	return int32(dst)
}

// string转float64
func StrToFloat64(src string) float64 {
	if src == "" {
		return 0
	}
	dst, err := strconv.ParseFloat(src, 64)
	if err != nil {
		log.Warn("srt to float64 err(%v).", err)
		return 0
	}
	return dst
}

// string转到bool
func StrToBool(src string) bool {
	if src == "" {
		return false
	}
	dst, err := strconv.Atoi(src)
	if err != nil {
		log.Warn("srt to bool err(%v).", err)
		return false
	}
	return dst != 0
}

// ---------------------------------------------------------------------------------------------
// int转到string
func IntToStr(src int) string {
	return strconv.Itoa(src)
}

// int64转到string
func Int64ToStr(src int64) string {
	return strconv.FormatInt(src, 10)
}

// int转到string
func Int32ToStr(src int32) string {
	return strconv.Itoa(int(src))
}

// ToString .
func ToString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

func Md5(data string) string {
	v := md5.New()
	v.Write([]byte(data))
	md5Data := v.Sum([]byte(""))
	return hex.EncodeToString(md5Data)
}

// ToJSON json string
func ToJSON(v interface{}) string {
	j, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(j)
}

// 转成base64格式字符串
func ToBase64String(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func ReplaceEscapeChar(src string) (des string) {
	des = src
	for k, v := range escapeCharMap {
		des = strings.Replace(des, k, v, -1)
	}
	return
}

func StructToJson(data interface{}) string {
	if reflect.TypeOf(data).Kind() == reflect.Ptr {
		data = reflect.ValueOf(data).Elem().Interface()
	}
	if bytes, err := json.Marshal(data); err != nil {
		return ""
	} else {
		str := string(bytes)
		return str
	}
}

func StructToJsonIndent(data interface{}) string {
	if reflect.TypeOf(data).Kind() == reflect.Ptr {
		data = reflect.ValueOf(data).Elem().Interface()
	}
	if bytes, err := json.MarshalIndent(data, "", " "); err != nil {
		return ""
	} else {
		str := string(bytes)
		return str
	}
}
