package ingest

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type Batcher interface {
	Sink(data []map[string]interface{}, endC chan struct{}) (count int, err error)
}

type BatcherWorker struct {
	sync.RWMutex
	Conn clickhouse.Conn
	// IterCount keeps track of the number of iterations of sinking of the
	// current worker
	IterCount int64
}

// Sync will receive a slice of data ready to be added to ClickHouse
// and will proceed to the dumping of those data in an efficient
// manner, givent the shape of the data.
func (b *BatcherWorker) Sink(data []map[string]interface{}) (count int, err error) {
	dataByChannel := GetDataByChannel(data)
	count = 0
	for channel, item := range dataByChannel {
		insertQuery := fmt.Sprintf("INSERT INTO %s", channel)
		batch, err := b.Conn.PrepareBatch(context.Background(), insertQuery)
		if err != nil {
			return count, err
		}
		for i := 0; i < len(item); i++ {
			// We got the map[string]interface{}, we will now extract the slice of
			// values ([]interface{}) to append
			_, _, sortedValues, err := SortMap(item[i])
			if err != nil {
				panic(err)
			}
			err = batch.Append(sortedValues...)
			if err != nil {
				log.Printf("Error inserting item: %v into %v. Error=%v", sortedValues, channel, err)
			} else {
				b.Lock()
				b.IterCount += 1
				b.Unlock()
				count += 1
			}
		}
	}
	return count, nil
}

// CREATE TABLE IF NOT EXISTS default (
// 	`_channel` Nullable(String),
// 	`_data` Tuple(
// 	  bar Nullable(String),
// 	  count Nullable(Int64),
// 	  foo Nullable(String)),
// 	`_level` Nullable(String),
// 	`_logid` String,
// 	`_message` Nullable(String),
// 	`_senderid` Nullable(String),
// 	`_timestamp` Nullable(Int64),
// 	`float64.keys` Array(Nullable(String)),
// 	`float64.values` Array(Nullable(Int64)),
// 	`string.keys` Array(Nullable(String)),
// 	`string.values` Array(Nullable(String)),
//   )
//   ENGINE = MergeTree
//   PRIMARY KEY (_logid)
//   ORDER BY _logid