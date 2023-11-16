package kafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/segmentio/kafka-go"
)

var (
	_ Sender   = (*KafkaSender)(nil)
	_ Receiver = (*kafkaReceiver)(nil)
	_ Event    = (*Message)(nil)
)

type Message struct {
	key   string
	value []byte
}

func (m *Message) Key() string {
	return fmt.Sprintf("%v:%v", uuid.New().String(), m.key)
}

func (m *Message) MsgKey() string {
	s := strings.Split(m.key, ":")
	if len(s) > 0 {
		return s[len(s)-1]
	}
	return ""
}

func (m *Message) Value() []byte {
	return m.value
}

func NewMessage(key string, value []byte) Event {
	return &Message{
		key:   key,
		value: value,
	}
}

type KafkaSender struct {
	writer *kafka.Writer
	topic  string
}

func (s *KafkaSender) Send(message Event) error {
	err := s.writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(message.Key()),
		Value: message.Value(),
	})
	if err != nil {
		log.Errorf("Send Err %v", err)
		return err
	}
	return nil
}

func (s *KafkaSender) Close() error {
	err := s.writer.Close()
	if err != nil {
		return err
	}
	return nil
}

func NewKafkaSender(address []string, topic string) (Sender, error) {
	w := &kafka.Writer{
		Topic: topic,
		Addr:  kafka.TCP(address...),
		//Balancer:     &kafka.LeastBytes{},
		Balancer:               &kafka.Hash{},
		BatchTimeout:           time.Millisecond * 100,
		AllowAutoTopicCreation: true,
	}
	return &KafkaSender{writer: w, topic: topic}, nil
}

type kafkaReceiver struct {
	reader  *kafka.Reader
	topic   string
	groupId string
}

func (k *kafkaReceiver) Receive(ctx context.Context, handler Handler) error {
	go func() {
		if k.groupId == "" {
			for {
				wctx := context.Background()
				m, err := k.reader.ReadMessage(wctx)
				if err != nil {
					if errors.Is(err, io.EOF) {
						log.Debugf("Kafka Receive Topic %v Close", k.topic)
						break
					}

					log.Errorf("Kafka Receive Topic %v FetchMessage %v", k.topic, err)
					continue
				}
				err = handler(ctx, &Message{
					key:   string(m.Key),
					value: m.Value,
				})
				if err != nil {
					log.Errorf("message handling exception: %v", err)
				}
			}
		} else {
			for {
				wctx := context.Background()
				m, err := k.reader.FetchMessage(wctx)
				if err != nil {
					if errors.Is(err, io.EOF) {
						log.Debugf("Kafka Receive Topic %v Close", k.topic)
						break
					}

					log.Errorf("Kafka Receive Topic %v FetchMessage %v", k.topic, err)
					continue
				}
				err = handler(ctx, &Message{
					key:   string(m.Key),
					value: m.Value,
				})
				if err != nil {
					log.Errorf("message handling exception: %v", err)
				}
				if err := k.reader.CommitMessages(wctx, m); err != nil {
					log.Errorf("failed to commit messages: %v", err)
				}
			}
		}
	}()
	return nil
}

func (k *kafkaReceiver) Close() error {
	err := k.reader.Close()
	if err != nil {
		return err
	}
	return nil
}

func NewKafkaReceiver(address []string, topic string, groupId string) (Receiver, error) {
	cfg := kafka.ReaderConfig{
		Brokers:  address,
		GroupID:  groupId,
		Topic:    topic,
		MinBytes: 1,    // 10KB
		MaxBytes: 10e6, // 10MB
		MaxWait:  500 * time.Millisecond,
	}
	if groupId != "" {
		cfg.StartOffset = kafka.LastOffset
		cfg.CommitInterval = time.Second
	}
	r := kafka.NewReader(cfg)
	if groupId == "" {
		err := r.SetOffsetAt(context.Background(), time.Now())
		if err != nil {
			log.Error("NewKafkaReceiver SetOffsetAt Error: %v", err)
		}
	}

	return &kafkaReceiver{reader: r, topic: topic, groupId: groupId}, nil
}
