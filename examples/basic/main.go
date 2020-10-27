package main

import (
	"log"
	"net/http"

	"github.com/smartystreets/cproxy/v2"
)

func main() {
	handler := cproxy.New()
	log.Println("Listening on:", "*:8080")
	_ = http.ListenAndServe(":8080", handler)
}
