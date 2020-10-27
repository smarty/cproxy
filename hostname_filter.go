package cproxy

import (
	"net/http"
)

type hostnameFilter struct {
	authorized []string
}

func NewHostnameFilter(authorized []string) Filter {
	return &hostnameFilter{authorized: authorized}
}

func (this hostnameFilter) IsAuthorized(request *http.Request) bool {
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
