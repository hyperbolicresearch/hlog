package ingest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/hyperbolicresearch/hlog/config"
	"github.com/hyperbolicresearch/hlog/internal/core"
	kafka_service "github.com/hyperbolicresearch/hlog/internal/kafka"
	"github.com/hyperbolicresearch/hlog/internal/mongodb"
)

type MongoDBIngester struct {
	sync.RWMutex
	*mongo.Database
	*kafka_service.KafkaWorker
	ConsumeInterval time.Duration
	// TopicCallback is the Kafka topic to produce to upon successful insertion.
	TopicCallback string
	CloseChan     chan struct{}
}

type MongoDBIngesterConfig struct {
	// Kafka-related configurations
	KafkaConfigs    kafka_service.KafkaConfigs
	KafkaTopics     []string
	ConsumeInterval time.Duration
	// MongoDB-related configurations
	MongoServer   string
	TopicCallback string
	Database      string
}

func NewMongoDBIngester(cfg *config.Config) *MongoDBIngester {
	mongoClient := mongodb.Client(cfg.MongoDB.Server)
	db := mongoClient.Database(cfg.Database)

	kw, err := kafka_service.NewKafkaWorker(cfg.Kafka)
	if err != nil {
		panic(err)
	}
	err = kw.ConfigurePubSub()
	if err != nil {
		panic(err)
	}
	err = kw.SubscribeTopics(cfg.MongoDB.KafkaTopics)
	if err != nil {
		panic(err)
	}
	m := &MongoDBIngester{
		ConsumeInterval: cfg.MongoDB.ConsumeInterval,
		Database:        db,
		KafkaWorker:     kw,
		TopicCallback:   cfg.MongoDB.TopicCallback,
		CloseChan:       make(chan struct{}, 1),
	}
	return m
}

// Start spins up everything and starts listening for incoming
// events from Kafka, and gets ready to sink them to the database.
func (m *MongoDBIngester) Start(stop chan os.Signal) {
	m.RLock()
	kw := m.KafkaWorker
	m.RUnlock()

	// we can only start a if it's a consumer
	if !kw.IsConsumer {
		return
	}

	run := true
	for run {
		select {
		case <-m.CloseChan:
			run = false
		default:
			m.Consume()
		}
	}
}

func (m *MongoDBIngester) Stop() error {
	m.CloseChan <- struct{}{}
	return nil
}

func (m *MongoDBIngester) Consume() error {
	m.RLock()
	ci := m.ConsumeInterval
	m.RUnlock()
	ev, err := m.KafkaWorker.Consumer.ReadMessage(ci)
	if err != nil {
		return nil
	}
	go m.Sink(ev)

	return err
}

func (m *MongoDBIngester) Sink(msg *kafka.Message) error {
	var value core.Log
	if err := json.Unmarshal(msg.Value, &value); err != nil {
		fmt.Printf("Error unmarshalling value %v", err)
	}

	m.Lock()
	col := m.Database.Collection(*msg.TopicPartition.Topic)
	defer m.Unlock()
	_, err := col.InsertOne(context.TODO(), value)
	if err != nil {
		panic(err)
	}

	// Produce to m.TopicCallback if any.
	// TODO : Probably export to a separate function ???
	if m.TopicCallback != "" {
		m.KafkaWorker.Lock()
		m.KafkaWorker.Producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &m.TopicCallback,
				Partition: kafka.PartitionAny},
			Value: []byte(msg.Value),
		}, nil)
		m.KafkaWorker.Unlock()
	}

	log.Printf("Successfully processed log from topic: %-10v Message: %v",
		*msg.TopicPartition.Topic, value.Message)

	return nil
}
