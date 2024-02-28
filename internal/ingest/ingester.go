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
	Start()
	Stop() error
	Consume() error
	Transform() error
	ExtractSchemas() error
	Sink() error
	Commit() error
}

type IngesterWorker struct {
	sync.RWMutex
	*kafka.Consumer
	IsRunning bool
	Messages  *Messages
	// BufferSchemas stores the different schemas (one schema per
	// channel) from the buffered messages.
	BufferSchemas []interface{}
	// CanCommit is a way for us to make sure that we have successfully
	// sink the buffered messages to ClickHouse, and that we can commit
	// to Kafka. This guarantess At-Least-Once delivery.
	CanCommit chan struct{}
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
	i.Lock()
	defer i.Unlock()
	i.IsRunning = false
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
	_ = i.ExtractSchemas()
	go i.Sink()

	// Waiting until we have the acks that we successfully sink
	// the data to ClickHouse (which should be sent from inside
	// Sink)
	<-i.CanCommit
	i.Commit()

	return nil
}

// Transform will flatten the message to the appropriate format
// that will be stored to ClickHouse, add metadata.
func (i *IngesterWorker) Transform() error {
	// 1. assume interface{}
	// 2. assert that Log.Data is of type map[string]interface
	// 3. get all the keys to a slice
	// 4. iterate over it and add corresponding keys with types to new map
	// 5. add metadata to the new map
	// 6. create struct from map (? is it necessary)

	return nil
}

func (i *IngesterWorker) ExtractSchemas() error {
	// 1. create a slice of channels from messages
	// 2. for each channel, generate a schema with ClickHouse
	return nil
}

func (i *IngesterWorker) Sink() error {
	// 1. alter table if needed for each slice
	// 2. sink the data to clickhouse
	// 3. write to CanCommit channel
	return nil
}

func (i *IngesterWorker) Commit() error {
	// 1. commit to current offset in kafka
	// 2. log about batch processing completion
	return nil
}
