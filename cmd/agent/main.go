package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/Dmitriy-net/pcmetrics/internal/logger"
)

const (
	defaultPollInterval   = 2 * time.Second
	defaultReportInterval = 10 * time.Second
	defaultServerURL      = "http://localhost:8080"
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

func reportMetrics(serverURL string) {
	for name, value := range metrics.Gauges {
		sendMetric(serverURL, "gauge", name, fmt.Sprintf("%g", value))
	}

	for name, value := range metrics.Counters {
		sendMetric(serverURL, "counter", name, fmt.Sprintf("%d", value))
	}
}

func sendMetric(serverURL, metricType, name, value string) {
	url := fmt.Sprintf("%s/update/%s/%s/%s", serverURL, metricType, name, value)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		log.Printf("Error creating request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Server returned non-OK status: %s\n", resp.Status)
		log.Printf("Server returned non-OK status: %s\n", resp.Status)
	}
}

func main() {
	time.Sleep(5000)

	logFile := logger.SetupLogging("agent.log")
	defer logFile.Close()

	address := flag.String("a", defaultServerURL, "HTTP server address")
	reportInterval := flag.Int("r", 10, "Report interval in seconds")
	pollInterval := flag.Int("p", 2, "Poll interval in seconds")

	flag.Parse()

	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		*address = envAddress
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		value, err := strconv.Atoi(envReportInterval)
		if err != nil {
			fmt.Printf("Invalid REPORT_INTERVAL value: %v\n", err)
			log.Printf("Invalid REPORT_INTERVAL value: %v\n", err)
		} else {
			*reportInterval = value
		}
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		value, err := strconv.Atoi(envPollInterval)
		if err != nil {
			fmt.Printf("Invalid POLL_INTERVAL value: %v\n", err)
			log.Printf("Invalid POLL_INTERVAL value: %v\n", err)
		} else {
			*pollInterval = value
		}
	}

	serverURL := *address
	pollIntervalDuration := time.Duration(*pollInterval) * time.Second
	reportIntervalDuration := time.Duration(*reportInterval) * time.Second
	fmt.Printf("Send requests to %s\n", *address)
	log.Printf("Send requests to %s\n", *address)

	for {
		updateMetrics()
		time.Sleep(pollIntervalDuration)

		reportMetrics(serverURL)
		time.Sleep(reportIntervalDuration - pollIntervalDuration)
	}
}
