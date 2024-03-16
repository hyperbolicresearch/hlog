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
	KafkaTopics     []string
	KafkaConfigs    Kafka
	ConsumeInterval time.Duration
	DefaultLevel    logger.Level
}
