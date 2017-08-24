package cproxy

import (
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestFilterFixture(t *testing.T) {
	gunit.Run(new(FilterFixture), t)
}

type FilterFixture struct {
	*gunit.Fixture
}

func (this *FilterFixture) TestAllowEverything() {
	filter := NewFilter()

	this.So(filter.IsAuthorized(nil), should.BeTrue)
	this.So(filter.IsAuthorized(httptest.NewRequest("GET", "/", nil)), should.BeTrue)
}
