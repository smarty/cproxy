package cproxy

import (
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestHostnameSuffixFilterFixture(t *testing.T) {
	gunit.Run(new(HostnameSuffixFilterFixture), t)
}

type HostnameSuffixFilterFixture struct {
	*gunit.Fixture

	filter *HostnameSuffixFilter
}

func (it *HostnameSuffixFilterFixture) Setup() {
	it.filter = NewHostnameSuffixFilter([]string{"domain:1234", ".domain2:5678"})
}

func (it *HostnameSuffixFilterFixture) TestDenied() {
	it.assertUnauthorized("")
	it.assertUnauthorized("a")
	it.assertUnauthorized("domain:12345")
	it.assertUnauthorized("DOMAIN:1234")
}

func (it *HostnameSuffixFilterFixture) assertUnauthorized(domain string) {
	request := httptest.NewRequest("CONNECT", domain, nil)
	it.So(it.filter.IsAuthorized(request), should.BeFalse)
}

func (it *HostnameSuffixFilterFixture) TestAuthorized() {
	it.assertAuthorized("domain:1234")
	it.assertAuthorized("whatever.domain2:5678")
	it.assertAuthorized("somedomain:1234") // only matches suffix
}

func (it *HostnameSuffixFilterFixture) assertAuthorized(domain string) {
	request := httptest.NewRequest("CONNECT", domain, nil)
	it.So(it.filter.IsAuthorized(request), should.BeTrue)
}
