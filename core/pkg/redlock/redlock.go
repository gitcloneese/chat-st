package redlock

import (
	"context"
	crand "crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
	"xy3-proto/pkg/log"

	redis "github.com/go-redis/redis/v8"
)

const (
	// DefaultRetryCount is the max retry times for lock acquire
	DefaultRetryCount = 10

	// DefaultRetryDelay is upper wait time in millisecond for lock acquire retry
	DefaultRetryDelay = 200

	// ClockDriftFactor is clock drift factor, more information refers to doc
	ClockDriftFactor = 0.01

	// UnlockScript is redis lua script to release a lock
	UnlockScript = `
        if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("del", KEYS[1])
        else
            return 0
        end
        `
)

// RedLock holds the redis lock
type RedLock struct {
	retryCount  int
	retryDelay  int
	driftFactor float64

	clients []*RedClient
	quorum  int

	cache KVCache
}

// RedClient holds client to redis
type RedClient struct {
	addr   string
	user   string
	passwd string
	cli    *redis.Client
}

func parseConnString(addr string) (*redis.Options, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	opts := &redis.Options{
		Network: u.Scheme,
		Addr:    u.Host,
	}

	dbStr := strings.Trim(u.Path, "/")
	if dbStr == "" {
		dbStr = "0"
	}
	db, err := strconv.Atoi(dbStr)
	if err != nil {
		return nil, err
	}
	opts.DB = db

	password, ok := u.User.Password()
	if ok {
		opts.Password = password
	}

	for k, v := range u.Query() {
		if k == "DialTimeout" {
			timeout, err := strconv.Atoi(v[0])
			if err != nil {
				return nil, err
			}
			opts.DialTimeout = time.Duration(timeout)
		}
		if k == "ReadTimeout" {
			timeout, err := strconv.Atoi(v[0])
			if err != nil {
				return nil, err
			}
			opts.ReadTimeout = time.Duration(timeout)
		}
		if k == "WriteTimeout" {
			timeout, err := strconv.Atoi(v[0])
			if err != nil {
				return nil, err
			}
			opts.WriteTimeout = time.Duration(timeout)
		}
	}

	return opts, nil
}

// NewRedLock creates a RedLock
func NewRedLock(addrs []string) (*RedLock, error) {
	if len(addrs)%2 == 0 {
		return nil, fmt.Errorf("error redis server list: %d", len(addrs))
	}

	var clients []*RedClient
	for _, addr := range addrs {
		opts, err := parseConnString(addr)
		if err != nil {
			return nil, err
		}
		cli := redis.NewClient(opts)
		val, err := cli.Ping(context.Background()).Result()
		if err != nil {
			log.Error("NewRedis ping val:%s err:%+v", val, err)
			panic(err)
		}
		clients = append(clients, &RedClient{addr: addr, cli: cli})
	}

	return &RedLock{
		retryCount:  DefaultRetryCount,
		retryDelay:  DefaultRetryDelay,
		driftFactor: ClockDriftFactor,
		quorum:      len(addrs)/2 + 1,
		clients:     clients,
		cache:       NewCacheImpl(CacheTypeSimple, nil),
	}, nil
}

// NewRedLockWithPasswd creates a RedLock
func NewRedLockWithPasswd(options []*redis.Options) (*RedLock, error) {
	var clients []*RedClient
	for _, opt := range options {
		cli := redis.NewClient(opt)
		val, err := cli.Ping(context.Background()).Result()
		if err != nil {
			log.Error("NewRedLockWithPasswd ping val:%s err:%+v", val, err)
			panic(err)
		}
		clients = append(clients, &RedClient{opt.Addr, opt.Username, opt.Password, cli})
	}

	return &RedLock{
		retryCount:  DefaultRetryCount,
		retryDelay:  DefaultRetryDelay,
		driftFactor: ClockDriftFactor,
		quorum:      len(options)/2 + 1,
		clients:     clients,
		cache:       NewCacheImpl(CacheTypeSimple, nil),
	}, nil
}

// SetCache resets cache based on cache type
func (r *RedLock) SetCache(cacheType string, opts map[string]interface{}) {
	r.cache = NewCacheImpl(cacheType, opts)
}

// SetRetryCount sets acquire lock retry count
func (r *RedLock) SetRetryCount(count int) {
	if count <= 0 {
		return
	}
	r.retryCount = count
}

// SetRetryDelay sets acquire lock retry max internal in millisecond
func (r *RedLock) SetRetryDelay(delay int) {
	if delay <= 0 {
		return
	}
	r.retryDelay = delay
}

func getRandStr() string {
	b := make([]byte, 16)
	_, err := crand.Read(b)
	if err != nil {
		log.Error("getRandStr err:", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

func lockInstance(client *RedClient, resource string, val string, ttl int, c chan bool) {
	if client.cli == nil {
		c <- false
		return
	}
	reply := client.cli.SetNX(context.Background(), resource, val, time.Duration(ttl)*time.Millisecond)
	if reply.Err() != nil || !reply.Val() {
		log.Warn("lockInstance cli:%v key:%s val:%v resval:%v error:%v", client, resource, val, reply.Val(), reply.Err())
		c <- false
		return
	}
	c <- true
}

func unlockInstance(client *RedClient, resource string, val string, c chan bool) {
	if client.cli != nil {
		client.cli.Eval(context.Background(), UnlockScript, []string{resource}, val)
	}
	c <- true
}

// Lock acquires a distribute lock
func (r *RedLock) Lock(resource string, ttl int) (int64, error) {
	now := time.Now()
	defer func(now time.Time) {
		_metricRedLockReqDur.Observe(int64(time.Since(now)/time.Millisecond), resource)
	}(now)

	val := getRandStr()
	for i := 0; i < r.retryCount; i++ {
		c := make(chan bool, len(r.clients))
		success := 0
		start := time.Now()

		for _, cli := range r.clients {
			go lockInstance(cli, resource, val, ttl, c)
		}
		for j := 0; j < len(r.clients); j++ {
			if <-c {
				success++
			}
		}

		drift := int(float64(ttl)*r.driftFactor) + 2
		costTime := time.Since(start).Nanoseconds() / 1e6
		validityTime := int64(ttl) - costTime - int64(drift)
		if success >= r.quorum && validityTime > 0 {
			_, err := r.cache.Set(resource, val, validityTime)
			if err != nil {
				log.Error("Lock error:%v", err)
			}
			return validityTime, err
		}
		cul := make(chan bool, len(r.clients))
		for _, cli := range r.clients {
			go unlockInstance(cli, resource, val, cul)
		}
		for j := 0; j < len(r.clients); j++ {
			<-cul
		}
		// Wait a random delay before to retry
		time.Sleep(time.Duration(rand.Intn(r.retryDelay)) * time.Millisecond)
	}

	return 0, errors.New("failed to require lock")
}

// UnLock releases an acquired lock
func (r *RedLock) UnLock(resource string) error {
	elem, err := r.cache.Get(resource)
	if err != nil {
		return err
	}
	if elem == nil {
		return nil
	}
	defer r.cache.Delete(resource)
	c := make(chan bool, len(r.clients))
	for _, cli := range r.clients {
		go unlockInstance(cli, resource, elem.Val, c)
	}
	for i := 0; i < len(r.clients); i++ {
		<-c
	}
	return nil
}
