package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hyperbolicresearch/hlog/config"
	"github.com/hyperbolicresearch/hlog/internal/core"
	kafka_service "github.com/hyperbolicresearch/hlog/internal/kafka"
	"github.com/hyperbolicresearch/hlog/pkg/logger"
)

func LiveTail(config *config.Livetail) {
	kw, err := kafka_service.NewKafkaWorker(&config.KafkaConfigs)
	if err != nil {
		panic(err)
	}
	kw.ConfigureConsumer()
	kw.SubscribeTopics(config.KafkaTopics)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	logger := logger.New(config.DefaultLevel, os.Stdout)

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
			// log.Printf("Topic=%-10v Message=%v", *ev.TopicPartition.Topic, string(ev.Value))
			var l core.Log
			if err := json.Unmarshal(ev.Value, &l); err != nil {
				fmt.Printf("error unmarshalling value %v: %v", ev.Value, err)
			} else {
				logger.Log(l)
			}
		}
	}

}
