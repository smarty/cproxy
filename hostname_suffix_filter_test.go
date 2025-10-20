package cproxy

import (
	"net/http/httptest"
	"testing"

	"github.com/smarty/gunit"
	"github.com/smarty/gunit/assert/should"
)

func TestHostnameSuffixFilterFixture(t *testing.T) {
	gunit.Run(new(HostnameSuffixFilterFixture), t)
}

type HostnameSuffixFilterFixture struct {
	*gunit.Fixture

	filter Filter
}

func (this *HostnameSuffixFilterFixture) Setup() {
	this.filter = NewHostnameSuffixFilter([]string{"domain:1234", ".domain2:5678"})
}

func (this *HostnameSuffixFilterFixture) TestDenied() {
	this.assertUnauthorized("")
	this.assertUnauthorized("a")
	this.assertUnauthorized("domain:12345")
	this.assertUnauthorized("DOMAIN:1234")
}

func (this *HostnameSuffixFilterFixture) assertUnauthorized(domain string) {
	request := httptest.NewRequest("CONNECT", domain, nil)
	this.So(this.filter.IsAuthorized(nil, request), should.BeFalse)
}

func (this *HostnameSuffixFilterFixture) TestAuthorized() {
	this.assertAuthorized("domain:1234")
	this.assertAuthorized("whatever.domain2:5678")
	this.assertAuthorized("somedomain:1234") // only matches suffix
}

func (this *HostnameSuffixFilterFixture) assertAuthorized(domain string) {
	request := httptest.NewRequest("CONNECT", domain, nil)
	this.So(this.filter.IsAuthorized(nil, request), should.BeTrue)
}
