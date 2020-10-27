package cproxy

import (
	"net/http"
	"time"
)

func New(options ...option) http.Handler {
	var this configuration
	Options.apply(options...)(&this)
	return newHandler(this.Filter, this.ClientConnector, this.ServerConnector, this.Monitor)
}

var Options singleton

type singleton struct{}
type option func(*configuration)

type configuration struct {
	DialTimeout     time.Duration
	Filter          Filter
	Dialer          Dialer
	LogConnections  bool
	ProxyProtocol   bool
	Initializer     initializer
	ClientConnector clientConnector
	ServerConnector serverConnector
	Monitor         monitor
	Logger          logger
}

func (singleton) DialTimeout(value time.Duration) option {
	return func(this *configuration) { this.DialTimeout = value }
}
func (singleton) Filter(value Filter) option {
	return func(this *configuration) { this.Filter = value }
}
func (singleton) ClientConnector(value clientConnector) option {
	return func(this *configuration) { this.ClientConnector = value }
}
func (singleton) Dialer(value Dialer) option {
	return func(this *configuration) { this.Dialer = value }
}
func (singleton) LogConnections(value bool) option {
	return func(this *configuration) { this.LogConnections = value }
}
func (singleton) ProxyProtocol(value bool) option {
	return func(this *configuration) { this.ProxyProtocol = value }
}
func (singleton) Initializer(value initializer) option {
	return func(this *configuration) { this.Initializer = value }
}
func (singleton) ServerConnector(value serverConnector) option {
	return func(this *configuration) { this.ServerConnector = value }
}
func (singleton) Monitor(value monitor) option {
	return func(this *configuration) { this.Monitor = value }
}
func (singleton) Logger(value logger) option {
	return func(this *configuration) { this.Logger = value }
}

func (singleton) apply(options ...option) option {
	return func(this *configuration) {
		for _, option := range Options.defaults(options...) {
			option(this)
		}

		if this.Dialer == nil {
			this.Dialer = newDialer(this)
		}

		if this.Initializer == nil && this.ProxyProtocol {
			this.Initializer = newProxyProtocolInitializer()
		}

		this.Initializer = newLoggingInitializer(this)

		if this.ServerConnector == nil {
			this.ServerConnector = newServerConnector(this.Dialer, this.Initializer)
		}
	}
}
func (singleton) defaults(options ...option) []option {
	const defaultDialTimeout = time.Second * 10
	var defaultFilter = newFilter()
	var defaultClientConnector = newClientConnector()
	var defaultInitializer = newInitializer()
	var defaultMonitor = nop{}
	var defaultLogger = nop{}
	return append([]option{
		Options.DialTimeout(defaultDialTimeout),
		Options.Filter(defaultFilter),
		Options.ClientConnector(defaultClientConnector),
		Options.Initializer(defaultInitializer),
		Options.Monitor(defaultMonitor),
		Options.Logger(defaultLogger),
	}, options...)
}

type nop struct{}

func (nop) Measure(int)                   {}
func (nop) Printf(string, ...interface{}) {}
