package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}
}

func main() {
	// 1. Configure Mongo and Kafka

	// KAFKA
	channels := os.Getenv("CHANNELS")
	clientId := os.Getenv("CLIENT_ID")
	kafkaServer := os.Getenv("KAFKA_SERVER")
	topics := strings.Split(channels, ",")
	configs := kafka.ConfigMap{
		"bootstrap.servers": kafkaServer,
		"group.id":          clientId,
		"auto.offset.reset": "earliest",
	}
	kw, err := NewKafkaWorker(&configs, topics)
	if err != nil {
		os.Exit(1)
	}

	// MONGODB
	mongodbUri := os.Getenv("MONGODB_URI")
	if mongodbUri == "" {
		panic("no set uri for mongodb")
	}
	clientOptions := options.Client().ApplyURI(mongodbUri)
	ctx := context.Background()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)
	db := client.Database(clientId)

	stop := make(chan struct{}, 1)
	go func() {
		run := true
		for run {
			select {
			case <-stop:
				run = false
			default:
				ev, err := kw.Consumer.ReadMessage(100 * time.Millisecond)
				if err != nil {
					continue
				}
				go ProcessLog(ev, db)
			}
		}
	}()
	// 5. Then they are added to ClickHouse, and all the metrics queries
	//    are runned again to update everything.
	// 6. From the dashboard, we can access all the logs by query them,
	//    access the metrics that are available.
	// 7. (future) You can add processing pipeline to process the logs
	//    which are essentially functions runned with the log as input.

	<-stop
}

// ProcessLog takes a new incomming message (log) and further process
// the storage
func ProcessLog(ev *kafka.Message, db *mongo.Database) {
	var value Log
	if err := json.Unmarshal(ev.Value, &value); err != nil {
		fmt.Printf("Error unmarshalling value %v", err)
	}
	
}
