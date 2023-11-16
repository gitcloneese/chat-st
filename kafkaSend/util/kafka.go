package util

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
	"x-server/core/apollo"
	"x-server/core/mq/kafka"

	pbCommon "xy3-proto/common"
	pb "xy3-proto/logger"
	"xy3-proto/pkg/conf/paladin"
	"xy3-proto/pkg/log"
)

type KafkaConfig struct {
	Addrs []string `toml:"addrs"`
}

var (
	sender     = new(sync.Map)
	_kfaConfig *KafkaConfig
)

func loadConfig() (*KafkaConfig, error) {
	var (
		ct  paladin.TOML
		cfg *KafkaConfig
	)
	// apollo配置优先
	v := apollo.Get(apollo.KafkaNS)
	if v == nil || v.Unmarshal(&ct) != nil {
		err := paladin.Get(apollo.KafkaNS).Unmarshal(&ct)
		if err != nil {
			log.Error("NewRedis UnmarshalTOML error: %v", err)
			return nil, err
		}
	}

	if err := ct.Get("Kafka").UnmarshalTOML(&cfg); err != nil {
		log.Warn("kafka cfg Kafka err %v", err)
		return nil, err
	}
	return cfg, nil
}

func newKafka(topic string) (kafka.Sender, error) {
	if _kfaConfig == nil {
		cfg, err := loadConfig()
		if err != nil {
			log.Warn("kafka loadConfig err %v", err)
			return nil, err
		}
		_kfaConfig = cfg
	}

	sender, err := kafka.NewKafkaSender(_kfaConfig.Addrs, topic)
	if err != nil {
		log.Warn("new kafka producer err %v", err)
		return nil, err
	}
	return sender, nil
}

// Record
// 写记录
func Record(req *pb.LogMsgs) error {
	return pubKafka(pb.LogServerTopic, req)
}

// RecordSingle
// 写单个记录
func RecordSingle(category pb.ELogCategory, os pbCommon.OSType, data interface{}) {
	log.Debug("RecordSingle category:%v", category)
	var buf []byte
	if data != nil {
		var err error
		buf, err = json.Marshal(data)
		if err != nil {
			log.Error("RecordSingle Data Marshal err %v", err)
			return
		}
	}

	msg := &pb.LogMsgs{
		Messages: []*pb.LogMsg{
			{
				Category: category,
				Os:       os,
				Time:     time.Now().Unix(),
				Json:     string(buf),
			},
		},
	}
	err := Record(msg)
	if err != nil {
		log.Error("RecordMany Record Error: %v", err)
	}
}

func RecordMany(category pb.ELogCategory, os pbCommon.OSType, list []interface{}) {
	// nowTime := time.Now().Format(model.DateLayout)
	nowTime := time.Now().Unix()
	msg := &pb.LogMsgs{Messages: make([]*pb.LogMsg, 0, len(list))}

	for _, obj := range list {
		if obj == nil {
			continue
		}
		buffer, err := json.Marshal(obj)
		if err != nil {
			log.Error("RecordMany marshal err! err:%v obj:%v", err, obj)
			continue
		}
		msg.Messages = append(msg.Messages, &pb.LogMsg{
			Category: category,
			Os:       os,
			Time:     nowTime,
			Json:     string(buffer),
		})
	}
	err := Record(msg)
	if err != nil {
		log.Error("RecordMany Record Error: %v", err)
	}
}

func pubKafka(topic string, req *pb.LogMsgs) error {
	send, exist := sender.Load(topic)
	if !exist {
		send, err := newKafka(topic)
		if err != nil {
			log.Error("pubKafka newKafka err %v", err)
			return err
		}
		sender.Store(topic, send)
	}

	data, err := json.Marshal(req)
	if err != nil {
		log.Warn("pubKafka Marshal err %v", err)
		return err
	}
	if len(data) == 0 {
		log.Warn("pubKafka data is Zero")
		return fmt.Errorf("pubKafka data is Zero")
	}

	msg := kafka.NewMessage(uuid.New().String(), data)
	if err := send.(kafka.Sender).Send(msg); err != nil {
		log.Error("pushKafka topic:%v err %v", topic, err)
		return err
	}

	return nil
}
