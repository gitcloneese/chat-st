package component

import (
	"context"

	v8 "github.com/go-redis/redis/v8"
	"xy3-proto/pkg/conf/paladin"
	"xy3-proto/pkg/log"
)

// 创建Redis对象
func NewRedis() (*v8.Client, error) {
	var cfg struct {
		Client *v8.Options
	}

	if err := paladin.Get("redis.txt").UnmarshalTOML(&cfg); err != nil {
		log.Error("redis.txt err %v", err)
		return nil, err
	}
	log.Debug("redis.txt %+v", cfg.Client)
	r := v8.NewClient(cfg.Client)
	res, err := r.Ping(context.Background()).Result()
	if err != nil {
		log.Error("Redis Ping err %v", err)
		return nil, err
	}
	log.Info("NewRedis Ping Result %v", res)

	return r, nil
}

// NewRedisCluster .
func NewRedisCluster() (r *v8.ClusterClient, err error) {
	var cfg struct {
		ClusterClient *v8.ClusterOptions
	}
	err = paladin.Get("redis.txt").UnmarshalTOML(&cfg)
	if err != nil {
		log.Error("NewRedisCluster err:%v", err)
		return
	}
	log.Debug("redis.txt %+v", cfg.ClusterClient)
	r = v8.NewClusterClient(cfg.ClusterClient)
	res, err := r.Ping(context.Background()).Result()
	if err != nil {
		log.Error("NewRedisCluster Ping err:%v", err)
		return
	}
	log.Info("NewRedisCluster Ping Result %v", res)
	return
}
