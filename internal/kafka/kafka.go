package kafka_service

import (
	"fmt"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/hyperbolicresearch/hlog/config"
)

type PubSubWorker interface {
	ConfigureConsumer() error
	ConfigureProducer() error
	ConfigurePubSub() error
	SubscribeTopics(topics []string) error
}

type KafkaWorker struct {
	sync.RWMutex
	*kafka.Consumer
	*kafka.Producer
	Configs    *config.Kafka
	IsConsumer bool
	IsProducer bool
}

type KafkaConfigs struct {
	Server           string
	GroupId          string
	AutoOffsetReset  string
	EnableAutoCommit bool
}

// NewKafkaWorker creates a new KafkaWorker and returns it
// with an error message.
func NewKafkaWorker(cfg *config.Kafka) (*KafkaWorker, error) {
	w := &KafkaWorker{
		Configs: cfg,
	}
	return w, nil
}

func (k *KafkaWorker) ConfigurePubSub() error {
	err := k.ConfigureConsumer()
	if err != nil {
		return err
	}
	err = k.ConfigureProducer()
	if err != nil {
		return err
	}
	return nil
}

// ConfigureConsumer makes the KafkaWorker a consumer
func (k *KafkaWorker) ConfigureConsumer() error {
	cfg := kafka.ConfigMap{
		"bootstrap.servers":  k.Configs.Server,
		"group.id":           k.Configs.GroupId,
		"auto.offset.reset":  k.Configs.AutoOffsetReset,
		"enable.auto.commit": k.Configs.EnableAutoCommit,
	}
	consumer, err := kafka.NewConsumer(&cfg)
	if err != nil {
		return fmt.Errorf("failed to create consumer: %v", err)
	}
	k.Consumer = consumer
	k.Lock()
	k.IsConsumer = true
	k.Unlock()
	return nil
}

// ConfigureProducer makes the KafkaWorker a producer
func (k *KafkaWorker) ConfigureProducer() error {
	cfg := kafka.ConfigMap{
		"bootstrap.servers": k.Configs.Server,
	}
	producer, err := kafka.NewProducer(&cfg)
	if err != nil {
		return fmt.Errorf("failed to create producer: %v", err)
	}
	k.Producer = producer
	k.Lock()
	k.IsProducer = true
	k.Unlock()
	return nil
}

// SubscribeTopics subscribes to a given list of topics for consuming
func (k *KafkaWorker) SubscribeTopics(topics []string) error {
	err := k.Consumer.SubscribeTopics(topics, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topics: %v, error: %v",
			topics, err)
	}
	return nil
}
