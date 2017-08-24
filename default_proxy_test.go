package cproxy

import (
	"bytes"
	"net"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestProxyFixture(t *testing.T) {
	gunit.Run(new(ProxyFixture), t)
}

type ProxyFixture struct {
	*gunit.Fixture

	client *TestSocket
	server *TestSocket

	proxy *DefaultProxy
}

func (this *ProxyFixture) Setup() {
	this.client = NewTestSocket()
	this.server = NewTestSocket()
	this.proxy = NewProxy(this.client, this.server)
}

func (this *ProxyFixture) TestProxyCopiesAndCloses() {
	this.client.readBuffer.WriteString("from client")
	this.server.readBuffer.WriteString("from server")

	this.proxy.Proxy()

	this.So(this.client.writeBuffer.String(), should.Equal, "from server")
	this.So(this.server.writeBuffer.String(), should.Equal, "from client")
	this.So(this.client.closeRead, should.Equal, 1)
	this.So(this.server.closeWrite, should.Equal, 1)
	this.So(this.client.closeWrite, should.Equal, 1)
	this.So(this.server.closeRead, should.Equal, 1)
	this.So(this.server.close, should.Equal, 1)
	this.So(this.client.close, should.Equal, 1)
}

//////////////////////////////////////////////////////////////

type TestSocket struct {
	readBuffer  *bytes.Buffer
	writeBuffer *bytes.Buffer
	reads       int
	writes      int
	closeRead   int
	closeWrite  int
	close       int
	address     string
	port        int
}

func NewTestSocket() *TestSocket {
	return &TestSocket{
		readBuffer:  bytes.NewBufferString(""),
		writeBuffer: bytes.NewBufferString(""),
	}
}

func (this *TestSocket) Close() error {
	this.close++
	return nil
}
func (this *TestSocket) CloseRead() error {
	this.closeRead++
	return nil
}
func (this *TestSocket) CloseWrite() error {
	this.closeWrite++
	return nil
}
func (this *TestSocket) Read(value []byte) (int, error) {
	this.reads++

	if this.close > 0 || this.closeRead > 0 {
		panic("Can't read from closed socket")
	}

	return this.readBuffer.Read(value)
}
func (this *TestSocket) Write(value []byte) (int, error) {
	this.writes++

	if this.close > 0 || this.closeWrite > 0 {
		panic("Can't write to closed socket")
	}
	return this.writeBuffer.Write(value)
}

func (this *TestSocket) RemoteAddr() net.Addr {
	return &net.TCPAddr{
		IP:   net.ParseIP(this.address),
		Port: this.port,
	}
}
