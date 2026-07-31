package main

import (
	"log"
	"net/http"
	"net/http/fcgi"
	"runtime"

	"github.com/pablomartinfranco/benchmark-go-web-servers/internal/benchmarkapp"
)

type queuedRequest struct {
	w    http.ResponseWriter
	r    *http.Request
	done chan struct{}
}

type dispatchHandler struct {
	next  http.Handler
	queue chan queuedRequest
}

func newDispatchHandler(next http.Handler, workers, queueSize int) http.Handler {
	h := &dispatchHandler{
		next:  next,
		queue: make(chan queuedRequest, queueSize),
	}

	for range workers {
		go func() {
			for req := range h.queue {
				func() {
					defer close(req.done)
					h.next.ServeHTTP(req.w, req.r)
				}()
			}
		}()
	}

	return h
}

func (h *dispatchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := queuedRequest{w: w, r: r, done: make(chan struct{})}

	select {
	case h.queue <- req:
		<-req.done
	default:
		h.next.ServeHTTP(w, r)
	}
}

func main() {
	workers := runtime.GOMAXPROCS(0) * 2
	queueSize := workers * 1024

	h := newDispatchHandler(benchmarkapp.Handler(), workers, queueSize)
	if err := fcgi.Serve(nil, h); err != nil {
		log.Fatalf("fastcgi dispatch server failed: %v", err)
	}
}
