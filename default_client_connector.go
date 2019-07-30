package cproxy

import "net/http"

type DefaultClientConnector struct {
}

func NewClientConnector() *DefaultClientConnector {
	return &DefaultClientConnector{}
}

func (it *DefaultClientConnector) Connect(response http.ResponseWriter) Socket {
	if hijacker, ok := response.(http.Hijacker); !ok {
		return nil
	} else if socket, _, _ := hijacker.Hijack(); socket == nil {
		return nil // this 'else if' exists to avoid the pointer nil != interface nil issue
	} else {
		return socket
	}
}
