package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReportMetrics(t *testing.T) {
	requestPaths := make(map[string]bool)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		requestPaths[r.URL.Path] = true
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	metrics.Gauges["Alloc"] = 123.456
	metrics.Counters["PollCount"] = 1

	reportMetrics(testServer.URL)

	expectedPaths := map[string]bool{
		"/update/gauge/Alloc/123.456": true,
		"/update/counter/PollCount/1": true,
	}

	for path := range expectedPaths {
		if !requestPaths[path] {
			t.Errorf("Expected URL path %s not found", path)
		}
	}
	for path := range requestPaths {
		if !expectedPaths[path] {
			t.Errorf("Unexpected URL path %s", path)
		}
	}
}
