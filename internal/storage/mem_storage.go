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

func (m *MemoryStorage) AddMetric(name string, metric entity.Metric) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if existMetric, ok := m.metrics[name]; ok {
		if existMetric.Type == entity.Counter {
			existMetric.Value = existMetric.Value.(int64) + metric.Value.(int64)
		} else {
			existMetric.Value = metric.Value
		}
		m.metrics[name] = existMetric
	} else {
		m.metrics[name] = metric
	}
}

func (m *MemoryStorage) GetMetric(name string) (entity.Metric, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	metric, ok := m.metrics[name]
	return metric, ok
}

func (m *MemoryStorage) GetAllMetrics() map[string]entity.Metric {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.metrics
}

func (m *MemoryStorage) ClearMetrics() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metrics = make(map[string]entity.Metric)
}
