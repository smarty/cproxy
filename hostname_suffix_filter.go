package cproxy

import (
	"net/http"
	"strings"
)

type HostnameSuffixFilter struct {
	authorized []string
}

func NewHostnameSuffixFilter(authorized []string) *HostnameSuffixFilter {
	return &HostnameSuffixFilter{authorized: authorized}
}

func (this HostnameSuffixFilter) IsAuthorized(request *http.Request) bool {
	host := request.URL.Host

	for _, authorized := range this.authorized {
		if strings.HasSuffix(host, authorized) {
			return true
		}
	}

	return false
}
