package utils

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	randomstring "github.com/xyproto/randomstring"
)

type Log struct {
	LogId     string                 `json:"log_id"`
	SenderId  string                 `json:"sender_id"`
	Timestamp int64                  `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data"`
}

// GenerateRandomLogs generates every in interval seconds
// logs in a choice of numTopics topics, simulating how many processes
// would produce logs in ra real-life scenario.
func GenerateRandomLogs(stop chan struct{}) {
	quit := make(chan struct{})
	ticker := time.NewTicker(time.Second * time.Duration(10))

	// Kafka
	kafkaConfigs := kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_SERVER"),
	}
	producer, err := kafka.NewProducer(&kafkaConfigs)
	if err != nil {
		panic(err)
	}

	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Error producing: %v", err)
				} else {
					log.Printf("Produced to topic=%-9v partition=%v",
						*ev.TopicPartition.Topic,
						ev.TopicPartition.Partition)
				}
			}
		}
	}()

	run := true
	for run {
		select {
		case <-quit:
			ticker.Stop()
			stop <- struct{}{}
			run = false
		case <-ticker.C:
			go Generate(producer)
		}
	}
}

// Generate generates a new log with random data and produces it
// to a random topic via the kafka producer provided to it.
func Generate(kafkaProducer *kafka.Producer) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := rnd.Intn(5)
	channels := os.Getenv("CHANNELS")
	topics := strings.Split(channels, ",")

	id := uuid.New().String()

	senderIds := []string{
		"client-0001",
		"client-0002",
		"client-0003",
		"client-0004",
		"client-0005",
	}
	senderId := senderIds[index]

	timestamp := time.Now().UnixNano()

	levels := []string{
		"debug",
		"info",
		"warn",
		"error",
		"fatal",
	}
	level := levels[index]

	message := randomstring.HumanFriendlyEnglishString(7)

	data := map[string]interface{}{
		"foo":   "foo",
		"bar":   "bar",
		"count": index,
	}

	sendableLog := Log{
		LogId:     id,
		SenderId:  senderId,
		Timestamp: timestamp,
		Level:     level,
		Message:   message,
		Data:      data,
	}
	value, err := json.Marshal(sendableLog)
	if err != nil {
		panic(err)
	}

	kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topics[index],
			Partition: kafka.PartitionAny},
		Value: []byte(value),
	}, nil)
}
