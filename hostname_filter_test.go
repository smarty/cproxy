package cproxy

import (
	"net/http/httptest"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestHostnameFilterFixture(t *testing.T) {
	gunit.Run(new(HostnameFilterFixture), t)
}

type HostnameFilterFixture struct {
	*gunit.Fixture

	filter Filter
}

func (this *HostnameFilterFixture) Setup() {
	this.filter = NewHostnameFilter([]string{"domain:1234", "sub.domain2:5678", "*.xyz.domain3:9012"})
}

func (this *HostnameFilterFixture) TestDenied() {
	this.assertUnauthorized("")
	this.assertUnauthorized("a")
	this.assertUnauthorized("domain:12345")
	this.assertUnauthorized("DOMAIN:1234")
	this.assertUnauthorized("whatever.domain2:5678")
	this.assertUnauthorized("whatever.sub.domain2:5678")
	this.assertUnauthorized("somedomain:1234") // must match exactly
	this.assertUnauthorized("domain3:1234")
	this.assertUnauthorized("xyz.domain3:9012")
}

func (this *HostnameFilterFixture) assertUnauthorized(domain string) {
	request := httptest.NewRequest("CONNECT", domain, nil)
	this.So(this.filter.IsAuthorized(nil, request), should.BeFalse)
}

func (this *HostnameFilterFixture) TestAuthorized() {
	this.assertAuthorized("domain:1234")
	this.assertAuthorized("sub.domain2:5678")
	this.assertAuthorized("a.xyz.domain3:9012")
	this.assertAuthorized("b.xyz.domain3:9012")
}

func (this *HostnameFilterFixture) assertAuthorized(domain string) {
	request := httptest.NewRequest("CONNECT", domain, nil)
	this.So(this.filter.IsAuthorized(nil, request), should.BeTrue)
}
