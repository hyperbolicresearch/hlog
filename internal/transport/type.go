package transport

type TransportLayer interface {
	// Start starts the listening porcess of the transport layer
	// it can be by example starting the HTTP server or start
	// consuming to a particular topic by the Kafka consumer.
	Start() error

	// Stop actually stops the listening process. It can mean by
	// example stopping the HTTP server or stopping the Kafka
	// consumer consuming on the particular topic.
	Stop() error
}