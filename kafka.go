package main

import (
	"errors"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaWorker struct {
	sync.RWMutex
	*kafka.Consumer
}

// NewKafkaWorker creates a new KafkaWorker and returns it
// with an error message.
func NewKafkaWorker(configs *kafka.ConfigMap, topics []string) (*KafkaWorker, error) {
	consumer, err := kafka.NewConsumer(configs)
	if err != nil {
		return nil, errors.New("failed to create consumer")
	}
	err = consumer.SubscribeTopics(topics, nil)
	if err != nil {
		return nil, errors.New("failed to subscribe topics")
	}

	w := &KafkaWorker{
		Consumer: consumer,
	}

	return w, nil
}

