package cproxy

import (
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestHostnameFilterFixture(t *testing.T) {
	gunit.Run(new(HostnameFilterFixture), t)
}

type HostnameFilterFixture struct {
	*gunit.Fixture

	filter *HostnameFilter
}

func (it *HostnameFilterFixture) Setup() {
	it.filter = NewHostnameFilter([]string{"domain:1234", "sub.domain2:5678", "*.xyz.domain3:9012"})
}

func (it *HostnameFilterFixture) TestDenied() {
	it.assertUnauthorized("")
	it.assertUnauthorized("a")
	it.assertUnauthorized("domain:12345")
	it.assertUnauthorized("DOMAIN:1234")
	it.assertUnauthorized("whatever.domain2:5678")
	it.assertUnauthorized("whatever.sub.domain2:5678")
	it.assertUnauthorized("somedomain:1234") // must match exactly
	it.assertUnauthorized("domain3:1234")
	it.assertUnauthorized("xyz.domain3:9012")
}

func (it *HostnameFilterFixture) assertUnauthorized(domain string) {
	request := httptest.NewRequest("CONNECT", domain, nil)
	it.So(it.filter.IsAuthorized(request), should.BeFalse)
}

func (it *HostnameFilterFixture) TestAuthorized() {
	it.assertAuthorized("domain:1234")
	it.assertAuthorized("sub.domain2:5678")
	it.assertAuthorized("a.xyz.domain3:9012")
	it.assertAuthorized("b.xyz.domain3:9012")
}

func (it *HostnameFilterFixture) assertAuthorized(domain string) {
	request := httptest.NewRequest("CONNECT", domain, nil)
	it.So(it.filter.IsAuthorized(request), should.BeTrue)
}
