package cproxy

import "net/http"

type DefaultFilter struct{}

func NewFilter() *DefaultFilter {
	return &DefaultFilter{}
}

func (this *DefaultFilter) IsAuthorized(request *http.Request) bool {
	return true
}
