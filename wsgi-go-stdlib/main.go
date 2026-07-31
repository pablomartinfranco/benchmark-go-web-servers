package main

import (
	"log"
	"net/http"
	"os"

	"github.com/pablomartinfranco/benchmark-go-web-servers/internal/benchmarkapp"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "18080"
	}

	addr := "127.0.0.1:" + port
	if err := http.ListenAndServe(addr, benchmarkapp.Handler()); err != nil {
		log.Fatalf("wsgi stdlib backend server failed: %v", err)
	}
}
