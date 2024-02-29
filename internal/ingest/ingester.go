package ingest

type Ingester interface {
	Start()
	Stop() error
	Consume() error
	Sink() error
	// Transform() error
	// ExtractSchemas() error
	// Commit() error
}