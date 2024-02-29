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
			m.RLock()
			ci := m.ConsumeInterval
			m.RUnlock()
			ev, err := kw.Consumer.ReadMessage(ci)
			if err != nil {
				continue
			}
			go m.processLog(ev, m.Database)
		}
	}
}

func (m *MongoDBIngester) Stop() error {
	return nil
}

func (m *MongoDBIngester) Consume() error {
	return nil
}

func (m *MongoDBIngester) Sink() error {
	return nil
}

// processLog is a helper function that will take a new message,
// and add it to the MongoDB database.
func (m *MongoDBIngester) processLog(ev *kafka.Message, db *mongo.Database) {
	var value core.Log
	if err := json.Unmarshal(ev.Value, &value); err != nil {
		fmt.Printf("Error unmarshalling value %v", err)
	}

	col := db.Collection(*ev.TopicPartition.Topic)
	_, err := col.InsertOne(context.TODO(), value)
	if err != nil {
		panic(err)
	}

	log.Printf("Successfully processed log from topic: %v", *ev.TopicPartition.Topic)
}
