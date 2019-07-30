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

func (it *WireupFixture) Setup() {
	it.clientSocket = NewTestSocket()
	it.serverSocket = NewTestSocket()
	it.clientConnector = NewTestClientConnector(it.clientSocket)
	it.dialer = NewTestDialer(it.serverSocket)
	it.handler = Configure(
		WithClientConnector(it.clientConnector),
		WithDialer(it.dialer),
	)
}

func (it *WireupFixture) TestEndToEnd() {
	request := httptest.NewRequest("CONNECT", "host", nil)
	response := httptest.NewRecorder()

	it.clientSocket.readBuffer.WriteString("from client")
	it.serverSocket.readBuffer.WriteString("from server")

	it.handler.ServeHTTP(response, request)

	it.So(it.dialer.address, should.Equal, "host")
	it.So(it.clientSocket.writeBuffer.String(), should.Equal, "HTTP/1.1 200 OK\r\n\r\nfrom server")
	it.So(it.serverSocket.writeBuffer.String(), should.Equal, "from client")
}
