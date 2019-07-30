package cproxy

import (
	"net"
	"time"

	"github.com/smartystreets/logging"
)

type DefaultDialer struct {
	timeout time.Duration
	logging bool
	logger  *logging.Logger
}

func NewDialer() *DefaultDialer {
	return NewDialerWithTimeout(time.Second * 10)
}

func NewDialerWithTimeout(timeout time.Duration) *DefaultDialer {
	return &DefaultDialer{timeout: timeout}
}

func (it *DefaultDialer) WithLogging() *DefaultDialer {
	it.logging = true
	return it
}

func (it *DefaultDialer) Dial(address string) Socket {
	if socket, err := net.DialTimeout("tcp", address, it.timeout); err == nil {
		return socket
	} else if it.logging {
		it.logger.Printf("[INFO] Unable to establish connection to [%s]: %s", address, err)
	}

	return nil
}
