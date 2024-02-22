package core

type Log struct {
	LogId     string                 `json:"log_id" bson:"log_id"`
	SenderId  string                 `json:"sender_id" bson:"sender_id"`
	Timestamp int64                  `json:"timestamp" bson:"timestamp"`
	Level     string                 `json:"level" bson:"level"`
	Message   string                 `json:"message" bson:"message"`
	Data      map[string]interface{} `json:"data" bson:"data"`
}
