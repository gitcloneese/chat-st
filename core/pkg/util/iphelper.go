package util

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func StringIpToInt(ipstring string) int {
	ipSegs := strings.Split(ipstring, ".")
	var ipInt int = 0
	var pos uint = 24
	for _, ipSeg := range ipSegs {
		tempInt, _ := strconv.Atoi(ipSeg)
		tempInt = tempInt << pos
		ipInt = ipInt | tempInt
		pos -= 8
	}
	return ipInt
}

func IpIntToString(ipInt int) string {
	ipSegs := make([]string, 4)
	var length int = len(ipSegs)
	buffer := bytes.NewBufferString("")
	for i := 0; i < length; i++ {
		tempInt := ipInt & 0xFF
		ipSegs[length-i-1] = strconv.Itoa(tempInt)
		ipInt = ipInt >> 8
	}
	for i := 0; i < length; i++ {
		buffer.WriteString(ipSegs[i])
		if i < length-1 {
			buffer.WriteString(".")
		}
	}
	return buffer.String()
}

// 获取本地的ip
func GetLocalIp() string {
	addrSlice, err := net.InterfaceAddrs()
	if nil != err {
		fmt.Printf("Get local IP addr failed!!!")
		return "localhost"
	}
	for _, addr := range addrSlice {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if nil != ipnet.IP.To4() {
				return ipnet.IP.String()
			}
		}
	}
	return "localhost"
}

// 获取客户端IP
func GetClientIP(ctx context.Context) string {
	c, ok := ctx.(*gin.Context)
	if !ok {
		return ""
	}

	return c.ClientIP()
}
