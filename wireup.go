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

func Configure() *Wireup {
	return &Wireup{}
}

func (this *Wireup) WithFilter(value Filter) *Wireup {
	this.filter = value
	return this
}

func (this *Wireup) WithClientConnector(value ClientConnector) *Wireup {
	this.clientConnector = value
	return this
}

func (this *Wireup) WithDialer(value Dialer) *Wireup {
	this.dialer = value
	return this
}
func (this *Wireup) WithInitializer(value Initializer) *Wireup {
	this.initializer = value
	return this
}
func (this *Wireup) WithServerConnector(value ServerConnector) *Wireup {
	this.serverConnector = value
	return this
}

func (this *Wireup) WithMeter(value Meter) *Wireup {
	this.meter = value
	return this
}

func (this *Wireup) Build() http.Handler {
	filter := this.filter
	if filter == nil {
		filter = NewFilter()
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
