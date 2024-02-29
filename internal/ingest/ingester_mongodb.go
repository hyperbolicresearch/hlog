package ingest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/hyperbolicresearch/hlog/internal/core"
	kafka_service "github.com/hyperbolicresearch/hlog/internal/kafka"
	"github.com/hyperbolicresearch/hlog/internal/mongodb"
)

type MongoDBIngester struct {
	sync.RWMutex
	*mongo.Database
	*kafka_service.KafkaWorker
	ConsumeInterval time.Duration
	// TopicsCallbackList is a list of topics to produce to upon
	// successful insertion.
	TopicsCallbackList []string
	CloseChan          chan struct{}
}

type MongoDBIngesterConfig struct {
	// Kafka related configurations
	KafkaConfigs    kafka_service.KafkaConfigs
	KafkaTopics     []string
	ConsumeInterval time.Duration
	// MongoDb related configurations
	MongoServer        string
	TopicsCallbackList []string
	Database           string
}

func NewMongoDBIngester(configs *MongoDBIngesterConfig) *MongoDBIngester {
	mongoClient := mongodb.Client(configs.MongoServer)
	db := mongoClient.Database(configs.Database)

	kw, err := kafka_service.NewKafkaWorker(&configs.KafkaConfigs)
	if err != nil {
		panic(err)
	}
	err = kw.ConfigurePubSub()
	if err != nil {
		panic(err)
	}
	m := &MongoDBIngester{
		ConsumeInterval:    configs.ConsumeInterval,
		Database:           db,
		KafkaWorker:        kw,
		TopicsCallbackList: configs.TopicsCallbackList,
		CloseChan:          make(chan struct{}, 1),
	}
	return m
}

// Start spins up everything and starts listening for incoming
// events from Kafka, and gets ready to sink them to the database.
func (m *MongoDBIngester) Start() {
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
	
	col := m.Database.Collection(*msg.TopicPartition.Topic)
	_, err := col.InsertOne(context.TODO(), value)
	if err != nil {
		panic(err)
	}
	
	log.Printf("Successfully processed log from topic: %v", 
		*msg.TopicPartition.Topic)
	return nil
}
