package ingest

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/hyperbolicresearch/hlog/internal/core"
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
	// BufferSchemas stores the different schemas (one schema per channel)
	//from the buffered messages.
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
	// given that MinCommitCount is met.
	MaxBatchableWait time.Duration
}

type Messages struct {
	sync.RWMutex
	Data           []*kafka.Message
	TransformedData []interface{}
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
	for _, msg := range i.Messages.Data {
		// metadata fields
		t := make(map[string]interface{})
		var entry core.Log
		if err := json.Unmarshal(msg.Value, &entry); err != nil {
			return fmt.Errorf("error unmarshalling value %v: %v", msg.Value, err)
		}
		values := reflect.ValueOf(entry)
		types := values.Type()
		for j := 0; j < values.NumField(); j++ {
			_type := fmt.Sprintf("_%s", strings.ToLower(types.Field(j).Name))
			_value := values.Field(j)
			t[_type] = _value
		}
		// data fields of the Log.Data field.
		// PS: Needed fast fast queries, leveraging materialized views of
		// ClickHouse. They are not intended to be normal fields since we do not
		// actually want to store each field in a separate column (which can
		// be devastating as fields are dynamic)
		dataValues := reflect.ValueOf(entry.Data)
		dataTypes := values.Type()
		for k := 0; k < dataValues.NumField(); k++ {
			_type := dataTypes.Field(k).Name
			_value := dataValues.Field(k)
			t[_type] = _value
		}
		// arrays of same-typed data fields.
		// PS: Needed for queries
		// TODO: Support more types
		for k, v := range entry.Data {
			switch reflect.TypeOf(v).Kind() {
			case reflect.String:
				t["string.keys"] = append(t["string.keys"].([]string), k)
				t["string.values"] = append(t["string.values"].([]string), v.(string))
			case reflect.Int:
				t["int.keys"] = append(t["int.keys"].([]string), k)
				t["int.values"] = append(t["int.values"].([]int), v.(int))
			case reflect.Float64:
				t["float64.keys"] = append(t["float64.keys"].([]string), k)
				t["float64.values"] = append(t["float64.values"].([]float64), v.(float64))
			}
		}

		i.Messages.TransformedData = append(i.Messages.TransformedData, t)
	}

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
