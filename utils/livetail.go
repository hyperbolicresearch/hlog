package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hyperbolicresearch/hlog/config"
	"github.com/hyperbolicresearch/hlog/internal/logs"
	kafkaservice "github.com/hyperbolicresearch/hlog/transport/kafka"
)

// LiveTail is a real-time, bridge between Kafka and a logging medium
// that allows the observation as they are occuring of newly ingested
// messages.
func LiveTail(config *config.Livetail, sigchan chan os.Signal) {
	kw, err := kafkaservice.NewKafkaWorker(&config.KafkaConfigs)
	if err != nil {
		panic(err)
	}
	kw.ConfigureConsumer()
	kw.SubscribeTopics(config.KafkaTopics)
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
			var l logs.Log
			if err := json.Unmarshal(ev.Value, &l); err != nil {
				fmt.Printf("error unmarshalling value %v: %v", ev.Value, err)
			} else {
				err := config.Logger.Log(l)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
