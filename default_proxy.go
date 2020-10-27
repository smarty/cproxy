package cproxy

import (
	"io"
	"sync"
)

type defaultProxy struct {
	client socket
	server socket
	waiter *sync.WaitGroup
}

func newProxy(client, server socket) *defaultProxy {
	waiter := &sync.WaitGroup{}
	waiter.Add(2) // wait on both client->server and server->client streams

	return &defaultProxy{
		waiter: waiter,
		client: client,
		server: server,
	}
}

func (this *defaultProxy) Proxy() {
	go this.streamAndClose(this.client, this.server)
	go this.streamAndClose(this.server, this.client)
	this.closeSockets()
}

func (this *defaultProxy) streamAndClose(reader, writer socket) {
	_, _ = io.Copy(writer, reader)

	tryCloseRead(reader)
	tryCloseWrite(writer)

	this.waiter.Done()
}
func tryCloseRead(socket socket) {
	if tcp, ok := socket.(tcpSocket); ok {
		_ = tcp.CloseRead()
	}
}
func tryCloseWrite(socket socket) {
	if tcp, ok := socket.(tcpSocket); ok {
		_ = tcp.CloseWrite()
	}
}

func (this *defaultProxy) closeSockets() {
	this.waiter.Wait()
	_ = this.client.Close()
	_ = this.server.Close()
}
