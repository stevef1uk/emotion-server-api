package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleHealth(t *testing.T) {
	srv := &Server{}
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	srv.handleHealth(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var data map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if data["status"] != "ok" {
		t.Errorf("expected status 'ok', got %q", data["status"])
	}
}

func TestHandleIndex(t *testing.T) {
	srv := &Server{}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	srv.handleIndex(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("expected text/html, got %q", contentType)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Emotion Server API") {
		t.Errorf("expected body to contain 'Emotion Server API'")
	}
	if !strings.Contains(body, "Try It Live!") {
		t.Errorf("expected body to contain 'Try It Live!'")
	}
}

func TestHandleAnalyzeOptionsCORS(t *testing.T) {
	srv := &Server{}
	req := httptest.NewRequest(http.MethodOptions, "/api/analyze", nil)
	w := httptest.NewRecorder()

	srv.handleAnalyze(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", resp.StatusCode)
	}
	if origin := resp.Header.Get("Access-Control-Allow-Origin"); origin != "*" {
		t.Errorf("expected Access-Control-Allow-Origin: *, got %q", origin)
	}
}

func TestHandleAnalyzeEmptyBody(t *testing.T) {
	srv := &Server{}
	req := httptest.NewRequest(http.MethodPost, "/api/analyze", strings.NewReader(`{"text":""}`))
	w := httptest.NewRecorder()

	srv.handleAnalyze(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}
