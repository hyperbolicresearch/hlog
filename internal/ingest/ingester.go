package ingest

type Ingester interface {
	Start()
	Stop() error
	Consume()
}
