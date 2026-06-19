package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthz(t *testing.T) {
	handler := New("test")

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
	if response.Body.String() != "ok\n" {
		t.Fatalf("expected ok response, got %q", response.Body.String())
	}
}

func TestVersion(t *testing.T) {
	handler := New("test-version")

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/version", nil)
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
	if !strings.Contains(response.Body.String(), `"version":"test-version"`) {
		t.Fatalf("expected version payload, got %q", response.Body.String())
	}
}

func TestUnknownPath(t *testing.T) {
	handler := New("test")

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/missing", nil)
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}
}

func TestMetrics(t *testing.T) {
	handler := New("test")

	for range 3 {
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		handler.ServeHTTP(response, request)
	}

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
	if !strings.Contains(response.Body.String(), "k3s_sample_http_requests_total 4") {
		t.Fatalf("expected request counter, got %q", response.Body.String())
	}
}
