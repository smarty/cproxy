package cproxy

import (
	"io"
	"net"
	"net/http"
)

type (
	Filter interface {
		IsAuthorized(*http.Request, http.ResponseWriter) bool
	}

	ClientConnector interface {
		Connect(w http.ResponseWriter) Socket
	}
)

type (
	Dialer interface {
		Dial(string) Socket
	}

	ServerConnector interface {
		Connect(Socket, string) Proxy
	}

	Initializer interface {
		Initialize(Socket, Socket) bool
	}

	Proxy interface {
		Proxy()
	}
)

type (
	Socket interface {
		io.ReadWriteCloser
		RemoteAddr() net.Addr
	}

	TCPSocket interface {
		Socket
		CloseRead() error
		CloseWrite() error
	}
)

type Meter interface {
	Measure(int)
}

const (
	MeasurementHTTPRequest int = iota
	MeasurementBadMethod
	MeasurementUnauthorizedRequest
	MeasurementClientConnectionFailed
	MeasurementServerConnectionFailed
	MeasurementProxyReady
	MeasurementProxyComplete
)
