package main

import (
	"log"
	"net/http"

	"github.com/smartystreets/cproxy"
)

func main() {
	handler := cproxy.Configure().Build()
	log.Println("Listening on:", "*:8080")
	http.ListenAndServe(":8080", handler)
}
