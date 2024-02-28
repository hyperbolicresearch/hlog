package core

type Log struct {
	// metadata
	LogId     string                 `json:"log_id" bson:"log_id"`
	SenderId  string                 `json:"sender_id" bson:"sender_id"`
	Timestamp int64                  `json:"timestamp" bson:"timestamp"`
	Level     string                 `json:"level" bson:"level"`
	// data
	Message   string                 `json:"message" bson:"message"`
	Data      map[string]interface{} `json:"data" bson:"data"`
}
