package dao

import (
	"x-server/core/apollo"
	"x-server/core/mq/kafka"
	"x-server/core/pkg/util"
	pb "xy3-proto/logger"
	"xy3-proto/pkg/conf/paladin"
	"xy3-proto/pkg/log"
)

func NewKafkaReceiver() (kafka.Receiver, func(), error) {
	ct := paladin.TOML{}
	cfg := util.KafkaConfig{}
	v := apollo.Get(apollo.KafkaNS)
	if v == nil || v.Unmarshal(&ct) != nil {
		if err := paladin.Get(apollo.KafkaNS).Unmarshal(&ct); err != nil {
			log.Error("kafka.txt cfg err %v", err)
			return nil, nil, err
		}
	}
	if err := ct.Get("Kafka").UnmarshalTOML(&cfg); err != nil {
		log.Error("kafka cfg Kafka err %v", err)
		return nil, nil, err
	}

	log.Info("kafka broker: %v", cfg.Addrs)
	receiver, err := kafka.NewKafkaReceiver(cfg.Addrs, pb.LogServerTopic, pb.GroupIdWriter)
	if err != nil {
		log.Error("new kafka consumer err %v", err)
		return nil, nil, err
	}
	cancel := func() {
		if err := receiver.Close(); err != nil {
			log.Error("Kafka Receiver Close err %v", err)
		}
	}
	log.Info("new kafka broker consumer sucessful")

	return receiver, cancel, nil
}
