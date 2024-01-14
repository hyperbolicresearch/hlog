package core

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Log struct {
	Id          uuid.UUID   `json:"id"`
	Level       string      `json:"level"`
	Message     string      `json:"message,omitempty"`
	Data        interface{} `json:"data,omitempty"`
	Timestamp   int64       `json:"timestamp,omitempty"`
	Writer      string      `json:"writer,omitempty"`
	*LogOptions `json:"options,omitempty"`
}

type LogOptions struct {
	Tags            []string  `json:"tags,omitempty"`
	TransactionId   uuid.UUID `json:"transaction_id,omitempty"`
	TransactionStep int       `json:"transaction_step,omitempty"`
}

type LogPayloadFromRequest struct {
	Level           string      `json:"level"`
	Message         string      `json:"message"`
	Data            interface{} `json:"data,omitempty"`
	Writer          string      `json:"writer"`
	Tags            []string    `json:"tags,omitempty"`
	TransactionId   string      `json:"transaction_id,omitempty"`
	TransactionStep int         `json:"transaction_step,omitempty"`
}

// NewLog creates and returns a new Log object.
func NewLog(level string, msg string, data interface{}, w string, options LogOptions) *Log {
	log := &Log{
		Id:        uuid.New(),
		Level:     level,
		Timestamp: time.Now().Unix(),
		Message:   msg,
		Data:      data,
		Writer:    w,
		LogOptions: &LogOptions{
			Tags:            options.Tags,
			TransactionId:   options.TransactionId,
			TransactionStep: options.TransactionStep,
		},
	}
	return log
}

// String returns a simplified string representation of the Log object.
func (l *Log) String() string {
	_time := time.Unix(l.Timestamp, 0).Format("Mon, Jan 02 2006, 15:04:05")
	s := fmt.Sprintf("[%v] %v >>> %v", l.Writer, _time, l.Message)
	return s
}
