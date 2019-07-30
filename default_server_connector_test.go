package cproxy

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestServerConnectorFixture(t *testing.T) {
	gunit.Run(new(ServerConnectorFixture), t)
}

type ServerConnectorFixture struct {
	*gunit.Fixture

	clientSocket *TestSocket
	serverSocket *TestSocket
	dialer       *TestDialer
	initializer  *TestInitializer

	connector *DefaultServerConnector
}

func (it *ServerConnectorFixture) Setup() {
	it.clientSocket = NewTestSocket()
	it.serverSocket = NewTestSocket()
	it.dialer = NewTestDialer(it.serverSocket)
	it.initializer = NewTestInitializer(true)
	it.connector = NewServerConnector(it.dialer, it.initializer)
}

//////////////////////////////////////////////////////////////

func (it *ServerConnectorFixture) TestFailedDialReturnsNoProxy() {
	it.dialer.socket = nil

	proxy := it.connect("address")

	it.So(proxy, should.BeNil)
	it.So(it.dialer.address, should.Equal, "address")
}

func (it *ServerConnectorFixture) TestFailedInitializationClosesServerSocketAndReturnsNoProxy() {
	it.initializer.success = false

	proxy := it.connect("a")

	it.So(proxy, should.BeNil)
	it.So(it.clientSocket.close, should.Equal, 0)
	it.So(it.serverSocket.close, should.Equal, 1)
}

func (it *ServerConnectorFixture) TestSuccessfulConnectionYieldsInitializedProxy() {
	proxy := it.connect("a")

	it.So(proxy, should.NotBeNil)
	it.So(proxy.(*DefaultProxy).client, should.Equal, it.clientSocket)
	it.So(proxy.(*DefaultProxy).server, should.Equal, it.serverSocket)
	it.So(it.initializer.client, should.Equal, it.clientSocket)
	it.So(it.initializer.server, should.Equal, it.serverSocket)
	it.So(it.clientSocket.close, should.Equal, 0)
	it.So(it.serverSocket.close, should.Equal, 0)
}

func (it *ServerConnectorFixture) connect(address string) Proxy {
	return it.connector.Connect(it.clientSocket, address)
}

//////////////////////////////////////////////////////////////

type TestDialer struct {
	address string
	socket  Socket
}

func NewTestDialer(socket Socket) *TestDialer {
	return &TestDialer{socket: socket}
}

func (it *TestDialer) Dial(address string) Socket {
	it.address = address
	return it.socket
}

//////////////////////////////////////////////////////////////

type TestInitializer struct {
	success bool
	client  Socket
	server  Socket
}

func NewTestInitializer(success bool) *TestInitializer {
	return &TestInitializer{success: success}
}

func (it *TestInitializer) Initialize(client, server Socket) bool {
	it.client = client
	it.server = server
	return it.success
}
