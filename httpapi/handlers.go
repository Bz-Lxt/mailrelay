package httpapi

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Bz-Lxt/mailrelay/types"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]string{"status": "ok"})
}

func (s *Server) handleEnqueue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, 405, map[string]string{"error": "method"})
		return
	}
	var env types.Envelope
	if err := json.NewDecoder(r.Body).Decode(&env); err != nil && err != io.EOF {
		writeJSON(w, 400, map[string]string{"error": err.Error()})
		return
	}
	got, err := s.relay.Enqueue(r.Context(), env)
	if err != nil {
		writeJSON(w, statusOf(err), map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 201, got)
}

func (s *Server) handleDispatch(w http.ResponseWriter, r *http.Request) {
	n, _ := strconv.Atoi(r.URL.Query().Get("n"))
	got, err := s.relay.Dispatch(r.Context(), n)
	if err != nil {
		writeJSON(w, statusOf(err), map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, got)
}

func (s *Server) handleList(w http.ResponseWriter, r *http.Request) {
	box := r.URL.Query().Get("box")
	got, err := s.relay.List(r.Context(), box)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, got)
}

func (s *Server) handleOne(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/mail/")
	got, err := s.relay.Get(r.Context(), id)
	if err != nil {
		writeJSON(w, statusOf(err), map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, got)
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	mailboxes, err := s.relay.Report(r.Context())
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{
		"metrics":   s.relay.Stats(),
		"mailboxes": mailboxes,
	})
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	p := filepath.Join(s.web, "index.html")
	if _, err := os.Stat(p); err == nil {
		http.ServeFile(w, r, p)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(`<!doctype html><html><body><h1>mailrelay</h1><p>ops</p></body></html>`))
}

func statusOf(err error) int {
	switch err {
	case types.ErrInvalid:
		return 400
	case types.ErrNotFound:
		return 404
	case types.ErrQuota, types.ErrConflict:
		return 409
	case types.ErrCanceled:
		return 499
	default:
		return 500
	}
}
