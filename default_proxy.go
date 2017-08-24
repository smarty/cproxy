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

func (this *DefaultProxy) Proxy() {
	go this.streamAndClose(this.client, this.server)
	go this.streamAndClose(this.server, this.client)
	this.closeSockets()
}

func (this *DefaultProxy) streamAndClose(reader, writer Socket) {
	io.Copy(writer, reader)

	tryCloseRead(reader)
	tryCloseWrite(writer)

	this.waiter.Done()
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

func (this *DefaultProxy) closeSockets() {
	this.waiter.Wait()
	this.client.Close()
	this.server.Close()
}
