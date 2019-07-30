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

func (it *FilterFixture) TestAllowEverything() {
	filter := NewFilter()

	it.So(filter.IsAuthorized(nil), should.BeTrue)
	it.So(filter.IsAuthorized(httptest.NewRequest("GET", "/", nil)), should.BeTrue)
}
