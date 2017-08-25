package cproxy

import "net/http"

type DefaultHandler struct {
	filter          Filter
	clientConnector ClientConnector
	serverConnector ServerConnector
	meter           Meter
}

func NewHandler(filter Filter, clientConnector ClientConnector, serverConnector ServerConnector, meter Meter) *DefaultHandler {
	return &DefaultHandler{
		filter:          filter,
		clientConnector: clientConnector,
		serverConnector: serverConnector,
		meter:           meter,
	}
}

func (this *DefaultHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	this.meter.Measure(MeasurementHTTPRequest)

	if request.Method != "CONNECT" {
		this.meter.Measure(MeasurementBadMethod)
		writeResponseStatus(response, http.StatusMethodNotAllowed)

	} else if !this.filter.IsAuthorized(request) {
		this.meter.Measure(MeasurementUnauthorizedRequest)
		writeResponseStatus(response, http.StatusUnauthorized)

	} else if client := this.clientConnector.Connect(response); client == nil {
		this.meter.Measure(MeasurementClientConnectionFailed)
		writeResponseStatus(response, http.StatusNotImplemented)

	} else if proxy := this.serverConnector.Connect(client, request.URL.Host); proxy == nil {
		this.meter.Measure(MeasurementServerConnectionFailed)
		client.Write(StatusBadGateway)
		client.Close()

	} else {
		this.meter.Measure(MeasurementProxyReady)
		client.Write(StatusReady)
		proxy.Proxy()
		this.meter.Measure(MeasurementProxyComplete)
	}
}

func writeResponseStatus(response http.ResponseWriter, statusCode int) {
	http.Error(response, http.StatusText(statusCode), statusCode)
}

var (
	StatusBadGateway = []byte("HTTP/1.1 502 Bad Gateway\r\n\r\n")
	StatusReady      = []byte("HTTP/1.1 200 OK\r\n\r\n")
)
