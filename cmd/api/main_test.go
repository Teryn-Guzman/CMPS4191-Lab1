package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testApplication() *applicationDependencies {
	return &applicationDependencies{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

func TestHealthcheck(t *testing.T) {
	app := testApplication()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	res := httptest.NewRecorder()

	routes(app).ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["status"] != "available" {
		t.Fatalf("expected available health status, got %q", body["status"])
	}
}

// TestReportsRoute verifies the versioned reports route exists and returns JSON.
func TestReportsRoute(t *testing.T) {
	app := testApplication()
	req := httptest.NewRequest(http.MethodPost, "/v1/reports", bytes.NewBufferString(""))
	rr := httptest.NewRecorder()

	routes(app).ServeHTTP(rr, req)

	if rr.Code == http.StatusNotFound {
		t.Fatalf("expected /v1/reports route to exist, got status %d", rr.Code)
	}
	if rr.Body.Len() == 0 {
		t.Fatal("expected /v1/reports to return a non-empty response body")
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected JSON response, got Content-Type %q", ct)
	}
}

func TestReportsRejectsInvalidRequest(t *testing.T) {
	app := testApplication()
	req := httptest.NewRequest(http.MethodPost, "/reports", strings.NewReader(`{"consumer_id":"","from":"2026-01-01T00:00:00Z","to":"2027-01-01T00:00:00Z"}`))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	routes(app).ServeHTTP(res, req)

	if res.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, res.Code)
	}
	if !strings.Contains(res.Body.String(), "consumer_id") {
		t.Fatalf("expected validation error for consumer_id, got %s", res.Body.String())
	}
}

func TestReportsRejectsUnsupportedMethod(t *testing.T) {
	app := testApplication()
	req := httptest.NewRequest(http.MethodGet, "/reports", nil)
	res := httptest.NewRecorder()

	routes(app).ServeHTTP(res, req)

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, res.Code)
	}
}

func TestReadUUIDParam(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
		fail bool
	}{
		{name: "id", path: "/consumers/abc", want: "abc"},
		{name: "trailing slash", path: "/consumers/abc/", want: "abc"},
		{name: "missing id", path: "/consumers/", fail: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			got, err := readUUIDParam("id", req)
			if tt.fail {
				if err == nil {
					t.Fatal("expected an error")
				}
				return
			}
			if err != nil {
				t.Fatalf("readUUIDParam returned error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
