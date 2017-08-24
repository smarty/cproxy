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

func (this *DefaultDialer) Dial(address string) (Socket, error) {
	return net.DialTimeout("tcp", address, this.timeout)
}
