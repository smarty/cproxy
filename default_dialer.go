package cproxy

import (
	"net"
	"time"
)

type DefaultDialer struct {
	timeout time.Duration
}

func NewDialer(timeout time.Duration) *DefaultDialer {
	return &DefaultDialer{timeout: timeout}
}

func (this *DefaultDialer) Dial(address string) Socket {
	if socket, err := net.DialTimeout("tcp", address, this.timeout); err != nil {
		return nil
	} else {
		return socket
	}
}
