package redsync

import (
	"context"
	"errors"
	"strings"
	"time"

	redis "github.com/go-redis/redis/v8"
)

type pool struct {
	delegate *redis.Client
}

func (p *pool) Get(ctx context.Context) (Conn, error) {
	c := p.delegate
	if ctx != nil {
		c = c.WithContext(ctx)
	}
	return &conn{c}, nil
}

// NewPool returns a Goredis-based pool implementation.
func NewPool(delegate *redis.Client) Pool {
	return &pool{delegate}
}

type conn struct {
	delegate *redis.Client
}

func (c *conn) Get(ctx context.Context, name string) (string, error) {
	value, err := c.delegate.Get(ctx, name).Result()
	return value, noErrNil(err)
}

func (c *conn) Set(ctx context.Context, name string, value string) (bool, error) {
	reply, err := c.delegate.Set(ctx, name, value, 0).Result()
	return reply == "OK", noErrNil(err)
}

func (c *conn) SetNX(ctx context.Context, name string, value string, expiry time.Duration) (bool, error) {
	ok, err := c.delegate.SetNX(ctx, name, value, expiry).Result()
	return ok, noErrNil(err)
}

func (c *conn) PTTL(ctx context.Context, name string) (time.Duration, error) {
	expiry, err := c.delegate.PTTL(ctx, name).Result()
	return expiry, noErrNil(err)
}

func (c *conn) Eval(ctx context.Context, script *Script, keysAndArgs ...interface{}) (interface{}, error) {
	keys := make([]string, script.KeyCount)
	args := keysAndArgs

	if script.KeyCount > 0 {
		for i := 0; i < script.KeyCount; i++ {
			keys[i] = keysAndArgs[i].(string)
		}

		args = keysAndArgs[script.KeyCount:]
	}

	v, err := c.delegate.EvalSha(ctx, script.Hash, keys, args...).Result()
	if err != nil && strings.HasPrefix(err.Error(), "NOSCRIPT ") {
		v, err = c.delegate.Eval(ctx, script.Src, keys, args...).Result()
	}
	return v, noErrNil(err)
}

func (c *conn) Close() error {
	// Not needed for this library
	return nil
}

func noErrNil(err error) error {
	if !errors.Is(err, redis.Nil) {
		return err
	}
	return nil
}
