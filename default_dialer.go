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

func (this *DefaultDialer) WithLogging() *DefaultDialer {
	this.logging = true
	return this
}

func (this *DefaultDialer) Dial(address string) Socket {
	if socket, err := net.DialTimeout("tcp", address, this.timeout); err == nil {
		return socket
	} else if this.logging {
		this.logger.Printf("[INFO] Unable to establish connection to [%s]: %s", address, err)
	}

	return nil
}
