package cproxy

import "github.com/smartystreets/logging"

type LoggingInitializer struct {
	logger *logging.Logger
	inner  Initializer
}

func NewLoggingInitializer(inner Initializer) *LoggingInitializer {
	return &LoggingInitializer{inner: inner}
}

func (this *LoggingInitializer) Initialize(client, server Socket) bool {
	result := this.inner.Initialize(client, server)

	if result {
		this.logger.Printf("[INFO] Established connection [%s] -> [%s]", client.RemoteAddr(), server.RemoteAddr())
	} else {
		this.logger.Printf("[INFO] Connection failed [%s] -> [%s]", client.RemoteAddr(), server.RemoteAddr())
	}

	return result
}
