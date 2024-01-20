package storage

import (
	"encoding/json"
	"github.com/mocrob/go_course.git/internal/entity"
	"github.com/pkg/errors"
	"log"
	"os"
	"sync"
)

type MemoryStorage struct {
	mu      sync.RWMutex
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
func NewMemoryStorageFromFile(fileStoragePath string) *MemoryStorage {
	file, err := os.OpenFile(fileStoragePath, os.O_RDONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("%+v", errors.Wrap(err, "failed to open file"))
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}(file)

	// Получаем размер файла
	fileStat, err := file.Stat()
	if err != nil {
		log.Fatalf("%+v", errors.Wrap(err, "failed to get file stats"))
	}

	if fileStat.Size() == 0 {
		// Файл пустой, возвращаем MemStorage с пустым map
		return NewMemoryStoragePlusMetrics(make(map[string]entity.Metric))
	}

	decoder := json.NewDecoder(file)
	initialMetrics := map[string]entity.Metric{}
	if err := decoder.Decode(&initialMetrics); err != nil {
		log.Fatalf("%+v", errors.Wrap(err, "failed to decode metrics"))
	}

	return NewMemoryStoragePlusMetrics(initialMetrics)
}
func (m *MemoryStorage) AddMetric(id string, metric entity.Metric) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	existingMetric, ok := m.metrics[id]
	if !ok {
		m.metrics[id] = metric
		return nil
	}

	if existingMetric.MType == entity.Gauge {
		existingMetric.Value = metric.Value
		m.metrics[id] = existingMetric
		return nil
	}

	if metric.Delta == nil {
		return errors.New("delta cannot be nil for counter metric")
	}

	if existingMetric.Delta == nil {
		existingMetric.Delta = new(int64)
	}

	*existingMetric.Delta += *metric.Delta
	m.metrics[id] = existingMetric

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
