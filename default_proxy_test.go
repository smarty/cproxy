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

func (it *ProxyFixture) Setup() {
	it.client = NewTestSocket()
	it.server = NewTestSocket()
	it.proxy = NewProxy(it.client, it.server)
}

func (it *ProxyFixture) TestProxyCopiesAndCloses() {
	it.client.readBuffer.WriteString("from client")
	it.server.readBuffer.WriteString("from server")

	it.proxy.Proxy()

	it.So(it.client.writeBuffer.String(), should.Equal, "from server")
	it.So(it.server.writeBuffer.String(), should.Equal, "from client")
	it.So(it.client.closeRead, should.Equal, 1)
	it.So(it.server.closeWrite, should.Equal, 1)
	it.So(it.client.closeWrite, should.Equal, 1)
	it.So(it.server.closeRead, should.Equal, 1)
	it.So(it.server.close, should.Equal, 1)
	it.So(it.client.close, should.Equal, 1)
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

func (it *TestSocket) Close() error {
	it.close++
	return nil
}
func (it *TestSocket) CloseRead() error {
	it.closeRead++
	return nil
}
func (it *TestSocket) CloseWrite() error {
	it.closeWrite++
	return nil
}
func (it *TestSocket) Read(value []byte) (int, error) {
	it.reads++

	if it.close > 0 || it.closeRead > 0 {
		panic("Can't read from closed socket")
	}

	return it.readBuffer.Read(value)
}
func (it *TestSocket) Write(value []byte) (int, error) {
	it.writes++

	if it.close > 0 || it.closeWrite > 0 {
		panic("Can't write to closed socket")
	}
	return it.writeBuffer.Write(value)
}

func (it *TestSocket) RemoteAddr() net.Addr {
	return &net.TCPAddr{
		IP:   net.ParseIP(it.address),
		Port: it.port,
	}
}
