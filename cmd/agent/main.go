package main

import (
	"github.com/mocrob/go_course.git/internal/agent"
	"github.com/mocrob/go_course.git/internal/repository"
	"github.com/mocrob/go_course.git/internal/storage"
)

func main() {
	parseFlags()
	memStorage := storage.NewMemoryStorage()
	repo := repository.MetricRepo(memStorage)
	stopSymb := make(chan struct{})
	go agent.MetricAgent(repo, "http://"+flagRunAddr+"/update", flagReportInterval, flagPollInterval, stopSymb)
	select {}
}
