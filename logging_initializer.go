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
	if !this.inner.Initialize(client, server) {
		return false
	}

	this.logger.Printf("[INFO] Established connection [%s] -> [%s]", client.RemoteAddr(), server.RemoteAddr())
	return true
}
