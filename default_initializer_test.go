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

func (this *InitializerFixture) TestAlwaysSuccessful() {
	initializer := newInitializer()

	this.So(initializer.Initialize(nil, nil), should.BeTrue)
	this.So(initializer.Initialize(NewTestSocket(), nil), should.BeTrue)
	this.So(initializer.Initialize(NewTestSocket(), NewTestSocket()), should.BeTrue)
}
