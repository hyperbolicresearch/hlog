package ingest

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/hyperbolicresearch/hlog/config"
	"github.com/hyperbolicresearch/hlog/internal/clickhouseservice"
	"github.com/hyperbolicresearch/hlog/internal/core"
	"github.com/hyperbolicresearch/hlog/internal/kafkaservice"
	"github.com/hyperbolicresearch/hlog/internal/mongodb"
)

// IngesterWorker is responsible the handle the end-to-end dumping
// of data to ClickHouse
type IngesterWorker struct {
	sync.RWMutex
	*BatcherWorker
	*kafkaservice.KafkaWorker
	MongoDatabase *mongo.Database
	IsRunning     bool
	IsIngesting   bool
	Messages      *Messages
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

// Messages is the data structure holding the messages that will be
// written to ClickHouse.
type Messages struct {
	sync.RWMutex
	Data            []*core.Log
	TransformedData []map[string]interface{}
	StorableData    []map[string]interface{}
}

// TODO: make configs
func NewClickHouseIngester(cfg *config.Config) *IngesterWorker {
	kw, err := kafkaservice.NewKafkaWorker(&cfg.ClickHouse.KafkaConfigs)
	if err != nil {
		panic("failed to create ingester")
	}
	err = kw.ConfigureConsumer()
	if err != nil {
		panic(err)
	}
	err = kw.SubscribeTopics(cfg.ClickHouse.KafkaTopics)
	if err != nil {
		panic(err)
	}
	mongoClient := mongodb.Client(cfg.MongoDB.Server)
	db := mongoClient.Database(cfg.MongoDB.Database)

	// TODO make configurable
	addrs := []string{"127.0.0.1:9000"}
	chConn, err := clickhouseservice.Conn(addrs)
	if err != nil {
		panic(err)
	}

	_i := &IngesterWorker{
		MongoDatabase: db,
		BatcherWorker: &BatcherWorker{
			Conn: chConn,
		},
		Messages:         &Messages{},
		KafkaWorker:      kw,
		ConsumeInterval:  cfg.ClickHouse.ConsumeInterval,
		MinBatchableSize: cfg.ClickHouse.MinBatchableSize,
		MaxBatchableSize: cfg.ClickHouse.MaxBatchableSize,
		MaxBatchableWait: cfg.ClickHouse.MaxBatchableWait,
	}
	return _i
}

// Start spins up the consuming process generally speaking. It runs as
// long as i.IsRunning is true and tracks new incomming logs
func (i *IngesterWorker) Start(stop chan os.Signal) {
	i.RLock()
	if i.IsRunning {
		// TODO: Should I really panic here?
		panic("Ingester worker already running...")
	}
	i.RUnlock()

	i.Lock()
	i.IsRunning = true
	i.Unlock()

	consumeTicker := time.NewTicker(i.ConsumeInterval)

	for i.IsRunning {
		select {
		case <-stop:
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
func (i *IngesterWorker) Consume() {
	// We will try to extract as much as possible messages from
	// Kafka given that j <= i.MaxBatchableSize.
	// Doing like that, we make sure that we always read less
	// or equal to the i.MaxBatchableSize.
	for j := 0; j < i.MaxBatchableSize; j++ {
		msg, err := i.KafkaWorker.Consumer.ReadMessage(i.ConsumeInterval)
		if err != nil {
			continue
		}
		fmt.Printf("%+v\n", string(msg.Value))
		var l core.Log
		if err := json.Unmarshal(msg.Value, &l); err != nil {
			 panic(fmt.Errorf("error unmarshalling value %v: %v", msg.Value, err))
		}
		i.Messages.Lock()
		i.Messages.Data = append(i.Messages.Data, &l)
		i.Messages.Unlock()
	}

	fmt.Println("At least we got here")
	// The following pipeline will perfom the preparation and the
	// sink of the buffered data to ClickHouse.
	_ = i.Transform()
	_ = i.ExtractSchemas()
	_, err := i.Sink(i.Messages.StorableData)
	if err != nil {
		panic(err)
	}

	// Waiting until we have the acks that we successfully sink
	// the data to ClickHouse (which should be sent from inside
	// Sink)
	_, err = i.Consumer.Commit()
	if err != nil {
		panic(err)
	}
}

// Transform will flatten the message to the appropriate format
// that will be stored to ClickHouse, add metadata.
func (i *IngesterWorker) Transform() error {
	fmt.Println("Transforming")
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
		sortedT, _, _, err := SortMap(t)
		if err != nil {
			return err
		}
		i.Messages.Lock()
		i.Messages.TransformedData = append(i.Messages.TransformedData, sortedT)
		i.Messages.Unlock()
	}
	return nil
}

// ExtractSchemas takes a bunch of data and extracts the SQL compatible
// schema out of them.
func (i *IngesterWorker) ExtractSchemas() error {
	fmt.Println("Extracting schema")
	i.Messages.RLock()
	data := i.Messages.TransformedData
	i.Messages.RUnlock()
	// get the part from messages that we are saving in the db
	// we only store metadata adn field arrays. fields are only
	// materialized when needed (in the future)
	storableData := GetStorableData(data)
	i.Messages.Lock()
	i.Messages.StorableData = append(i.Messages.StorableData, storableData...)
	i.Messages.Unlock()
	// group messages by channel
	dataByChannel := GetDataByChannel(storableData)

	// TODO Make config
	addrs := []string{"localhost:9000"}
	chConn, err := clickhouseservice.Conn(addrs)
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
		var chFields []string
		for rows.Next() {
			if err := rows.Scan(vars...); err != nil {
				return fmt.Errorf("error reading clickhouse description: %v", err)
			}
			for ndx, v := range vars {
				if ndx == 2 {
					break
				} // we are just interested in the first 2 ones
				switch v := v.(type) {
				case *string:
					chFields = append(chFields, *v)
				}
			}
		}
		i.processFields(channel, chFields)
	}
	return nil
}

// processFields will take a slice of the form [column_name, column_type, ...]
// and produce an intermediate representation with it that will later be used
// in the batching steps to define how to create or alter tables before sinking
func (i *IngesterWorker) processFields(channel string, chFields []string) error {
	// we create a map-representation of the channel fields (chFields)
	// in the form: {...field_name, ...field_type}
	repr := map[string]interface{}{}
	for j := 0; j <= len(chFields)-2; j += 2 {
		key := chFields[j]
		value := chFields[j+1]
		repr[key] = value
	}

	i.RLock()
	col := i.MongoDatabase.Collection("_sqlschemas")
	i.RUnlock()
	filter := bson.D{}
	var result map[string]interface{}
	var toCreate bool = false
	err := col.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		// means that either there is no schema because that's the first
		// time this clientID sends logs or that the data was corrupted.
		// we therefore create a new one
		_, err := col.InsertOne(context.TODO(), repr)
		if err != nil {
			panic(err)
		}
		result = repr
		toCreate = true
	}
	if toCreate {
		// CREATE TABLE...
		err := GenerateSQLAndApply(result, channel, false)
		if err != nil {
			panic(err)
		}
	} else {
		// that means the document was found, which implies that we
		// should verify whether we should update fields or not
		toUpdate := map[string]interface{}{}
		for key, value := range repr {
			if _, ok := result[key]; !ok {
				toUpdate[key] = value
			}
		}
		if len(toUpdate) > 0 {
			// ALTER TABLE...
			err := GenerateSQLAndApply(toUpdate, channel, true)
			if err != nil {
				panic(err)
			}
		}
	}
	return nil
}
