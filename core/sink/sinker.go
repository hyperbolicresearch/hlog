package sink

type Sinker interface {
	Sink(data []map[string]interface{}, endC chan struct{}) (count int, err error)
}