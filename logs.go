package main

import "github.com/google/uuid"

type Log struct {
	Id        uuid.UUID   `bson:"id"`
	SenderId  string      `bson:"sender_id"`
	Timestamp int64       `bson:"timestamp"`
	Level     string      `bson:"level"`
	Message   string      `bson:"message"`
	Data      interface{} `bson:"data"`
}

