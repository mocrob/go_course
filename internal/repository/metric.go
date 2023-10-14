package repository

import (
	"github.com/mocrob/go_course.git/internal/entity"
)

type MetricRepo interface {
	AddMetric(name string, metric entity.Metric)
	GetMetric(name string) (entity.Metric, bool)
	GetAllMetrics() map[string]entity.Metric
	ClearMetrics()
}
