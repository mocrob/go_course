package main

import (
	"github.com/mocrob/go_course.git/internal/agent"
	"github.com/mocrob/go_course.git/internal/storage"
	"time"
)

func main() {
	parseFlags()
	memStorage := storage.NewMemoryStorage()
	stopSymb := make(chan struct{})
	agent.MetricAgent(memStorage, "http://"+flagRunAddr+"/update", time.Duration(flagReportInterval), time.Duration(flagPollInterval), stopSymb)
}
