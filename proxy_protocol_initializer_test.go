package cproxy

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestProxyProtocolInitializerFixture(t *testing.T) {
	gunit.Run(new(ProxyProtocolInitializerFixture), t)
}

type ProxyProtocolInitializerFixture struct {
	*gunit.Fixture

	client      *TestSocket
	server      *TestSocket
	initializer *ProxyProtocolInitializer
}

func (it *ProxyProtocolInitializerFixture) Setup() {
	it.client = NewTestSocket()
	it.server = NewTestSocket()
	it.initializer = NewProxyProtocolInitializer()
}

func (it *ProxyProtocolInitializerFixture) TestIPv4ProtocolV1() {
	it.client.address = "1.1.1.1"
	it.client.port = 1234
	it.server.address = "2.2.2.2"
	it.server.port = 5678

	result := it.initializer.Initialize(it.client, it.server)

	it.So(result, should.BeTrue)
	it.So(it.server.writeBuffer.String(), should.Equal, "PROXY TCP4 1.1.1.1 2.2.2.2 1234 5678\r\n")
	it.So(it.client.writeBuffer.String(), should.BeEmpty)
	it.So(it.server.close, should.Equal, 0)
	it.So(it.client.close, should.Equal, 0)
}

func (it *ProxyProtocolInitializerFixture) TestIPv6ProtocolV1() {
	it.client.address = "2001:db8::68"
	it.client.port = 1234
	it.server.address = "2.2.2.2"
	it.server.port = 5678

	result := it.initializer.Initialize(it.client, it.server)

	it.So(result, should.BeTrue)
	it.So(it.server.writeBuffer.String(), should.Equal, "PROXY TCP6 2001:db8::68 2.2.2.2 1234 5678\r\n")
	it.So(it.client.writeBuffer.String(), should.BeEmpty)
	it.So(it.server.close, should.Equal, 0)
	it.So(it.client.close, should.Equal, 0)
}
