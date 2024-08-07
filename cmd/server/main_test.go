package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestUpdateMetricHandler(t *testing.T) {
	memStorage := NewMemStorage()
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", updateMetricHandler(memStorage))

	tests := []struct {
		method      string
		metricType  string
		metricName  string
		metricValue string
		statusCode  int
	}{
		{"POST", "gauge", "testGauge", "12.34", http.StatusOK},
		{"POST", "counter", "testCounter", "56", http.StatusOK},
		{"POST", "invalid", "testInvalid", "123", http.StatusBadRequest},
		{"GET", "gauge", "testGauge", "12.34", http.StatusMethodNotAllowed},
		{"POST", "gauge", "testGauge", "invalidValue", http.StatusBadRequest},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, "/update/"+tt.metricType+"/"+tt.metricName+"/"+tt.metricValue, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != tt.statusCode {
			t.Errorf("handler returned wrong status code: got %v want %v",
				rr.Code, tt.statusCode)
		}
	}
}

func TestGetValueHandler(t *testing.T) {
	memStorage := NewMemStorage()
	r := chi.NewRouter()
	r.Get("/value/{type}/{name}", getValueHandler(memStorage))

	memStorage.UpdateGauge("testGauge", 12.34)
	memStorage.UpdateCounter("testCounter", 56)

	tests := []struct {
		metricType string
		metricName string
		expected   string
		statusCode int
	}{
		{"gauge", "testGauge", "12.34", http.StatusOK},
		{"counter", "testCounter", "56", http.StatusOK},
		{"gauge", "nonExistentGauge", "", http.StatusNotFound},
		{"counter", "nonExistentCounter", "", http.StatusNotFound},
		{"invalid", "testInvalid", "", http.StatusBadRequest},
	}

	for _, tt := range tests {
		req, err := http.NewRequest("GET", "/value/"+tt.metricType+"/"+tt.metricName, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != tt.statusCode {
			t.Errorf("handler returned wrong status code: got %v want %v",
				rr.Code, tt.statusCode)
		}
		if tt.statusCode == http.StatusOK && rr.Body.String() != tt.expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), tt.expected)
		}
	}
}

func TestListMetricsHandler(t *testing.T) {
	memStorage := NewMemStorage()
	r := chi.NewRouter()
	r.Get("/", listMetricsHandler(memStorage))

	memStorage.UpdateGauge("testGauge", 12.34)
	memStorage.UpdateCounter("testCounter", 56)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Gauge: testGauge = 12.34") {
		t.Errorf("handler returned unexpected body: got %v want it to contain %v",
			body, "Gauge: testGauge = 12.34")
	}
	if !strings.Contains(body, "Counter: testCounter = 56") {
		t.Errorf("handler returned unexpected body: got %v want it to contain %v",
			body, "Counter: testCounter = 56")
	}
}
