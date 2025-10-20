package cproxy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smarty/gunit"
	"github.com/smarty/gunit/assert/should"
)

func TestConfigFixture(t *testing.T) {
	gunit.Run(new(ConfigFixture), t)
}

type ConfigFixture struct {
	*gunit.Fixture

	clientSocket    *TestSocket
	serverSocket    *TestSocket
	clientConnector *TestClientConnector
	dialer          *TestDialer
	handler         http.Handler
}

func (this *ConfigFixture) Setup() {
	this.clientSocket = NewTestSocket()
	this.serverSocket = NewTestSocket()
	this.clientConnector = NewTestClientConnector(this.clientSocket)
	this.dialer = NewTestDialer(this.serverSocket)
	this.handler = New(
		Options.ClientConnector(this.clientConnector),
		Options.Dialer(this.dialer),
	)
}

func (this *ConfigFixture) TestEndToEnd() {
	request := httptest.NewRequest("CONNECT", "host", nil)
	response := httptest.NewRecorder()

	this.clientSocket.readBuffer.WriteString("from client")
	this.serverSocket.readBuffer.WriteString("from server")

	this.handler.ServeHTTP(response, request)

	this.So(this.dialer.address, should.Equal, "host")
	this.So(this.clientSocket.writeBuffer.String(), should.Equal, "HTTP/1.1 200 OK\r\n\r\nfrom server")
	this.So(this.serverSocket.writeBuffer.String(), should.Equal, "from client")
}
