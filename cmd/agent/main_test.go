package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateMetrics(t *testing.T) {
	oldMetrics := metrics

	updateMetrics()

	for _, key := range []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc", "RandomValue"} {
		if _, ok := metrics.Gauges[key]; !ok {
			t.Errorf("Metric %s not found", key)
		}
	}

	if metrics.Counters["PollCount"] != oldMetrics.Counters["PollCount"]+1 {
		t.Errorf("PollCount should be incremented")
	}
}

func TestSendMetric(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		expectedPath := "/update/gauge/testName/1.234"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		if r.Header.Get("Content-Type") != "text/plain" {
			t.Errorf("Expected Content-Type text/plain, got %s", r.Header.Get("Content-Type"))
		}

		w.WriteHeader(http.StatusOK)
	}))

	defer ts.Close()

	sendMetric(ts.URL, "gauge", "testName", "1.234")
}

func TestReportMetrics(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/update/gauge/RandomValue/1.234" {
			t.Errorf("Expected path /update/gauge/RandomValue/1.234, got %s", r.URL.Path)
		}

		if r.Header.Get("Content-Type") != "text/plain" {
			t.Errorf("Expected Content-Type text/plain, got %s", r.Header.Get("Content-Type"))
		}

		w.WriteHeader(http.StatusOK)
	}))

	defer ts.Close()

	metrics.Gauges["RandomValue"] = 1.234
	metrics.Counters["TestCounter"] = 42

	reportMetrics(ts.URL)
}
