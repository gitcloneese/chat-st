package main

import (
	"context"
	"fmt"
	v9 "github.com/redis/go-redis/v9"
	"log"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	rdb := v9.NewClient(&v9.Options{
		Addr:     Addr,
		Username: User,
		Password: Passwd, // no password set
		DB:       DB,     // use default DB
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("开始压测redis...\n")
	now := time.Now()
	var i int
	wg := &sync.WaitGroup{}
	wg.Add(n)
	for i < n {
		i++
		go test(wg, rdb)
	}
	wg.Wait()
	qs := atomic.LoadInt64(&Qs)
	t := time.Since(now).Seconds()
	log.Printf("压测完成 qs:%v time:%v qps:%v\n", qs, t, float64(qs)/t)
	os.Exit(0)
}

func Key() string {
	r := rand.Intn(100000)
	return fmt.Sprintf("test:%v-%v", time.Now().UnixNano(), r)
}

const (
	Set = iota
	Incr
	HSet
	HGet
)

func test(wg *sync.WaitGroup, client *v9.Client) {
	defer wg.Done()
	key := Key()
	var err error
	switch T {
	case Set:
		err = set(wg, client)
	case Incr:
		err = incr(wg, client)
	case HSet:
		err = hSet(wg, client, key)
	case HGet:
		err = hGet(wg, client)
	}
	if err != nil {
		log.Print(err)
	}
}

func set(wg *sync.WaitGroup, client *v9.Client) error {
	defer wg.Done()
	err := client.Set(context.TODO(), Key(), 1, time.Second*500).Err()
	atomic.AddInt64(&Qs, 1)
	return err
}

func incr(wg *sync.WaitGroup, client *v9.Client) error {
	defer wg.Done()
	key := "incr"
	n := int64(rand.Int31n(5) + 1)
	err := client.IncrBy(context.TODO(), key, n).Err()
	atomic.AddInt64(&Qs, 1)
	return err
}

func hSet(wg *sync.WaitGroup, client *v9.Client, keys ...string) error {
	var key string
	if len(key) == 0 {
		key = keys[0]
	} else {
		key = Key()
	}
	defer wg.Done()
	err := client.HSet(context.TODO(), key, Key(), 1).Err()
	atomic.AddInt64(&Qs, 1)
	return err
}

func hGet(wg *sync.WaitGroup, client *v9.Client, keys ...string) error {
	var key string
	if len(key) == 0 {
		key = keys[0]
	} else {
		key = Key()
	}
	defer wg.Done()
	_, err := client.HGet(context.TODO(), key, "incr").Result()
	atomic.AddInt64(&Qs, 1)
	return err
}
