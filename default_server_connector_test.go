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

	connector *defaultServerConnector
}

func (this *ServerConnectorFixture) Setup() {
	this.clientSocket = NewTestSocket()
	this.serverSocket = NewTestSocket()
	this.dialer = NewTestDialer(this.serverSocket)
	this.initializer = NewTestInitializer(true)
	this.connector = newServerConnector(this.dialer, this.initializer)
}

//////////////////////////////////////////////////////////////

func (this *ServerConnectorFixture) TestFailedDialReturnsNoProxy() {
	this.dialer.socket = nil

	proxy := this.connect("address")

	this.So(proxy, should.BeNil)
	this.So(this.dialer.address, should.Equal, "address")
}

func (this *ServerConnectorFixture) TestFailedInitializationClosesServerSocketAndReturnsNoProxy() {
	this.initializer.success = false

	proxy := this.connect("a")

	this.So(proxy, should.BeNil)
	this.So(this.clientSocket.close, should.Equal, 0)
	this.So(this.serverSocket.close, should.Equal, 1)
}

func (this *ServerConnectorFixture) TestSuccessfulConnectionYieldsInitializedProxy() {
	proxy := this.connect("a")

	this.So(proxy, should.NotBeNil)
	this.So(proxy.(*defaultProxy).client, should.Equal, this.clientSocket)
	this.So(proxy.(*defaultProxy).server, should.Equal, this.serverSocket)
	this.So(this.initializer.client, should.Equal, this.clientSocket)
	this.So(this.initializer.server, should.Equal, this.serverSocket)
	this.So(this.clientSocket.close, should.Equal, 0)
	this.So(this.serverSocket.close, should.Equal, 0)
}

func (this *ServerConnectorFixture) connect(address string) proxy {
	return this.connector.Connect(this.clientSocket, address)
}

//////////////////////////////////////////////////////////////

type TestDialer struct {
	address string
	socket  socket
}

func NewTestDialer(socket socket) *TestDialer {
	return &TestDialer{socket: socket}
}

func (this *TestDialer) Dial(address string) socket {
	this.address = address
	return this.socket
}

//////////////////////////////////////////////////////////////

type TestInitializer struct {
	success bool
	client  socket
	server  socket
}

func NewTestInitializer(success bool) *TestInitializer {
	return &TestInitializer{success: success}
}

func (this *TestInitializer) Initialize(client, server socket) bool {
	this.client = client
	this.server = server
	return this.success
}
