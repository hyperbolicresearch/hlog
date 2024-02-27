package ingest

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Ingester interface {
	// Start will start the process of ingesting.
	Start()
	// Stop stops the process of ingesting.
	Stop() error
	// Consume will read from Kafka and store the messages
	Consume() error
	// Transform will perform the transformations of the messages
	// in the format we store them in ClickHouse
	Transform() error
	// ExtractSchema will extact the schema from the buffered
	// messages.
	ExtractSchema() error
	// SendForBatching makes the stored messages available
	// for insertion to ClickHouse by the Batcher
	SendForBatching() error
	// Commit will commit to Kafka after receiving the acks from
	// the Batcher for the corresponding offset.
	Commit() error
}

type IngesterWorker struct {
	sync.RWMutex
	*kafka.Consumer

	IsRunning    bool
	Messages     *Messages
	BufferSchema interface{}
	// CanCommit is a way for us to make sure that we have successfully
	// sink the buffered messages to ClickHouse, and that we can commit
	// to Kafka. This guarantess At-Least-Once delivery.
	CanCommit    chan struct{}
	// ConsumeInterval is the periodic interval to consume messages
	// from kafka.
	ConsumeInterval time.Duration
	// MinBatchableSize is the minimum number of messages that should
	// be stored to allow committing or batching.
	MinBatchableSize int
	// MaxBatchableSize is the maximum number of messages that we can
	// buffer to Messages before committing or batching.
	MaxBatchableSize int
	// MaxBatchableWait is the threshold for committing and batching
	// givent that MinCommitCount is met.
	MaxBatchableWait time.Duration
}

type Messages struct {
	sync.RWMutex
	Data []*kafka.Message
}

// Start spins up the consuming process generally speaking. It runs as
// long as i.IsRunning is true and tracks
func (i *IngesterWorker) Start() {
	// TODO: Should I really panic here?
	i.RLock()
	if i.IsRunning {
		panic("Ingester worker already running...")
	}
	i.RUnlock()

	i.Lock()
	i.IsRunning = true
	i.Unlock()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	consumeTicker := time.NewTicker(i.ConsumeInterval)

	for i.IsRunning {
		select {
		case <-sigchan:
			i.Stop()
		case <-consumeTicker.C:
			go i.Consume()
		}
	}
}

// Stop will stop all the ongoing processes gracefully, including
// gracefully shutting down the consumming process from Kafka and
// the sinking process to ClickHouse.
func (i *IngesterWorker) Stop() error {
	return nil
}

// Consume reads up to i.MaxBatchableSize messages from Kafka and
// orchestrate the further processing of these by invoking the
// subsequent methods.
func (i *IngesterWorker) Consume() error {
	// We will try to extract as much as possible messages from
	// Kafka given that j <= i.MaxBatchableSize.
	// Doing like that, we make sure that we always read less
	// or equal to the i.MaxBatchableSize.
	for j := 0; j < i.MaxBatchableSize; j++ {
		msg, err := i.Consumer.ReadMessage(i.ConsumeInterval)
		if err != nil {
			continue
		}
		i.Messages.Lock()
		i.Messages.Data = append(i.Messages.Data, msg)
		i.Messages.Unlock()
	}

	// The following pipeline will perfom the preparation and the
	// sink of the buffered data to ClickHouse.
	_ = i.Transform()
	_ = i.ExtractSchema()
	go i.SendForBatching()

	// Waiting until we have the acks that we successfully sink
	// the data to ClickHouse (which should be sent from inside
	// SendForBatching)
	<-i.CanCommit
	i.Commit()

	return nil
}

func(i *IngesterWorker) Transform() error {
	return nil
}

func (i *IngesterWorker) ExtractSchema() error {
	return nil
}

func (i *IngesterWorker) SendForBatching() error {
	return nil
}

func (i *IngesterWorker) Commit() error {
	return nil
}
