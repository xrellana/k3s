package httpapi

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	version        string
	hostname       string
	started        time.Time
	requests       prometheus.Counter
	metricsHandler http.Handler
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
		requests: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "k3s_sample",
			Name:      "http_requests_total",
			Help:      "Total HTTP requests handled.",
		}),
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(server.requests)
	registry.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace: "k3s_sample",
		Name:      "uptime_seconds",
		Help:      "Process uptime in seconds.",
	}, func() float64 {
		return time.Since(server.started).Seconds()
	}))
	server.metricsHandler = promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", server.handleIndex)
	mux.HandleFunc("GET /healthz", server.handleHealthz)
	mux.HandleFunc("GET /readyz", server.handleReadyz)
	mux.HandleFunc("GET /version", server.handleVersion)
	mux.Handle("GET /metrics", server.metricsHandler)

	return server.countRequests(mux)
}

func (s *Server) countRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.requests.Inc()
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
