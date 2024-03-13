package ingest

import (
	"sync"
)

type Batcher interface {
	Sink(data []map[string]interface{}, endC chan struct{}) (int, error)
	
}

type BatcherWorker struct {
	sync.RWMutex
}

// Sync will receive a slice of data ready to be added to ClickHouse
// and will proceed to the dumping of those data in an efficient
// manner, givent the shape of the data.
func (b *BatcherWorker) Sink(data []map[string]interface{}, endC chan struct{}) (int, error) {
	endC <- struct{}{}
	return 0, nil
}
