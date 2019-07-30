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

func (it *HandlerFixture) Setup() {
	it.filter = NewTestFilter(true)
	it.socket = &DummySocket{}
	it.clientConnector = NewTestClientConnector(it.socket)
	it.serverConnector = NewTestServerConnector()
	it.meter = NewTestMeter()

	it.handler = NewHandler(it.filter, it.clientConnector, it.serverConnector, it.meter)

	it.request = httptest.NewRequest("CONNECT", "domain:443", nil)
	it.response = httptest.NewRecorder()
}

//////////////////////////////////////////////////////////////

func (it *HandlerFixture) TestForbidsUnknownMethods() {
	it.request.Method = "GET"

	it.serveHTTP()

	it.shouldHaveResponse(405, "Method Not Allowed")
	it.So(it.meter.calls, should.Resemble, []int{MeasurementHTTPRequest, MeasurementBadMethod})
}

func (it *HandlerFixture) TestsDisallowsUnauthorizedRequests() {
	it.filter.authorized = false

	it.serveHTTP()

	it.So(it.filter.request, should.Equal, it.request)
	it.shouldHaveResponse(401, "Unauthorized")
	it.So(it.meter.calls, should.Resemble, []int{MeasurementHTTPRequest, MeasurementUnauthorizedRequest})
}

func (it *HandlerFixture) TestRejectClientWhichCannotBeConnected() {
	it.clientConnector.socket = nil

	it.serveHTTP()

	it.So(it.clientConnector.response, should.Equal, it.response)
	it.shouldHaveResponse(501, "Not Implemented")
	it.So(it.meter.calls, should.Resemble, []int{MeasurementHTTPRequest, MeasurementClientConnectionFailed})
}

func (it *HandlerFixture) TestRejectBadGatewayRequest() {
	it.serverConnector.proxy = nil

	it.serveHTTP()

	it.So(it.serverConnector.socket, should.Equal, it.socket)
	it.So(it.serverConnector.address, should.Equal, "domain:443")
	it.So(it.socket.String(), should.Equal, "HTTP/1.1 502 Bad Gateway\r\n\r\n")
	it.So(it.socket.closed, should.Equal, 1)
	it.shouldHaveResponse(200, "") // ResponseRecorder defaults to 200
	it.So(it.meter.calls, should.Resemble, []int{MeasurementHTTPRequest, MeasurementServerConnectionFailed})
}

func (it *HandlerFixture) TestProxyInvoked() {
	it.serveHTTP()

	it.So(it.socket.String(), should.Equal, "HTTP/1.1 200 OK\r\n\r\n")
	it.So(it.serverConnector.proxy.calls, should.Equal, 1)
	it.shouldHaveResponse(200, "") // ResponseRecorder defaults to 200
	it.So(it.meter.calls, should.Resemble, []int{MeasurementHTTPRequest, MeasurementProxyReady, MeasurementProxyComplete})
}

func (it *HandlerFixture) serveHTTP() {
	it.handler.ServeHTTP(it.response, it.request)
}
func (it *HandlerFixture) shouldHaveResponse(statusCode int, statusText string) {
	it.So(it.response.Code, should.Equal, statusCode)
	it.So(it.response.Body.String(), should.EqualTrimSpace, statusText)
}

//////////////////////////////////////////////////////////////

type TestFilter struct {
	authorized bool
	request    *http.Request
}

func NewTestFilter(authorized bool) *TestFilter {
	return &TestFilter{authorized: authorized}
}
func (it *TestFilter) IsAuthorized(request *http.Request) bool {
	it.request = request
	return it.authorized
}

//////////////////////////////////////////////////////////////

type TestClientConnector struct {
	socket   Socket
	response http.ResponseWriter
}

func NewTestClientConnector(socket Socket) *TestClientConnector {
	return &TestClientConnector{socket: socket}
}

func (it *TestClientConnector) Connect(response http.ResponseWriter) Socket {
	it.response = response
	return it.socket
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

func (it *TestServerConnector) Connect(socket Socket, address string) Proxy {
	it.socket = socket
	it.address = address

	if it.proxy == nil {
		return nil // Golang nil != nil issue
	}

	return it.proxy
}

type TestProxy struct{ calls int }

func (it *TestProxy) Proxy() {
	if it != nil {
		it.calls++
	}
}

//////////////////////////////////////////////////////////////

type DummySocket struct {
	written []byte
	closed  int
}

func (it *DummySocket) Read([]byte) (int, error) {
	panic("this test shouldn't read")
}

func (it *DummySocket) Write(buffer []byte) (int, error) {
	it.written = append(it.written, buffer...)
	return len(buffer), nil
}
func (it *DummySocket) Close() error {
	it.closed++
	return nil
}

func (it *DummySocket) RemoteAddr() net.Addr {
	panic("shouldn't be called")
}

func (it *DummySocket) String() string { return string(it.written) }

//////////////////////////////////////////////////////////////

type TestMeter struct {
	calls []int
}

func NewTestMeter() *TestMeter {
	return &TestMeter{}
}
func (it *TestMeter) Measure(value int) {
	it.calls = append(it.calls, value)
}
