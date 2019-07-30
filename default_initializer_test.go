package cproxy

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestInitializerFixture(t *testing.T) {
	gunit.Run(new(InitializerFixture), t)
}

type InitializerFixture struct {
	*gunit.Fixture
}

func (it *InitializerFixture) TestAlwaysSuccessful() {
	initializer := NewInitializer()

	it.So(initializer.Initialize(nil, nil), should.BeTrue)
	it.So(initializer.Initialize(NewTestSocket(), nil), should.BeTrue)
	it.So(initializer.Initialize(NewTestSocket(), NewTestSocket()), should.BeTrue)
}
