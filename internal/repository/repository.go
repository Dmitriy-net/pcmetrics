package repository

type Repository interface {
	UpdateGauge(name string, value float64) error
	UpdateCounter(name string, value int64) error

	GetGauge(name string) (float64, bool, error)
	GetCounter(name string) (int64, bool, error)

	ListMetrics() (map[string]float64, map[string]int64, error)
}
