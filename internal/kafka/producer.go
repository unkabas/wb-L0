package kafka

import (
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"strings"
)

const (
	flushTimout = 5000
)

var errUnknownType = errors.New("unknown event type")

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(address []string) (*Producer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(address, ","),
	}
	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, fmt.Errorf("error creating producer: %w", err)
	}
	return &Producer{producer: p}, nil
}

func (p *Producer) Produce(message, topic, key string) error {
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(message),
		Key:   []byte(key),
	}
	kafkaChan := make(chan kafka.Event)
	if err := p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
		return fmt.Errorf("error producing message: %w", err)
	}
	e := <-kafkaChan
	switch ev := e.(type) {
	case *kafka.Message:
		return nil
	case kafka.Error:
		return ev
	default:
		return errUnknownType
	}
}

func (p *Producer) Close() {
	p.producer.Flush(flushTimout)
	p.producer.Close()
}
