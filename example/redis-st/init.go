package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

// 每个玩家默认1s发送一个聊天
func init() {
	rand.Seed(time.Now().Unix())
	addFlag(flag.CommandLine)
}

var (
	n int

	Qs     int64
	Addr   string
	User   string
	Passwd string
	DB     int
	T      int
)

func addFlag(fs *flag.FlagSet) {
	fs.IntVar(&n, "n", 1000, fmt.Sprint("协程并发度"))
	fs.IntVar(&DB, "db", 0, fmt.Sprint("测试redis数据库,默认0库"))
	fs.StringVar(&Addr, "addr", "127.0.0.1:6379", fmt.Sprint("redis地址 默认127.0.0.1:6379"))
	fs.StringVar(&Passwd, "passwd", "", fmt.Sprint("redis密码"))
	fs.StringVar(&User, "user", "", fmt.Sprint("redis账户"))
	fs.IntVar(&T, "t", 0, "测试类型 默认测试set 1:incr 2:hset 3:hget 4...")
	flag.Parse()
}
