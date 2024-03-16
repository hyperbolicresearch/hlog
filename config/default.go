package config

import (
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/hyperbolicresearch/hlog/pkg/logger"
)

var (
	// DefaultConfig is the default top-level configuration for the
	// whole system.
	DefaultConfig = Config{
		Kafka:      &DefaultKafkaConfig,
		MongoDB:    &DefaultMongoDBConfig,
		ClickHouse: &DefaultClickHouseConfig,
		Livetail:   &DefaultLivetailConfig,
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
		KafkaConfigs:     DefaultKafkaConfig,
		KafkaTopics:      []string{"default"},
		ConsumeInterval:  time.Duration(100) * time.Millisecond,
		MinBatchableSize: 1,
		MaxBatchableSize: 1000,
		MaxBatchableWait: time.Duration(30) * time.Second,
	}

	// DefaultLivetailConfig is the default Livetail configuration.
	DefaultLivetailConfig = Livetail{
		KafkaTopics: []string{"default"},
		KafkaConfigs: Kafka{
			Server:           "0.0.0.0:65007",
			GroupId:          "hlog-livetail-default",
			AutoOffsetReset:  "earliest",
			EnableAutoCommit: true,
		},
		ConsumeInterval: time.Duration(100) * time.Millisecond,
		DefaultLevel:    logger.DEBUG,
	}
)
