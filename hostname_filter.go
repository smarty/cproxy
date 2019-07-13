package cproxy

import (
	"net/http"
)

type HostnameFilter struct {
	authorized []string
}

func NewHostnameFilter(authorized []string) *HostnameFilter {
	return &HostnameFilter{authorized: authorized}
}

func (this HostnameFilter) IsAuthorized(request *http.Request) bool {
	host := request.URL.Host

	for _, authorized := range this.authorized {
		if authorized[:2] == "*." {
			have, want := len(host), len(authorized)-1
			if have > want && authorized[1:] == host[len(host)-want:] {
				return true
			}
		} else if authorized == host {
			return true
		}
	}

	return false
}
