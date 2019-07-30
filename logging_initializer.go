package cproxy

import "github.com/smartystreets/logging"

type LoggingInitializer struct {
	logger *logging.Logger
	inner  Initializer
}

func NewLoggingInitializer(inner Initializer) *LoggingInitializer {
	return &LoggingInitializer{inner: inner}
}

func (it *LoggingInitializer) Initialize(client, server Socket) bool {
	result := it.inner.Initialize(client, server)

	if result {
		it.logger.Printf("[INFO] Established connection [%s] -> [%s]", client.RemoteAddr(), server.RemoteAddr())
	} else {
		it.logger.Printf("[INFO] Connection failed [%s] -> [%s]", client.RemoteAddr(), server.RemoteAddr())
	}

	return result
}
