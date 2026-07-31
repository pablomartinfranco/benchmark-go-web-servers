package main

import (
	"log"
	"net/http/fcgi"

	"github.com/pablomartinfranco/benchmark-go-web-servers/internal/benchmarkapp"
)

func main() {
	if err := fcgi.Serve(nil, benchmarkapp.Handler()); err != nil {
		log.Fatalf("fastcgi stdlib server failed: %v", err)
	}
}
