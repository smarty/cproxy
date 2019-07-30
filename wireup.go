package cproxy

import "net/http"

type Wireup struct {
	filter          Filter
	dialer          Dialer
	initializer     Initializer
	clientConnector ClientConnector
	serverConnector ServerConnector
	meter           Meter
}

func Configure(options ...Option) http.Handler {
	wireup := &Wireup{}
	for _, option := range options {
		option(wireup)
	}
	return wireup.build()
}

type Option func(*Wireup)

func WithFilter(value Filter) Option {
	return func(it *Wireup) { it.filter = value }
}
func WithClientConnector(value ClientConnector) Option {
	return func(it *Wireup) { it.clientConnector = value }
}
func WithDialer(value Dialer) Option {
	return func(it *Wireup) { it.dialer = value }
}
func WithInitializer(value Initializer) Option {
	return func(it *Wireup) { it.initializer = value }
}
func WithServerConnector(value ServerConnector) Option {
	return func(it *Wireup) { it.serverConnector = value }
}
func WithMeter(value Meter) Option {
	return func(it *Wireup) { it.meter = value }
}

func (it *Wireup) build() http.Handler {
	if it.filter == nil {
		it.filter = NewFilter()
	}

	if it.clientConnector == nil {
		it.clientConnector = NewClientConnector()
	}

	if it.meter == nil {
		it.meter = NewMeter()
	}

	return NewHandler(it.filter, it.clientConnector, it.buildServerConnector(), it.meter)
}
func (it *Wireup) buildServerConnector() ServerConnector {
	if it.dialer == nil {
		it.dialer = NewDialer()
	}

	if it.initializer == nil {
		it.initializer = NewInitializer()
	}

	if it.serverConnector == nil {
		it.serverConnector = NewServerConnector(it.dialer, it.initializer)
	}

	return it.serverConnector
}
