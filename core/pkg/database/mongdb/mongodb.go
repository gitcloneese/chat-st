package mongdb

import (
	"xy3-proto/pkg/log"
	"xy3-proto/pkg/net/netutil/breaker"
	xtime "xy3-proto/pkg/time"
)

// Config sqlserver config.
type Config struct {
	URI             string          // likes mongodb://foo:bar@localhost:27017
	ConnectTimeout  xtime.Duration  // connection mongodb timeout
	QueryTimeout    xtime.Duration  // query mongodb timeout
	ExecTimeout     xtime.Duration  // execute mongodb timeout
	Breaker         *breaker.Config // breaker
	MaxPoolSize     uint64
	MaxConnIdleTime xtime.Duration
}

// NewMongoDB NewSQLServer new db and retry connection when has error.
func NewMongoDB(c *Config) (db *DB, err error) {
	if c.ConnectTimeout == 0 || c.QueryTimeout == 0 || c.ExecTimeout == 0 {
		panic("mongo must be set query/execute/connect timeout")
	}
	db, err = Open(c)
	if err != nil {
		log.Error("open mongodb error(%v)", err)
		return
	}
	return
}
