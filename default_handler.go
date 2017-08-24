package cproxy

import "net/http"

type DefaultHandler struct {
	filter          Filter
	clientConnector ClientConnector
	serverConnector ServerConnector
}

func NewHandler(filter Filter, clientConnector ClientConnector, serverConnector ServerConnector) *DefaultHandler {
	return &DefaultHandler{
		filter:          filter,
		clientConnector: clientConnector,
		serverConnector: serverConnector,
	}
}

func (this *DefaultHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if request.Method != "CONNECT" {
		writeResponseStatus(response, http.StatusMethodNotAllowed)
	} else if !this.filter.IsAuthorized(request) {
		writeResponseStatus(response, http.StatusUnauthorized)
	} else if client := this.clientConnector.Connect(response); client == nil {
		writeResponseStatus(response, http.StatusNotImplemented)
	} else if proxy := this.serverConnector.Connect(client, request.URL.Host); proxy == nil {
		client.Write(StatusBadGateway)
		client.Close()
	} else {
		client.Write(StatusReady)
		proxy.Proxy()
	}
}

func writeResponseStatus(response http.ResponseWriter, statusCode int) {
	http.Error(response, http.StatusText(statusCode), statusCode)
}

var (
	StatusBadGateway = []byte("HTTP/1.1 502 Bad Gateway\r\n\r\n")
	StatusReady      = []byte("HTTP/1.1 200 OK\r\n\r\n")
)
