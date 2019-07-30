package cproxy

import (
	"io"
	"sync"
)

type DefaultProxy struct {
	client Socket
	server Socket
	waiter *sync.WaitGroup
}

func NewProxy(client, server Socket) *DefaultProxy {
	waiter := &sync.WaitGroup{}
	waiter.Add(2) // wait on both client->server and server->client streams

	return &DefaultProxy{
		waiter: waiter,
		client: client,
		server: server,
	}
}

func (it *DefaultProxy) Proxy() {
	go it.streamAndClose(it.client, it.server)
	go it.streamAndClose(it.server, it.client)
	it.closeSockets()
}

func (it *DefaultProxy) streamAndClose(reader, writer Socket) {
	io.Copy(writer, reader)

	tryCloseRead(reader)
	tryCloseWrite(writer)

	it.waiter.Done()
}
func tryCloseRead(socket Socket) {
	if tcp, ok := socket.(TCPSocket); ok {
		tcp.CloseRead()
	}
}
func tryCloseWrite(socket Socket) {
	if tcp, ok := socket.(TCPSocket); ok {
		tcp.CloseWrite()
	}
}

func (it *DefaultProxy) closeSockets() {
	it.waiter.Wait()
	it.client.Close()
	it.server.Close()
}
