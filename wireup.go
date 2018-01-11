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
	return func(this *Wireup) { this.filter = value }
}
func WithClientConnector(value ClientConnector) Option {
	return func(this *Wireup) { this.clientConnector = value }
}
func WithDialer(value Dialer) Option {
	return func(this *Wireup) { this.dialer = value }
}
func WithInitializer(value Initializer) Option {
	return func(this *Wireup) { this.initializer = value }
}
func WithServerConnector(value ServerConnector) Option {
	return func(this *Wireup) { this.serverConnector = value }
}
func WithMeter(value Meter) Option {
	return func(this *Wireup) { this.meter = value }
}

func (this *Wireup) build() http.Handler {
	if this.filter == nil {
		this.filter = NewFilter()
	}

	if this.clientConnector == nil {
		this.clientConnector = NewClientConnector()
	}

	if this.meter == nil {
		this.meter = NewMeter()
	}

	return NewHandler(this.filter, this.clientConnector, this.buildServerConnector(), this.meter)
}
func (this *Wireup) buildServerConnector() ServerConnector {
	if this.dialer == nil {
		this.dialer = NewDialer()
	}

	if this.initializer == nil {
		this.initializer = NewInitializer()
	}

	if this.serverConnector == nil {
		this.serverConnector = NewServerConnector(this.dialer, this.initializer)
	}

	return this.serverConnector
}
