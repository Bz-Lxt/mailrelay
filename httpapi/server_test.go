package httpapi

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Bz-Lxt/mailrelay/config"
	"github.com/Bz-Lxt/mailrelay/engine"
)

func TestHealthAndEnqueue(t *testing.T) {
	r, err := engine.Open(config.Config{Dir: t.TempDir(), Addr: ":0"})
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	s := New(r, ":0", "web")
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	s.http.Handler.ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Fatalf("health %d", rr.Code)
	}
	body := bytes.NewBufferString(`{"from":"a@local","to":"b@local","subject":"x","body":"y"}`)
	req = httptest.NewRequest(http.MethodPost, "/v1/enqueue", body)
	rr = httptest.NewRecorder()
	s.http.Handler.ServeHTTP(rr, req)
	if rr.Code != 201 {
		t.Fatalf("enqueue %d %s", rr.Code, rr.Body.String())
	}
}
