package cproxy

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/logging"
)

func TestLoggingInitializerFixture(t *testing.T) {
	gunit.Run(new(LoggingInitializerFixture), t)
}

type LoggingInitializerFixture struct {
	*gunit.Fixture

	client      *TestSocket
	server      *TestSocket
	fakeInner   *TestInitializer
	initializer *LoggingInitializer
}

func (it *LoggingInitializerFixture) Setup() {
	it.client = NewTestSocket()
	it.server = NewTestSocket()
	it.client.address = "1.2.3.4"
	it.client.port = 4321
	it.server.address = "5.6.7.8"
	it.server.port = 8765

	it.fakeInner = NewTestInitializer(true)
	it.initializer = NewLoggingInitializer(it.fakeInner)
	it.initializer.logger = logging.Capture()
}

func (it *LoggingInitializerFixture) TestInnerInitializerCalled() {
	result := it.initializer.Initialize(it.client, it.server)

	it.So(result, should.BeTrue)
	it.So(it.fakeInner.client, should.Equal, it.client)
	it.So(it.fakeInner.server, should.Equal, it.server)
}

func (it *LoggingInitializerFixture) TestLoggingOnFailure() {
	it.fakeInner.success = false

	it.initializer.Initialize(it.client, it.server)

	it.So(it.initializer.logger.Calls, should.Equal, 1)
	it.So(it.initializer.logger.Log.String(), should.EndWith,
		"[INFO] Connection failed [1.2.3.4:4321] -> [5.6.7.8:8765]\n")
}

func (it *LoggingInitializerFixture) TestLoggingOnSuccess() {
	it.initializer.Initialize(it.client, it.server)

	it.So(it.initializer.logger.Calls, should.Equal, 1)
	it.So(it.initializer.logger.Log.String(), should.EndWith,
		"[INFO] Established connection [1.2.3.4:4321] -> [5.6.7.8:8765]\n")
}
