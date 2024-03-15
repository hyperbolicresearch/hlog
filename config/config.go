package config

import (
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

// Config is the top-level configuration for hlog
type Config struct {
	*Kafka
	*MongoDB
	*ClickHouse
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
}

// FromYAML reads configs.yaml and extracts the configurations
func FromYAML(filename string) (*Config, error) {
	c := Config{}
	return &c, nil
}
