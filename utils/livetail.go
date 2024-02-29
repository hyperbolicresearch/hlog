package utils

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	kafka_service "github.com/hyperbolicresearch/hlog/internal/kafka"
)

func LiveTail() {
	topics := []string{"hyperclusters-1415-livetail"}
	kafkaServer := os.Getenv("KAFKA_SERVER")

	kafkaConfigs := kafka_service.KafkaConfigs{
		Server:          kafkaServer,
		GroupId:         "experimental-livetail",
		AutoOffsetReset: "earliest",
		EnableAutoCommit: true,
	}
	kw, err := kafka_service.NewKafkaWorker(&kafkaConfigs)
	if err != nil {
		panic(err)
	}
	kw.ConfigureConsumer()
	kw.SubscribeTopics(topics)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case <-sigchan:
			log.Printf("Caught signal: %v", sigchan)
			run = false
		default:
			ev, err := kw.Consumer.ReadMessage(time.Duration(100) * time.Millisecond)
			if err != nil {
				continue
			}
			log.Printf("Topic=%-10v Message=%v", *ev.TopicPartition.Topic, string(ev.Value))
		}
	}

}
