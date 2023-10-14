package main

import (
	"fmt"
	"github.com/mocrob/go_course.git/internal/entity/metric"
	"github.com/mocrob/go_course.git/internal/service"
	"net/http"
	"strings"
)

func main() {
	storage := metric.NewMemStorage()

	metricService := service.MetricService{}

	http.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {
		handleUpdate(w, r, storage, metricService)
	})
	http.ListenAndServe(":8080", nil)
}

func handleUpdate(w http.ResponseWriter, r *http.Request, storage *metric.MemStorage, metricService service.MetricService) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	remainingPath := strings.TrimPrefix(r.URL.Path, "/update/")
	pathSegments := strings.Split(remainingPath, "/")

	if len(pathSegments) < 3 {
		http.Error(w, "Некорректный URL", http.StatusNotFound)
		return
	}

	metricTypeName := pathSegments[0]
	metricName := pathSegments[1]
	value := pathSegments[2]

	if metricName == "" {
		http.Error(w, "Метрика не найдена", http.StatusNotFound)
		return
	}

	var err error
	switch metricTypeName {
	case "counter":
		err = metricService.AddCounterMetric(storage, metricName, value)
	case "gauge":
		err = metricService.AddGaugeMetric(storage, metricName, value)
	default:
		http.Error(w, "Тип метрики не найден", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка создания метрики: %s", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
