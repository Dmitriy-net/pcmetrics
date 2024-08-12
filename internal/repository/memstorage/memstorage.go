package memstorage

import "log"

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

func (ms *MemStorage) UpdateGauge(name string, value float64) error {
	log.Printf("Updating gauge: %s with value: %f\n", name, value)
	ms.gauges[name] = value
	return nil
}

func (ms *MemStorage) UpdateCounter(name string, value int64) error {
	log.Printf("Updating counter: %s with value: %d\n", name, value)
	ms.counters[name] += value
	return nil
}

func (ms *MemStorage) GetGauge(name string) (float64, bool, error) {
	value, exists := ms.gauges[name]
	log.Printf("Getting gauge: %s, exists: %v, value: %f\n", name, exists, value)
	return value, exists, nil
}

func (ms *MemStorage) GetCounter(name string) (int64, bool, error) {
	value, exists := ms.counters[name]
	log.Printf("Getting counter: %s, exists: %v, value: %d\n", name, exists, value)
	return value, exists, nil
}

func (ms *MemStorage) ListMetrics() (map[string]float64, map[string]int64, error) {
	log.Println("Listing all metrics")
	return ms.gauges, ms.counters, nil
}
