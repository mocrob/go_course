package agent

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/mocrob/go_course.git/internal/entity"
	"github.com/mocrob/go_course.git/internal/repository"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

func MetricAgent(repo repository.MetricRepo, hookPath string, stopSymb <-chan struct{}) {
	ticker := time.NewTicker(2 * time.Second)
	sendTicker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			gaugeMetrics := collectGaugeMetrics()
			for name, value := range gaugeMetrics {
				metric := entity.Metric{
					Type:  entity.Gauge,
					Name:  name,
					Value: value,
				}
				repo.AddMetric(name, metric)
			}

			counterMetrics := collectCounterMetrics()
			for name, value := range counterMetrics {
				metric := entity.Metric{
					Type:  entity.Counter,
					Name:  name,
					Value: value,
				}
				repo.AddMetric(name, metric)
			}
		case <-sendTicker.C:
			sendMetrics(repo, hookPath)
			repo.ClearMetrics()
		case <-stopSymb:
			return
		}
	}
}

func collectGaugeMetrics() map[string]float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return map[string]float64{
		"Alloc":         float64(m.Alloc),
		"BuckHashSys":   float64(m.BuckHashSys),
		"Frees":         float64(m.Frees),
		"GCCPUFraction": m.GCCPUFraction,
		"GCSys":         float64(m.GCSys),
		"HeapAlloc":     float64(m.HeapAlloc),
		"HeapIdle":      float64(m.HeapIdle),
		"HeapInuse":     float64(m.HeapInuse),
		"HeapObjects":   float64(m.HeapObjects),
		"HeapReleased":  float64(m.HeapReleased),
		"HeapSys":       float64(m.HeapSys),
		"LastGC":        float64(m.LastGC),
		"Lookups":       float64(m.Lookups),
		"MCacheInuse":   float64(m.MCacheInuse),
		"MCacheSys":     float64(m.MCacheSys),
		"MSpanInuse":    float64(m.MSpanInuse),
		"MSpanSys":      float64(m.MSpanSys),
		"Mallocs":       float64(m.Mallocs),
		"NextGC":        float64(m.NextGC),
		"NumForcedGC":   float64(m.NumForcedGC),
		"NumGC":         float64(m.NumGC),
		"OtherSys":      float64(m.OtherSys),
		"PauseTotalNs":  float64(m.PauseTotalNs),
		"StackInuse":    float64(m.StackInuse),
		"StackSys":      float64(m.StackSys),
		"Sys":           float64(m.Sys),
		"TotalAlloc":    float64(m.TotalAlloc),
		"RandomValue":   rand.Float64(),
	}
}

func collectCounterMetrics() map[string]int64 {
	return map[string]int64{
		"PollCount": 1,
	}
}

func sendMetrics(repo repository.MetricRepo, hookPath string) {
	allMetrics := repo.GetAllMetrics()
	for _, metric := range allMetrics {
		url := fmt.Sprintf("%s/%s/%s/%v", hookPath, metric.Type, metric.Name, metric.Value)
		req := resty.New().R()
		req.Method = http.MethodPost
		req.URL = url
		_, err := req.Send()
		if err != nil {
			fmt.Println("Failed to send metric:", metric, "Error:", err)
		}
	}
}
