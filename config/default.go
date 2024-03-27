package config

import (
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/hyperbolicresearch/hlog/pkg/logger"
)

var (
	// DefaultConfig is the default top-level configuration for the whole system.
	DefaultConfig = Config{
		Kafka:      &DefaultKafkaConfig,
		MongoDB:    &DefaultMongoDBConfig,
		ClickHouse: &DefaultClickHouseConfig,
		Livetail:   &DefaultLivetailConfig,
		Simulator:  &DefaultSimulatorConfig,
		APIv1:      &DefaultAPIConfig,
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
		Server:        "mongodb://localhost:27017/",
		Database:      "hlog-default",
		TopicCallback: "hlog-mongodb-callback",
		KafkaConfigs: Kafka{
			Server:           "0.0.0.0:65007",
			GroupId:          "hlog-default-mongodb",
			AutoOffsetReset:  "earliest",
			EnableAutoCommit: true,
		},
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
		KafkaConfigs: Kafka{
			Server:           "0.0.0.0:65007",
			GroupId:          "hlog-default-clickhouse",
			AutoOffsetReset:  "earliest",
			EnableAutoCommit: false,
		},
		KafkaTopics:      []string{"default"},
		ConsumeInterval:  time.Duration(5) * time.Second,
		MinBatchableSize: 1,
		MaxBatchableSize: 100,
		MaxBatchableWait: time.Duration(10) * time.Second,
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
		ConsumeInterval:         time.Duration(100) * time.Millisecond,
		DefaultLevel:            logger.DEBUG,
		InitLogsLoadedCount:     100,
		Logger:                  logger.New(logger.DEBUG, os.Stdout),
		MaxWebsocketConnections: 5,
		WebsocketPort:           1337,
	}

	// DefaultSimulatorConfig is the default Simulator configuration.
	DefaultSimulatorConfig = Simulator{
		KafkaTopics: []string{"default"},
		KafkaConfigs: Kafka{
			Server: "0.0.0.0:65007",
		},
		ProduceInterval: time.Duration(5) * time.Second,
		MessageLength:   7,
	}

	DefaultAPIConfig = APIv1{
		ServerAddr:                      "localhost:1542",
		LivetailLogger:                  logger.New(logger.DEBUG, os.Stdout),
		MaxLiveTailWebsocketConnections: 10,
		KafkaTopics:                     []string{"default"},
		KafkaConfigs: Kafka{
			Server:           "0.0.0.0:65007",
			GroupId:          "hlog-livetail-default",
			AutoOffsetReset:  "earliest",
			EnableAutoCommit: true,
		},
		ConsumeInterval:               time.Duration(100) * time.Millisecond,
		DefaultLevel:                  logger.DEBUG,
		InitLogsLoadedCount:           100,
		GeneralObservablesLogger:      logger.New(logger.DEBUG, os.Stdout),
		MaxGenObsWebsocketConnections: 10,
		PushInterval:                  time.Duration(5) * time.Second,
		SendGeneralObservables:        true,
		SendGeneralSystemObservables:  true,
	}
)
