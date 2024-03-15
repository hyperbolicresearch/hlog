package config

import (
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

var (
	// DefaultConfig is the default top-level configuration for the
	// whole system.
	DefaultConfig = Config{
		Kafka:      &DefaultKafkaConfig,
		MongoDB:    &DefaultMongoDBConfig,
		ClickHouse: &DefaultClickHouseConfig,
	}

	// DefaultKafkaConfig is the default kafka configuration.
	DefaultKafkaConfig = Kafka{
		Server:           "0.0.0.0:65007",
		GroupId:          "hlog-default",
		AutoOffsetReset:  "earliest",
		EnableAutoCommit: false,
	}

	// DefaultMongoDBConfig is the default MongoDB configuration.
	DefaultMongoDBConfig = MongoDB{
		Server:          "mongodb://localhost:27017/",
		Database:        "hlog-default",
		TopicCallback:   "hlog-mongodb-callback",
		KafkaConfigs:    DefaultKafkaConfig,
		KafkaTopics:     []string{"default"},
		ConsumeInterval: time.Millisecond * time.Duration(100),
	}

	// DefaultClickHouseConfig is the default ClickHouse configuration.
	DefaultClickHouseConfig = ClickHouse{
		Options: &clickhouse.Options{
			Addr: []string{"localhost:9000"},
			Auth: clickhouse.Auth{
				Username: "default",
				Database: "default",
			},
		},
	}
)
