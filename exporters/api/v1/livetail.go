package v1

import (
	"encoding/json"
	"os"

	"github.com/hyperbolicresearch/hlog/config"
	"github.com/hyperbolicresearch/hlog/internal/logs"
	kafkaservice "github.com/hyperbolicresearch/hlog/transport/kafka"
)

// LiveTail is the data structure that manages the process of pushing in
// real-time the incoming new entries in the system to the Observer.
type LiveTail struct {
	config      *config.APIv1
	kafkaWorker *kafkaservice.KafkaWorker
}

// NewLiveTail creates a new LiveTail instance
func NewLiveTail(cfg *config.APIv1) *LiveTail {
	// NOTE: Since the system is ingesting data only through Kafka,
	// there is no need to make this more generic. The live tailing
	// process if done directly from Kafka, without confirmation
	// of ingestion into the database. That's why we are configuring
	// Kafka here.
	kw, err := kafkaservice.NewKafkaWorker(&cfg.KafkaConfigs)
	if err != nil {
		panic(err)
	}
	kw.ConfigureConsumer()
	kw.SubscribeTopics(cfg.KafkaTopics)

	return &LiveTail{
		config:      cfg,
		kafkaWorker: kw,
	}
}

// Start will start listening for new entries and pushing them through
// the list of listeners at l.Logger.Writers
func (l *LiveTail) Start(sig chan os.Signal) error {
	run := true
	for run {
		select {
		case <-sig:
			run = false
		default:
			msg, err := l.kafkaWorker.Consumer.ReadMessage(l.config.ConsumeInterval)
			if err != nil {
				continue
			}
			var log logs.Log
			if err := json.Unmarshal(msg.Value, &log); err != nil {
				// TODO: handle error
				return err
			}
			// The actual process is happening here. Remark how this leverages
			// the logger's ability to write to io.Writer to actually write
			// to the Websocket connections.
			// NOTE: The websocket connections are actually added to the
			// logger by the API server's endpoints.
			err = l.config.LivetailLogger.Log(log)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
