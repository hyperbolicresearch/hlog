// We need a dedicated file to handle the process of adding the data effiently to ClickHouse.

package ingest

import (
	"sync"
)

type Batcher interface {
	
}

type BatcherWorker struct {
	sync.RWMutex
}
