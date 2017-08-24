package cproxy

import "net/http"

type HostnameFilter struct {
	authorized []string
}

func NewHostnameFilter(authorized []string) *HostnameFilter {
	return &HostnameFilter{authorized: authorized}
}

func (this HostnameFilter) IsAuthorized(request *http.Request) bool {
	host := request.URL.Host

	for _, authorized := range this.authorized {
		if authorized == host {
			return true
		}
	}

	return false
}
