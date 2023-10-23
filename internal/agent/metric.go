package agent

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/mocrob/go_course.git/internal/entity"
	"github.com/mocrob/go_course.git/internal/repository"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

func MetricAgent(repo repository.MetricRepo, hookPath string, reportInterval time.Duration, pollInterval time.Duration, stopSymb <-chan struct{}) {
	ticker := time.NewTicker(pollInterval * time.Second)
	sendTicker := time.NewTicker(reportInterval * time.Second)
	for {
		select {
		case <-ticker.C:
			collectGaugeMetrics(repo)
			increasePollIteration(repo)
		case <-sendTicker.C:
			sendMetrics(repo, hookPath)
			err := repo.ClearMetrics()
			if err != nil {
				log.Fatal(err)
			}
		case <-stopSymb:
			return
		}
	}
}

func collectGaugeMetrics(repo repository.MetricRepo) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	addGaugeMetricToStorage("Alloc", float64(m.Alloc), repo)
	addGaugeMetricToStorage("BuckHashSys", float64(m.BuckHashSys), repo)
	addGaugeMetricToStorage("Frees", float64(m.Frees), repo)
	addGaugeMetricToStorage("GCCPUFraction", float64(m.GCCPUFraction), repo)
	addGaugeMetricToStorage("GCSys", float64(m.GCSys), repo)
	addGaugeMetricToStorage("HeapAlloc", float64(m.HeapAlloc), repo)
	addGaugeMetricToStorage("HeapIdle", float64(m.HeapIdle), repo)
	addGaugeMetricToStorage("HeapInuse", float64(m.HeapInuse), repo)
	addGaugeMetricToStorage("HeapObjects", float64(m.HeapObjects), repo)
	addGaugeMetricToStorage("HeapReleased", float64(m.HeapReleased), repo)
	addGaugeMetricToStorage("HeapSys", float64(m.HeapSys), repo)
	addGaugeMetricToStorage("LastGC", float64(m.LastGC), repo)
	addGaugeMetricToStorage("Lookups", float64(m.Lookups), repo)
	addGaugeMetricToStorage("MCacheInuse", float64(m.MCacheInuse), repo)
	addGaugeMetricToStorage("MCacheSys", float64(m.MCacheSys), repo)
	addGaugeMetricToStorage("MSpanInuse", float64(m.MSpanInuse), repo)
	addGaugeMetricToStorage("MSpanSys", float64(m.MSpanSys), repo)
	addGaugeMetricToStorage("Mallocs", float64(m.Mallocs), repo)
	addGaugeMetricToStorage("NextGC", float64(m.NextGC), repo)
	addGaugeMetricToStorage("NumForcedGC", float64(m.NumForcedGC), repo)
	addGaugeMetricToStorage("NumGC", float64(m.NumGC), repo)
	addGaugeMetricToStorage("OtherSys", float64(m.OtherSys), repo)
	addGaugeMetricToStorage("PauseTotalNs", float64(m.PauseTotalNs), repo)
	addGaugeMetricToStorage("StackInuse", float64(m.StackInuse), repo)
	addGaugeMetricToStorage("StackSys", float64(m.StackSys), repo)
	addGaugeMetricToStorage("Sys", float64(m.Sys), repo)
	addGaugeMetricToStorage("TotalAlloc", float64(m.TotalAlloc), repo)
	addGaugeMetricToStorage("RandomValue", rand.Float64(), repo)
}

func addGaugeMetricToStorage(name string, value float64, repo repository.MetricRepo) {
	metric := entity.Metric{
		Type:  entity.Gauge,
		Name:  name,
		Value: value,
	}

	err := repo.AddMetric(name, metric)
	if err != nil {
		log.Fatal(err)
	}
}

func increasePollIteration(repo repository.MetricRepo) {
	addCounterMetricToStorage("PollCount", 1, repo)
}

func addCounterMetricToStorage(name string, value int64, repo repository.MetricRepo) {
	metric := entity.Metric{
		Type:  entity.Counter,
		Name:  name,
		Value: value,
	}

	err := repo.AddMetric(name, metric)
	if err != nil {
		log.Fatal(err)
	}
}

func sendMetrics(repo repository.MetricRepo, hookPath string) {
	allMetrics, err := repo.GetAllMetrics()
	if err != nil {
		log.Fatal(err)
	}

	for _, metric := range allMetrics {
		url := fmt.Sprintf("%s/%s/%s/%v", hookPath, metric.Type, metric.Name, metric.Value)
		req := resty.New().R()
		req.Method = http.MethodPost
		req.URL = url
		res, err := req.Send()
		if err != nil {
			fmt.Println("Failed to send metric:", metric, "Error:", err)
		}
		if res.StatusCode() != http.StatusOK {
			fmt.Println("Failed to send metric: ", metric, "Wrong response code: ", res.StatusCode())
		}
	}
}
