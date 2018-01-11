package cproxy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestWireupFixture(t *testing.T) {
	gunit.Run(new(WireupFixture), t)
}

type WireupFixture struct {
	*gunit.Fixture

	clientSocket    *TestSocket
	serverSocket    *TestSocket
	clientConnector *TestClientConnector
	dialer          *TestDialer
	handler         http.Handler
}

func (this *WireupFixture) Setup() {
	this.clientSocket = NewTestSocket()
	this.serverSocket = NewTestSocket()
	this.clientConnector = NewTestClientConnector(this.clientSocket)
	this.dialer = NewTestDialer(this.serverSocket)
	this.handler = Configure(
		WithClientConnector(this.clientConnector),
		WithDialer(this.dialer),
	)
}

func (this *WireupFixture) TestEndToEnd() {
	request := httptest.NewRequest("CONNECT", "host", nil)
	response := httptest.NewRecorder()

	this.clientSocket.readBuffer.WriteString("from client")
	this.serverSocket.readBuffer.WriteString("from server")

	this.handler.ServeHTTP(response, request)

	this.So(this.dialer.address, should.Equal, "host")
	this.So(this.clientSocket.writeBuffer.String(), should.Equal, "HTTP/1.1 200 OK\r\n\r\nfrom server")
	this.So(this.serverSocket.writeBuffer.String(), should.Equal, "from client")
}
