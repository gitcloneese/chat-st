package util

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/stringx"
)

const (
	randomLen       = 16
	millisPerSecond = 1000
)

type distLockOpt struct {
	lockKey       string
	ttlSecs       int  // in secs
	releaseOnDone bool // if no release, lock would be avaiable again after ttl
	printLog      bool
}

type DistLockOption func(o *distLockOpt)

func WithLockKey(k string) DistLockOption {
	return func(o *distLockOpt) {
		o.lockKey = k
	}
}

func WithTTLSecs(secs int) DistLockOption {
	return func(o *distLockOpt) {
		o.ttlSecs = secs
	}
}

func WithReleaseOnDone(p bool) DistLockOption {
	return func(o *distLockOpt) {
		o.releaseOnDone = p
	}
}

func WithPrintLog(p bool) DistLockOption {
	return func(o *distLockOpt) {
		o.printLog = p
	}
}

type DistLock interface {
	// default ttl 5 secs, release on done & would print log
	TryDo(f func() error, opts ...DistLockOption) error
}

type redisDistLock struct {
	client *redis.Client
}

func NewRedisDistLock(client *redis.Client) DistLock {
	return &redisDistLock{client}
}

func (l *redisDistLock) TryDo(f func() error, opts ...DistLockOption) error {
	opt := &distLockOpt{
		ttlSecs:       5,
		releaseOnDone: true,
		printLog:      true,
	}

	for _, o := range opts {
		o(opt)
	}
	if opt.lockKey == "" {
		return fmt.Errorf("invalid empty lock key provided")
	}
	return l.tryDo(f, opt)
}

func (l *redisDistLock) tryDo(f func() error, opt *distLockOpt) error {
	success, token := l.acquire(opt.lockKey, opt.ttlSecs)
	if !success && opt.printLog {
		log.Debugf("failed to acquire dist lock for key: %s", opt.lockKey)
		return nil
	}
	if opt.printLog {
		log.Debugf("dist lock acquired successfully for key: %s, token: %s, ttl: %d\n", opt.lockKey, token, opt.ttlSecs)
	}
	defer func() {
		if opt.releaseOnDone {
			tolerance := rand.Intn(500) // milliseconds
			time.Sleep(time.Duration(tolerance) * time.Millisecond)
			success := l.release(opt.lockKey, token)
			if !success && opt.printLog {
				log.Debugf("failed to release dist lock for key: %s ,using token: %s", opt.lockKey, token)
				return
			}
			if opt.printLog {
				log.Debugf("dist lock manually released successfully for key: %s, token: %s, ttl: %d\n", opt.lockKey, token, opt.ttlSecs)
			}
		}
	}()
	if err := f(); err != nil {
		return fmt.Errorf("failed to exec func, err: %w", err)
	}
	return nil
}

// acquire acquires the lock.
func (l *redisDistLock) acquire(k string, lockSecs int) (success bool, token string) {
	script := `return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])`
	token = stringx.Randn(randomLen)
	acquire := redis.NewScript(script)
	tolerance := rand.Intn(500) // milliseconds
	res, err := acquire.Run(context.Background(), l.client, []string{k}, []string{token, strconv.Itoa(lockSecs*millisPerSecond + tolerance)}).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Errorf("Error on acquiring lock for key: %s, %v", k, err)
		}
		return false, ""
	}
	if res == nil {
		return false, ""
	}
	reply, ok := res.(string)
	if !ok || reply != "OK" {
		log.Errorf("Unknown reply when acquiring lock for key: %s, res: %v", k, res)
		return false, ""
	}
	return true, token
}

// release releases the lock.
func (l *redisDistLock) release(k, token string) (success bool) {
	script := `
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end
`
	release := redis.NewScript(script)
	res, err := release.Run(context.Background(), l.client, []string{k}, []string{token}).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Errorf(" DistLock failed to release release for key:%s, error: %v", k, err)
		}
		return false
	}
	reply, ok := res.(int64)
	if !ok {
		return false
	}
	return reply == 1
}
