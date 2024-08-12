package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Dmitriy-net/pcmetrics/internal/repository"

	"github.com/go-chi/chi/v5"
)

func UpdateMetricHandler(repo repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Printf("Invalid method: %s\n", r.Method)
			return
		}

		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "name")
		metricValue := chi.URLParam(r, "value")

		var err error

		switch metricType {
		case "gauge":
			value, parseErr := strconv.ParseFloat(metricValue, 64)
			if parseErr != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				log.Printf("Failed to parse gauge value: %s\n", parseErr)
				return
			}
			err = repo.UpdateGauge(metricName, value)
		case "counter":
			value, parseErr := strconv.ParseInt(metricValue, 10, 64)
			if parseErr != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				log.Printf("Failed to parse counter value: %s\n", parseErr)
				return
			}
			err = repo.UpdateCounter(metricName, value)
		default:
			http.Error(w, "Bad request", http.StatusBadRequest)
			log.Printf("Invalid metric type: %s\n", metricType)
			return
		}

		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Printf("Error updating metric: %s\n", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200 OK"))
		log.Printf("Metric updated successfully: %s %s %s\n", metricType, metricName, metricValue)
	}
}

func GetValueHandler(repo repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "name")

		var (
			value  interface{}
			exists bool
			err    error
		)

		switch metricType {
		case "gauge":
			value, exists, err = repo.GetGauge(metricName)
		case "counter":
			value, exists, err = repo.GetCounter(metricName)
		default:
			http.Error(w, "Bad request", http.StatusBadRequest)
			log.Printf("Invalid metric type: %s\n", metricType)
			return
		}

		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Printf("Error getting metric value: %s\n", err)
			return
		}

		if !exists {
			http.Error(w, "Not found", http.StatusNotFound)
			log.Printf("Metric not found: %s %s\n", metricType, metricName)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%v", value)))
		log.Printf("Returned metric value: %s %s = %v\n", metricType, metricName, value)
	}
}

func ListMetricsHandler(repo repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gauges, counters, err := repo.ListMetrics()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Printf("Error listing metrics: %s\n", err)
			return
		}

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
		log.Println("Listed all metrics")
	}
}
