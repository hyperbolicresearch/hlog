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

	"github.com/hyperbolicresearch/hlog/config"
	"github.com/hyperbolicresearch/hlog/internal/core"
)

// GenerateRandomLogs generates logs every in k seconds
// in a choice of numTopics topics, simulating how many processes
// would produce logs in ra real-life scenario.
func GenerateRandomLogs(cfg *config.Config, stop chan os.Signal) {
	kafkaConfigs := kafka.ConfigMap{
		"bootstrap.servers": cfg.Simulator.KafkaConfigs.Server,
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

	ticker := time.NewTicker(time.Second * time.Duration(2))
	run := true
	for run {
		select {
		case <-stop:
			ticker.Stop()
			run = false
		case <-ticker.C:
			rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
			amount := rnd.Intn(50)
			for range amount {
				go Generate(producer, cfg)
			}
		}
	}
}

// Generate generates a new log with random data and produces it
// to a random topic via the kafka producer provided to it.
func Generate(kafkaProducer *kafka.Producer, cfg *config.Config) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := rnd.Intn(len(cfg.Simulator.KafkaTopics))
	topics := cfg.Simulator.KafkaTopics
	channel := topics[index]

	id := uuid.New().String()

	index = rnd.Intn(5)
	senderIds := []string{
		"client-0001",
		"client-0002",
		"client-0003",
		"client-0004",
		"client-0005",
	}
	senderId := senderIds[index]

	timestamp := time.Now().Unix()

	index = rnd.Intn(5)
	levels := []string{
		"debug",
		"info",
		"warn",
		"error",
		"fatal",
	}
	level := levels[index]

	message := ""
	foo := ""
	bar := ""
	for i := 0; i < cfg.Simulator.MessageLength; i++ {
		word := randomstring.HumanFriendlyEnglishString(7) + " "
		wfoo := randomstring.HumanFriendlyEnglishString(7) + " "
		wbar := randomstring.HumanFriendlyEnglishString(7) + " "
		message += word
		foo += wfoo
		bar += wbar
	}
	message = strings.TrimSpace(message)

	data := map[string]interface{}{
		"foo":   foo,
		"bar":   bar,
		"count": index,
		"firstname": "John",
		"lastname": "Doe",
		"company": "Acme Inc.",
	}

	sendableLog := core.Log{
		Channel:   channel,
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
			Topic:     &channel,
			Partition: kafka.PartitionAny},
		Value: []byte(value),
	}, nil)
}
