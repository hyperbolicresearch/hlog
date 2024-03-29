package sink

type MongoDBSinkerWorker struct {

}

// Sink 
func (w *MongoDBSinkerWorker) Sink(data []map[string]interface{}) (count int, err error) {

	return 0, nil
}