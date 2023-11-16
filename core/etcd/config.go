package etcd

import (
	"context"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"xy3-proto/pkg/conf/paladin"
	"xy3-proto/pkg/log"
)

func Update(cli *clientv3.Client, key string, value string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	_, err = cli.Put(ctx, key, value)
	cancel()
	if err != nil {
		log.Error("etcd Update err:%v", err)
		return
	}
	return
}

func Read(cli *clientv3.Client, key string) (resp *clientv3.GetResponse, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	resp, err = cli.Get(ctx, key)
	cancel()
	if err != nil {
		log.Error("etcd Read get err:%v", err)
		return
	}
	return
}

func UpdateConfig(cli *clientv3.Client, key string, value string) (err error) {
	resp, err := Read(cli, key)
	if err != nil {
		log.Error("etcd UpdateConfig Read key:%v err:%v", key, err)
		return err
	}
	if resp.Count != 0 {
		return
	}
	err = Update(cli, key, value)
	if err != nil {
		log.Error("etcd UpdateConfig Update Key:%v Value:%v err:%v", key, value, err)
		return err
	}
	log.Info("etcd UpdateConfig key:%v", key)
	return
}

func WatchConfig(cli *clientv3.Client, setters map[string]paladin.Setter) {
	for key, setter := range setters {
		go func(cli *clientv3.Client, key string, setter paladin.Setter) {
			rch := cli.Watch(context.Background(), key) // type WatchChan <-chan WatchResponse
			for wresp := range rch {
				for _, ev := range wresp.Events {
					if ev.Type == mvccpb.PUT {
						err := setter.Set(ev.Kv.Value)
						if err != nil {
							log.Error("etcd WatchConfig Config Set err:%v key:%v value:%v", err, ev.Kv.Key, ev.Kv.Value)
						}
					}
				}
			}
		}(cli, key, setter)
	}
}
