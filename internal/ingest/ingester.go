package ingest

import (
	"context"
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

	clickhouse_connector "github.com/hyperbolicresearch/hlog/internal/clickhouse"
	"github.com/hyperbolicresearch/hlog/internal/core"
	kafka_service "github.com/hyperbolicresearch/hlog/internal/kafka"
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
	*kafka_service.KafkaWorker
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
	Data            []*core.Log
	TransformedData []map[string]interface{}
}

// TODO: make configs
func NewIngesterWorker() *IngesterWorker {
	channels := "mainnet,testnet,subnet,intranet,darknet"
	topics := strings.Split(channels, ",")
	groupId := "hyperclusters-1415"
	kafkaServer := "0.0.0.0:65007"
	kConfigs := kafka.ConfigMap{
		"bootstrap.servers":  kafkaServer,
		"group.id":           groupId,
		"enable.auto.commit": false,
	}
	kw, err := kafka_service.NewKafkaWorker(&kConfigs, topics)
	if err != nil {
		panic("failed to create ingester")
	}

	_i := &IngesterWorker{
		Messages:         &Messages{},
		BufferSchemas: make([]interface{}, 0, len(topics)),
		KafkaWorker:      kw,
		CanCommit:        make(chan struct{}, 1),
		ConsumeInterval:  time.Duration(10) * time.Second,
		MinBatchableSize: 1,
		MaxBatchableSize: 1000,
		MaxBatchableWait: time.Duration(5) * time.Second,
	}
	return _i
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
		msg, err := i.KafkaWorker.Consumer.ReadMessage(i.ConsumeInterval)
		if err != nil {
			continue
		}
		var l core.Log
		if err := json.Unmarshal(msg.Value, &l); err != nil {
			return fmt.Errorf("error unmarshalling value %v: %v", msg.Value, err)
		}
		i.Messages.Lock()
		i.Messages.Data = append(i.Messages.Data, &l)
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
	for _, entry := range i.Messages.Data {
		// metadata fields
		t := make(map[string]interface{})

		values := reflect.ValueOf(*entry)
		types := values.Type()
		for j := 0; j < types.NumField(); j++ {
			_type := fmt.Sprintf("_%s", strings.ToLower(types.Field(j).Name))
			_value := values.Field(j).Interface()
			t[_type] = _value
		}
		// data fields
		for k, v := range entry.Data {
			t[k] = v
		}
		// arrays of same-typed data fields.
		for k, v := range entry.Data {
			switch reflect.TypeOf(v).Kind() {
			case reflect.String:
				t["string.keys"] = []string{}
				t["string.values"] = []string{}
				t["string.keys"] = append(t["string.keys"].([]string), k)
				t["string.values"] = append(t["string.values"].([]string), v.(string))
			case reflect.Int:
				t["int.keys"] = []string{}
				t["int.values"] = []int{}
				t["int.keys"] = append(t["int.keys"].([]string), k)
				t["int.values"] = append(t["int.values"].([]int), v.(int))
			case reflect.Float64:
				t["float64.keys"] = []string{}
				t["float64.values"] = []float64{}
				t["float64.keys"] = append(t["float64.keys"].([]string), k)
				t["float64.values"] = append(t["float64.values"].([]float64), v.(float64))
			}
		}
		i.Messages.Lock()
		i.Messages.TransformedData = append(i.Messages.TransformedData, t)
		i.Messages.Unlock()
	}
	return nil
}

func (i *IngesterWorker) ExtractSchemas() error {
	// i.Messages.RLock()
	// data := i.Messages.TransformedData
	// i.Messages.RUnlock()

	data := []map[string]interface{}{
		{
			// metadata
			"_channel":   "testnet",
			"_logid":     "0000-0000-0000-0000-0000",
			"_senderid":  "test-1234",
			"_timestamp": int64(1709118220916),
			"_level":     "debug",
			"_message":   "lorem ipsum dolor",
			"_data":      map[string]interface{}{"bar": "helloworld", "foo": 1},
			// fields
			"foo": 1,
			"bar": "helloworld",
			// field arrays
			"int.keys":      []string{"foo"},
			"int.values":    []int{1},
			"string.keys":   []string{"bar"},
			"string.values": []string{"helloworld"},
		},
	}

	// get the part from messages that we are saving in the db
	// we only store metadata adn field arrays. fields are only
	// materialized when needed (in the future)
	storableData := make([]map[string]interface{}, 0, len(data))
	for _, item := range data {
		kv := map[string]interface{}{}
		for k, v := range item {
			if strings.HasPrefix(k, "_") || strings.Contains(k, ".") {
				kv[k] = v
			}
		}
		storableData = append(storableData, kv)
	}

	// group messages by channel
	dataByChannel := make(map[string][]interface{})
	for _, item := range storableData {
		// we are grouping by channel (channel <=> table)
		key := item["_channel"].(string)
		if _, exists := dataByChannel[key]; !exists {
			dataByChannel[key] = []interface{}{item}
		} else {
			dataByChannel[key] = append(dataByChannel[key], item)
		}
	}

	addrs := []string{"localhost:9000"}
	chConn, err := clickhouse_connector.Conn(addrs)
	if err != nil {
		return err
	}

	// For each channel, we extract the schema
	for channel, channelValue := range dataByChannel {
		jsonStorableData, err := json.Marshal(channelValue)
		if err != nil {
			return fmt.Errorf("error marshalling storable data: %v", err)
		}
		rows, err := chConn.Query(
			context.Background(),
			"DESC format(JSONEachRow, $1)",
			string(jsonStorableData))
		if err != nil {
			return fmt.Errorf("error describing data: %v", err)
		}

		var (
			columnTypes = rows.ColumnTypes()
			vars        = make([]interface{}, len(columnTypes))
		)
		for j := range columnTypes {
			vars[j] = reflect.New(columnTypes[j].ScanType()).Interface()
		}


		// We will populate this with the pairs (column_name, column_type)
		// for further processing.

		var chFields []string
		for rows.Next() {
			if err := rows.Scan(vars...); err != nil {
				return fmt.Errorf("error reading clickhouse description: %v", err)
			}
			for ndx, v := range vars {
				if ndx == 2 { break }  // we are just interested in the first 2 ones
				switch v := v.(type) {
				case *string:
					chFields = append(chFields, *v)
				}
			}
		}
		// process the fields that we extracted
		i.processFields(channel, chFields)
	}
	return nil
}

func (i *IngesterWorker) Sink() error {
	// 1. alter table if needed for each slice
	// 2. sink the data to clickhouse
	// 3. write to CanCommit channel
	i.CanCommit <- struct{}{}
	return nil
}

func (i *IngesterWorker) Commit() error {
	// 1. commit to current offset in kafka
	// 2. log about batch processing completion
	return nil
}

// processFields will take a slice of the form [column_name, column_type, ...]
// and produce an intermediate representation with it that will later be used 
// in the batching steps to define how to create or alter tables before sinking
func (i *IngesterWorker) processFields(channel string, chFields []string) error {
	return nil
}
