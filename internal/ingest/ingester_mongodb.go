package ingest

import (
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	kafka_service "github.com/hyperbolicresearch/hlog/internal/kafka"
	"github.com/hyperbolicresearch/hlog/internal/mongodb"
)

type MongoDBIngester struct {
	*sync.RWMutex
	*mongo.Database
	*kafka_service.KafkaWorker
	ConsumeInterval time.Duration
	// TopicsCallbackList is a list of topics to produce to upon
	// successful insertion.
	TopicsCallbackList []string
}

type MongoDBIngesterConfig struct {
	// Kafka related configurations
	KafkaServer     string
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
	}
	return m
}

// Start spins up everythign and starts listening for incoming
// events from Kafka, and gets ready to sink them to the database.
func (m *MongoDBIngester) Start() {

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
