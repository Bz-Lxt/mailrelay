package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Bz-Lxt/mailrelay/engine"
)

type Server struct {
	relay *engine.Relay
	addr  string
	web   string
	http  *http.Server
}

func New(r *engine.Relay, addr, web string) *Server {
	if addr == "" {
		addr = ":8080"
	}
	if web == "" {
		web = "web"
	}
	s := &Server{relay: r, addr: addr, web: web}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/v1/enqueue", s.handleEnqueue)
	mux.HandleFunc("/v1/dispatch", s.handleDispatch)
	mux.HandleFunc("/v1/mailboxes", s.handleList)
	mux.HandleFunc("/v1/mail/", s.handleOne)
	mux.HandleFunc("/v1/stats", s.handleStats)
	mux.HandleFunc("/", s.handleIndex)
	s.http = &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	return s
}

func (s *Server) ListenAndServe() error { return s.http.ListenAndServe() }

func (s *Server) Shutdown(ctx context.Context) error { return s.http.Shutdown(ctx) }

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
