package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (ms *MemStorage) UpdateGauge(name string, value float64) {
	ms.gauges[name] = value
}

func (ms *MemStorage) UpdateCounter(name string, value int64) {
	ms.counters[name] += value
}

func (ms *MemStorage) GetGauge(name string) (float64, bool) {
	value, exists := ms.gauges[name]
	return value, exists
}

func (ms *MemStorage) GetCounter(name string) (int64, bool) {
	value, exists := ms.counters[name]
	return value, exists
}

func (ms *MemStorage) ListMetrics() (map[string]float64, map[string]int64) {
	return ms.gauges, ms.counters
}

func updateMetricHandler(ms *MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		metricType := MetricType(chi.URLParam(r, "type"))
		metricName := chi.URLParam(r, "name")
		metricValue := chi.URLParam(r, "value")

		switch metricType {
		case Gauge:
			value, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			ms.UpdateGauge(metricName, value)
		case Counter:
			value, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			ms.UpdateCounter(metricName, value)
		default:
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200 OK"))
	}
}

func getValueHandler(ms *MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := MetricType(chi.URLParam(r, "type"))
		metricName := chi.URLParam(r, "name")

		switch metricType {
		case Gauge:
			if value, exists := ms.GetGauge(metricName); exists {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprintf("%g", value)))
			} else {
				http.Error(w, "Not found", http.StatusNotFound)
			}
		case Counter:
			if value, exists := ms.GetCounter(metricName); exists {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprintf("%d", value)))
			} else {
				http.Error(w, "Not found", http.StatusNotFound)
			}
		default:
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
	}
}

func listMetricsHandler(ms *MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gauges, counters := ms.ListMetrics()
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<html><body><h1>Metrics</h1><ul>"))

		for name, value := range gauges {
			w.Write([]byte(fmt.Sprintf("<li>Gauge: %s = %g</li>", name, value)))
		}

		for name, value := range counters {
			w.Write([]byte(fmt.Sprintf("<li>Counter: %s = %d</li>", name, value)))
		}

		w.Write([]byte("</ul></body></html>"))
	}
}

func main() {
	address := flag.String("a", "localhost:8080", "HTTP server address")

	flag.Parse()

	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		*address = envAddress
	}

	r := chi.NewRouter()
	memStorage := NewMemStorage()

	r.Post("/update/{type}/{name}/{value}", updateMetricHandler(memStorage))
	r.Get("/value/{type}/{name}", getValueHandler(memStorage))
	r.Get("/", listMetricsHandler(memStorage))

	fmt.Printf("Server is listening on http://%s\n", *address)
	http.ListenAndServe(*address, r)
}
