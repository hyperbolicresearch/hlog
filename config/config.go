package config

import (
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"

	"github.com/hyperbolicresearch/hlog/pkg/logger"
)

// Config is the top-level configuration for hlog
type Config struct {
	*Kafka
	*MongoDB
	*ClickHouse
	*Livetail
	*Simulator
	*APIv1
}

// Kafka holds the configuration for Kafka
type Kafka struct {
	Server           string
	GroupId          string
	AutoOffsetReset  string
	EnableAutoCommit bool
}

// MongoDB holds the configuration for MongoDB
type MongoDB struct {
	Server   string
	Database string
	// Kafka-related
	TopicCallback   string
	KafkaConfigs    Kafka
	KafkaTopics     []string
	ConsumeInterval time.Duration
}

// ClickHouse holds the configuration for ClickHouse
type ClickHouse struct {
	*clickhouse.Options

	KafkaConfigs     Kafka
	KafkaTopics      []string
	ConsumeInterval  time.Duration
	MinBatchableSize int
	MaxBatchableSize int
	MaxBatchableWait time.Duration
}

// FromYAML reads configs.yaml and extracts the configurations
func FromYAML(filename string) (*Config, error) {
	return nil, fmt.Errorf("config file %v not found", filename)
}

// Livetail holds the configuration for the terminal-base
// visualization of the entering logs
type Livetail struct {
	KafkaTopics             []string
	KafkaConfigs            Kafka
	ConsumeInterval         time.Duration
	DefaultLevel            logger.Level
	InitLogsLoadedCount     int
	Logger                  *logger.Logger
	MaxWebsocketConnections int
	WebsocketPort           int
}

// Simulator holds the configurations for the log producing simulator
type Simulator struct {
	KafkaTopics     []string
	KafkaConfigs    Kafka
	ProduceInterval time.Duration
	MessageLength   int
}

// API contains the configurations for the API
type APIv1 struct {
	// ServerAddr is the address on which the API is listening
	ServerAddr string
	
	
	// Livetail configurations ___________________________________________
	
	// LivetailLogger writes the livetail logs
	LivetailLogger *logger.Logger
	// MaxLiveTailWebsocketConnections is the number of concurrent readers
	// of v1.LiveTail
	MaxLiveTailWebsocketConnections int
	KafkaTopics                     []string
	KafkaConfigs                    Kafka
	ConsumeInterval                 time.Duration
	DefaultLevel                    logger.Level
	// InitLogsLoadedCount is the number of logs to load on the fly before
	// starting the livetailing process
	InitLogsLoadedCount int
	
	
	// ObservablesTail configurations ____________________________________
	
	// GeneralObservablesLogger writes the general observable metrics
	GeneralObservablesLogger *logger.Logger
	// MaxGenObsWebsocketConnections is the number of concurrent
	// readers of v1.GeneralObservables
	MaxGenObsWebsocketConnections int
	// PushInterval determines how often to push updated versions of
	// the GeneralObservables
	PushInterval time.Duration
	// SendGeneralObservables determines whether or not to send the
	// v1.GeneralObservables
	SendGeneralObservables bool
	// SendGeneralSystemObservables determines whether or not to send the
	// v1.CustomObservables
	SendGeneralSystemObservables bool
}
