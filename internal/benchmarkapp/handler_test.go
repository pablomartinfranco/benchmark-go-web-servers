package benchmarkapp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerRoot(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d", rr.Code)
	}
	if got := rr.Header().Get("Content-Type"); !strings.Contains(got, "application/json") {
		t.Fatalf("unexpected content type: %q", got)
	}
	if !strings.Contains(rr.Body.String(), `"message":"Hello, World!"`) {
		t.Fatalf("unexpected body: %q", rr.Body.String())
	}
}

func TestHandlerHealthz(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d", rr.Code)
	}
	if rr.Body.String() != "ok" {
		t.Fatalf("unexpected body: %q", rr.Body.String())
	}
}
