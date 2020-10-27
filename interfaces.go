package cproxy

import (
	"io"
	"net"
	"net/http"
)

type (
	Filter interface {
		IsAuthorized(http.ResponseWriter, *http.Request) bool
	}

	clientConnector interface {
		Connect(http.ResponseWriter) socket
	}
)

type (
	dialer interface {
		Dial(string) socket
	}

	serverConnector interface {
		Connect(socket, string) proxy
	}

	initializer interface {
		Initialize(socket, socket) bool
	}

	proxy interface {
		Proxy()
	}
)

type (
	socket interface {
		io.ReadWriteCloser
		RemoteAddr() net.Addr
	}

	tcpSocket interface {
		socket
		CloseRead() error
		CloseWrite() error
	}
)

type (
	monitor interface {
		Measure(int)
	}
	logger interface {
		Printf(string, ...interface{})
	}
)

const (
	MeasurementHTTPRequest int = iota
	MeasurementBadMethod
	MeasurementUnauthorizedRequest
	MeasurementClientConnectionFailed
	MeasurementServerConnectionFailed
	MeasurementProxyReady
	MeasurementProxyComplete
)
