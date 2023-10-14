package handler

import (
	"github.com/mocrob/go_course.git/internal/entity"
	"github.com/mocrob/go_course.git/internal/repository"
	"net/http"
	"strconv"
	"strings"
)

func MetricUpdateHandler(repo repository.MetricRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/update/")
		parts := strings.Split(path, "/")

		if len(parts) != 3 {
			http.Error(w, "Некорректный URL", http.StatusNotFound)
			return
		}
		metricType := parts[0]
		if metricType != string(entity.Gauge) && metricType != string(entity.Counter) {
			http.Error(w, "Некорректный тип", http.StatusBadRequest)
			return
		}
		metricName := parts[1]
		var metricValue interface{}
		var err error
		switch metricType {
		case string(entity.Gauge):
			metricValue, err = strconv.ParseFloat(parts[2], 64)
		case string(entity.Counter):
			metricValue, err = strconv.ParseInt(parts[2], 10, 64)
		}
		if err != nil {
			http.Error(w, "Некорректное значение", http.StatusBadRequest)
			return
		}
		metric := entity.Metric{
			Type:  entity.Type(metricType),
			Name:  metricName,
			Value: metricValue,
		}
		repo.AddMetric(metricName, metric)
		w.WriteHeader(http.StatusOK)
	}
}
