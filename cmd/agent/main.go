package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
	serverURL      = "http://localhost:8080"
)

type Metrics struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

var metrics Metrics

func init() {
	metrics = Metrics{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}
}

func updateMetrics() {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	metrics.Gauges["Alloc"] = float64(rtm.Alloc)
	metrics.Gauges["BuckHashSys"] = float64(rtm.BuckHashSys)
	metrics.Gauges["Frees"] = float64(rtm.Frees)
	metrics.Gauges["GCCPUFraction"] = rtm.GCCPUFraction
	metrics.Gauges["GCSys"] = float64(rtm.GCSys)
	metrics.Gauges["HeapAlloc"] = float64(rtm.HeapAlloc)
	metrics.Gauges["HeapIdle"] = float64(rtm.HeapIdle)
	metrics.Gauges["HeapInuse"] = float64(rtm.HeapInuse)
	metrics.Gauges["HeapObjects"] = float64(rtm.HeapObjects)
	metrics.Gauges["HeapReleased"] = float64(rtm.HeapReleased)
	metrics.Gauges["HeapSys"] = float64(rtm.HeapSys)
	metrics.Gauges["LastGC"] = float64(rtm.LastGC)
	metrics.Gauges["Lookups"] = float64(rtm.Lookups)
	metrics.Gauges["MCacheInuse"] = float64(rtm.MCacheInuse)
	metrics.Gauges["MCacheSys"] = float64(rtm.MCacheSys)
	metrics.Gauges["MSpanInuse"] = float64(rtm.MSpanInuse)
	metrics.Gauges["MSpanSys"] = float64(rtm.MSpanSys)
	metrics.Gauges["Mallocs"] = float64(rtm.Mallocs)
	metrics.Gauges["NextGC"] = float64(rtm.NextGC)
	metrics.Gauges["NumForcedGC"] = float64(rtm.NumForcedGC)
	metrics.Gauges["NumGC"] = float64(rtm.NumGC)
	metrics.Gauges["OtherSys"] = float64(rtm.OtherSys)
	metrics.Gauges["PauseTotalNs"] = float64(rtm.PauseTotalNs)
	metrics.Gauges["StackInuse"] = float64(rtm.StackInuse)
	metrics.Gauges["StackSys"] = float64(rtm.StackSys)
	metrics.Gauges["Sys"] = float64(rtm.Sys)
	metrics.Gauges["TotalAlloc"] = float64(rtm.TotalAlloc)

	metrics.Gauges["RandomValue"] = rand.Float64()

	metrics.Counters["PollCount"]++
}

func reportMetrics() {
	for name, value := range metrics.Gauges {
		sendMetric("gauge", name, fmt.Sprintf("%f", value))
	}

	for name, value := range metrics.Counters {
		sendMetric("counter", name, fmt.Sprintf("%d", value))
	}
}

func sendMetric(metricType, name, value string) {
	url := fmt.Sprintf("%s/update/%s/%s/%s", serverURL, metricType, name, value)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Server returned non-OK status: %s\n", resp.Status)
	}
}

func main() {
	for {
		updateMetrics()
		time.Sleep(pollInterval)
		reportMetrics()
		time.Sleep(reportInterval - pollInterval)
	}
}
