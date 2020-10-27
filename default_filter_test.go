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
	filter := newFilter()

	this.So(filter.IsAuthorized(nil, nil), should.BeTrue)
	this.So(filter.IsAuthorized(nil, httptest.NewRequest("GET", "/", nil)), should.BeTrue)
}
