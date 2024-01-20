package handler

import (
	"github.com/mocrob/go_course.git/internal/entity"
)

type MetricRepo interface {
	AddMetric(name string, metric entity.Metric) error
	GetMetric(name string) (entity.Metric, bool, error)
	GetAllMetrics() (map[string]entity.Metric, error)
	ClearMetrics() error
}
