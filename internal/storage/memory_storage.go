package storage

import (
	"github.com/mocrob/go_course.git/internal/entity"
	"sync"
)

type MemoryStorage struct {
	mu      sync.Mutex
	metrics map[string]entity.Metric
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		metrics: make(map[string]entity.Metric),
	}
}

func NewMemoryStoragePlusMetrics(initMetrics map[string]entity.Metric) *MemoryStorage {
	return &MemoryStorage{
		metrics: initMetrics,
	}
}

func (m *MemoryStorage) AddMetric(name string, metric entity.Metric) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	existingMetric, ok := m.metrics[name]
	if !ok {
		m.metrics[name] = metric
		return nil
	}

	if existingMetric.Type == entity.Counter {
		existingMetric.Value = existingMetric.Value.(int64) + metric.Value.(int64)
	} else {
		existingMetric.Value = metric.Value
	}
	m.metrics[name] = existingMetric

	return nil
}

func (m *MemoryStorage) GetMetric(name string) (entity.Metric, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	metric, ok := m.metrics[name]
	return metric, ok, nil
}

func (m *MemoryStorage) GetAllMetrics() (map[string]entity.Metric, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.metrics, nil
}

func (m *MemoryStorage) ClearMetrics() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metrics = make(map[string]entity.Metric)
	return nil
}
