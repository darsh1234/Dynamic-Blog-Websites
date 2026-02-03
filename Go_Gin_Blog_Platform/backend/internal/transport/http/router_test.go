package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type fakeHealthChecker struct {
	err error
}

func (f fakeHealthChecker) Ping(_ context.Context) error {
	return f.err
}

func TestHealthz(t *testing.T) {
	r := NewRouter(slog.Default(), RouterDependencies{
		HealthChecker:      fakeHealthChecker{},
		HealthCheckTimeout: time.Second,
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/healthz", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	body, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var payload map[string]string
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("failed to parse response JSON: %v", err)
	}

	if payload["status"] != "ok" {
		t.Fatalf("expected status field to be 'ok', got %q", payload["status"])
	}
}

func TestHealthzUnhealthy(t *testing.T) {
	r := NewRouter(slog.Default(), RouterDependencies{
		HealthChecker:      fakeHealthChecker{err: errors.New("db down")},
		HealthCheckTimeout: time.Second,
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/healthz", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d", w.Code)
	}
}

func TestNotImplementedShape(t *testing.T) {
	r := NewRouter(slog.Default(), RouterDependencies{
		HealthCheckTimeout: time.Second,
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotImplemented {
		t.Fatalf("expected status 501, got %d", w.Code)
	}

	body, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("failed to parse response JSON: %v", err)
	}

	errObj, ok := payload["error"].(map[string]any)
	if !ok {
		t.Fatalf("expected error object in response")
	}

	if errObj["code"] != "not_implemented" {
		t.Fatalf("expected code not_implemented, got %v", errObj["code"])
	}
}
