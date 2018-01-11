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

	filter := this.filter
	if filter == nil {
		filter = NewFilter()
func (this *Wireup) build() http.Handler {
	}

	clientConnector := this.clientConnector
	if clientConnector == nil {
		clientConnector = NewClientConnector()
	}

	serverConnector := this.buildServerConnector()

	meter := this.meter
	if meter == nil {
		meter = NewMeter()
	}

	return NewHandler(filter, clientConnector, serverConnector, meter)
}
func (this *Wireup) buildServerConnector() ServerConnector {
	dialer := this.dialer
	if dialer == nil {
		dialer = NewDialer()
	}

	initializer := this.initializer
	if initializer == nil {
		initializer = NewInitializer()
	}

	serverConnector := this.serverConnector
	if serverConnector == nil {
		serverConnector = NewServerConnector(dialer, initializer)
	}

	return serverConnector
}
