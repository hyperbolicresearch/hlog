package ingest

import (
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Ingester interface {
	// Start will start the process of ingesting.
	Start() error

	// Stop stops the process of ingesting.
	Stop() error

	// Consume will read from Kafka and store the messages
	Consume() error

	// SendForBatching makes the stored messages available
	// for insertion to ClickHouse by the Batcher
	SendForBatching() error

	// ExtractSchema will extact the schema from the buffered
	// messages.
	ExtractSchema() error

	// Transform will perform the transformations of the messages
	// in the format we store them in ClickHouse
	Transform() error

	// Commit will commit to Kafka after receiving the acks from
	// the Batcher for the corresponding offset.
	Commit() error
}

type IngesterWorker struct {
	sync.RWMutex
	*kafka.Consumer

	Messages     []Messages
	IsRunning    bool
	BufferSchema interface{}
	// ConsumeInterval is the periodic interval to consume messages
	// from kafka.
	ConsumeInterval time.Duration
	// MinCommitCount is the minimum number of messages that should
	// be stored to allow committing or batching.
	MinCommitCount int
	// MaxCommitCount is the threshold for committing and batching
	// givent that MinCommitCount is met.
	MaxCommitWait time.Duration
}

type Messages struct {
	sync.RWMutex
	Data []*kafka.Message
}
