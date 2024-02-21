package kafka

import (
	"errors"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer struct {
	*kafka.Consumer
}

func NewConsumer() *Consumer {
	return &Consumer{
		
	}
}

func (c *Consumer) Start() error { return errors.New("") }

func (c *Consumer) Stop() error { return errors.New("") }
