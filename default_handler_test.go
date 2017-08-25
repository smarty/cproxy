package cproxy

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestHandlerFixture(t *testing.T) {
	gunit.Run(new(HandlerFixture), t)
}

type HandlerFixture struct {
	*gunit.Fixture

	handler         *DefaultHandler
	filter          *TestFilter
	socket          *DummySocket
	clientConnector *TestClientConnector
	serverConnector *TestServerConnector
	meter           *TestMeter

	request  *http.Request
	response *httptest.ResponseRecorder
}

func (this *HandlerFixture) Setup() {
	this.filter = NewTestFilter(true)
	this.socket = &DummySocket{}
	this.clientConnector = NewTestClientConnector(this.socket)
	this.serverConnector = NewTestServerConnector()
	this.meter = NewTestMeter()

	this.handler = NewHandler(this.filter, this.clientConnector, this.serverConnector, this.meter)

	this.request = httptest.NewRequest("CONNECT", "domain:443", nil)
	this.response = httptest.NewRecorder()
}

//////////////////////////////////////////////////////////////

func (this *HandlerFixture) TestForbidsUnknownMethods() {
	this.request.Method = "GET"

	this.serveHTTP()

	this.shouldHaveResponse(405, "Method Not Allowed")
	this.So(this.meter.calls, should.Resemble, []int{MeasurementHTTPRequest, MeasurementBadMethod})
}

func (this *HandlerFixture) TestsDisallowsUnauthorizedRequests() {
	this.filter.authorized = false

	this.serveHTTP()

	this.So(this.filter.request, should.Equal, this.request)
	this.shouldHaveResponse(401, "Unauthorized")
	this.So(this.meter.calls, should.Resemble, []int{MeasurementHTTPRequest, MeasurementUnauthorizedRequest})
}

func (this *HandlerFixture) TestRejectClientWhichCannotBeConnected() {
	this.clientConnector.socket = nil

	this.serveHTTP()

	this.So(this.clientConnector.response, should.Equal, this.response)
	this.shouldHaveResponse(501, "Not Implemented")
	this.So(this.meter.calls, should.Resemble, []int{MeasurementHTTPRequest, MeasurementClientConnectionFailed})
}

func (this *HandlerFixture) TestRejectBadGatewayRequest() {
	this.serverConnector.proxy = nil

	this.serveHTTP()

	this.So(this.serverConnector.socket, should.Equal, this.socket)
	this.So(this.serverConnector.address, should.Equal, "domain:443")
	this.So(this.socket.String(), should.Equal, "HTTP/1.1 502 Bad Gateway\r\n\r\n")
	this.So(this.socket.closed, should.Equal, 1)
	this.shouldHaveResponse(200, "") // ResponseRecorder defaults to 200
	this.So(this.meter.calls, should.Resemble, []int{MeasurementHTTPRequest, MeasurementServerConnectionFailed})
}

func (this *HandlerFixture) TestProxyInvoked() {
	this.serveHTTP()

	this.So(this.socket.String(), should.Equal, "HTTP/1.1 200 OK\r\n\r\n")
	this.So(this.serverConnector.proxy.calls, should.Equal, 1)
	this.shouldHaveResponse(200, "") // ResponseRecorder defaults to 200
	this.So(this.meter.calls, should.Resemble, []int{MeasurementHTTPRequest, MeasurementProxyReady, MeasurementProxyComplete})
}

func (this *HandlerFixture) serveHTTP() {
	this.handler.ServeHTTP(this.response, this.request)
}
func (this *HandlerFixture) shouldHaveResponse(statusCode int, statusText string) {
	this.So(this.response.Code, should.Equal, statusCode)
	this.So(this.response.Body.String(), should.EqualTrimSpace, statusText)
}

//////////////////////////////////////////////////////////////

type TestFilter struct {
	authorized bool
	request    *http.Request
}

func NewTestFilter(authorized bool) *TestFilter {
	return &TestFilter{authorized: authorized}
}
func (this *TestFilter) IsAuthorized(request *http.Request) bool {
	this.request = request
	return this.authorized
}

//////////////////////////////////////////////////////////////

type TestClientConnector struct {
	socket   Socket
	response http.ResponseWriter
}

func NewTestClientConnector(socket Socket) *TestClientConnector {
	return &TestClientConnector{socket: socket}
}

func (this *TestClientConnector) Connect(response http.ResponseWriter) Socket {
	this.response = response
	return this.socket
}

//////////////////////////////////////////////////////////////

type TestServerConnector struct {
	socket  Socket
	address string
	proxy   *TestProxy
}

func NewTestServerConnector() *TestServerConnector {
	return &TestServerConnector{proxy: &TestProxy{}}
}

func (this *TestServerConnector) Connect(socket Socket, address string) Proxy {
	this.socket = socket
	this.address = address

	if this.proxy == nil {
		return nil // Golang nil != nil issue
	}

	return this.proxy
}

type TestProxy struct{ calls int }

func (this *TestProxy) Proxy() {
	if this != nil {
		this.calls++
	}
}

//////////////////////////////////////////////////////////////

type DummySocket struct {
	written []byte
	closed  int
}

func (this *DummySocket) Read([]byte) (int, error) {
	panic("this test shouldn't read")
}

func (this *DummySocket) Write(buffer []byte) (int, error) {
	this.written = append(this.written, buffer...)
	return len(buffer), nil
}
func (this *DummySocket) Close() error {
	this.closed++
	return nil
}

func (this *DummySocket) RemoteAddr() net.Addr {
	panic("shouldn't be called")
}

func (this *DummySocket) String() string { return string(this.written) }

//////////////////////////////////////////////////////////////

type TestMeter struct {
	calls []int
}

func NewTestMeter() *TestMeter {
	return &TestMeter{}
}
func (this *TestMeter) Measure(value int) {
	this.calls = append(this.calls, value)
}
