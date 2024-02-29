package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"github.com/hyperbolicresearch/hlog/internal/ingest"
	kafka_service "github.com/hyperbolicresearch/hlog/internal/kafka"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}
}

func main() {
	log.Println("Hlog engine started...")

	channels := os.Getenv("CHANNELS")
	clientId := os.Getenv("CLIENT_ID")
	kafkaServer := os.Getenv("KAFKA_SERVER")
	topics := strings.Split(channels, ",")
	mongodbUri := os.Getenv("MONGODB_URI")

	stop := make(chan struct{}, 1)

	mongodbConfigs := ingest.MongoDBIngesterConfig{
		KafkaConfigs: kafka_service.KafkaConfigs{
			Server:          kafkaServer,
			GroupId:         clientId,
			AutoOffsetReset: "earliest",
			EnableAutoCommit: true,
		},
		KafkaTopics:     topics,
		ConsumeInterval: time.Duration(100) * time.Millisecond,
		MongoServer:     mongodbUri,
		TopicCallback:   "hyperclusters-1415-livetail",
		Database:        clientId,
	}
	MDBIngester := ingest.NewMongoDBIngester(&mongodbConfigs)
	go MDBIngester.Start()

	<-stop
}
