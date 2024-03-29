package sink

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type ClickHouseBatcherWorker struct {
	sync.RWMutex
	Conn clickhouse.Conn
	// IterCount keeps track of the number of iterations of sinking of the
	// current worker
	IterCount int64
}

// Sink will receive a slice of data ready to be added to ClickHouse
// and will proceed to the dumping of those data in an efficient
// manner, givent the shape of the data.
func (b *ClickHouseBatcherWorker) Sink(data []map[string]interface{}) (count int, err error) {
	fmt.Println("Sinking")
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
				log.Printf("Error inserting item: %v into %v. Error=%v",
					sortedValues, channel, err)
			} else {
				log.Println("Added to batch...")
			}
		}
		// Committing changes
		if err := batch.Send(); err != nil {
			return 0, err
		}
		// Updating the counter
		b.Lock()
		b.IterCount += 1
		b.Unlock()
		count += 1
		log.Printf("Batch %v inserted successfully into %s", count, channel)
	}
	return count, nil
}
