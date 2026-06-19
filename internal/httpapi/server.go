package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

type Server struct {
	version  string
	hostname string
	started  time.Time
	requests atomic.Uint64
}

func New(version string) http.Handler {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	server := &Server{
		version:  version,
		hostname: hostname,
		started:  time.Now().UTC(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", server.handleIndex)
	mux.HandleFunc("GET /healthz", server.handleHealthz)
	mux.HandleFunc("GET /readyz", server.handleReadyz)
	mux.HandleFunc("GET /version", server.handleVersion)
	mux.HandleFunc("GET /metrics", server.handleMetrics)

	return server.countRequests(mux)
}

func (s *Server) countRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.requests.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"name":     "k3s-sample",
		"version":  s.version,
		"hostname": s.hostname,
		"time":     time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	respondText(w, http.StatusOK, "ok\n")
}

func (s *Server) handleReadyz(w http.ResponseWriter, r *http.Request) {
	respondText(w, http.StatusOK, "ready\n")
}

func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"version": s.version,
	})
}

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(s.started).Seconds()

	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "# HELP k3s_sample_http_requests_total Total HTTP requests handled.\n")
	fmt.Fprintf(w, "# TYPE k3s_sample_http_requests_total counter\n")
	fmt.Fprintf(w, "k3s_sample_http_requests_total %d\n", s.requests.Load())
	fmt.Fprintf(w, "# HELP k3s_sample_uptime_seconds Process uptime in seconds.\n")
	fmt.Fprintf(w, "# TYPE k3s_sample_uptime_seconds gauge\n")
	fmt.Fprintf(w, "k3s_sample_uptime_seconds %.0f\n", uptime)
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func respondText(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}
