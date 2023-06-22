package cproxy

import (
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestProxyProtocolInitializerFixture(t *testing.T) {
	gunit.Run(new(ProxyProtocolInitializerFixture), t)
}

type ProxyProtocolInitializerFixture struct {
	*gunit.Fixture

	client      *TestSocket
	server      *TestSocket
	initializer *proxyProtocolInitializer
}

func (this *ProxyProtocolInitializerFixture) Setup() {
	this.client = NewTestSocket()
	this.server = NewTestSocket()
	this.initializer = newProxyProtocolInitializer()
}

func (this *ProxyProtocolInitializerFixture) TestIPv4ProtocolV1() {
	this.client.address = "1.1.1.1"
	this.client.port = 1234
	this.server.address = "2.2.2.2"
	this.server.port = 5678

	result := this.initializer.Initialize(this.client, this.server)

	this.So(result, should.BeTrue)
	this.So(this.server.writeBuffer.String(), should.Equal, "PROXY TCP4 1.1.1.1 2.2.2.2 1234 5678\r\n")
	this.So(this.client.writeBuffer.String(), should.BeEmpty)
	this.So(this.server.close, should.Equal, 0)
	this.So(this.client.close, should.Equal, 0)
}

func (this *ProxyProtocolInitializerFixture) TestIPv6ProtocolV1() {
	this.client.address = "2001:db8::68"
	this.client.port = 1234
	this.server.address = "2.2.2.2"
	this.server.port = 5678

	result := this.initializer.Initialize(this.client, this.server)

	this.So(result, should.BeTrue)
	this.So(this.server.writeBuffer.String(), should.Equal, "PROXY TCP6 2001:db8::68 2.2.2.2 1234 5678\r\n")
	this.So(this.client.writeBuffer.String(), should.BeEmpty)
	this.So(this.server.close, should.Equal, 0)
	this.So(this.client.close, should.Equal, 0)
}
