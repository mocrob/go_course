package main

import (
	"github.com/mocrob/go_course.git/internal/agent"
	"github.com/mocrob/go_course.git/internal/entity"
	"github.com/mocrob/go_course.git/internal/handler"
	"github.com/mocrob/go_course.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAgent_MetricAgent(t *testing.T) {
	agentStorage := storage.NewMemoryStoragePlusMetrics(make(map[string]entity.Metric))
	serverStorage := storage.NewMemoryStoragePlusMetrics(make(map[string]entity.Metric))

	server := httptest.NewServer(handler.MetricUpdateHandler(serverStorage))
	defer server.Close()

	stopSymb := make(chan struct{})
	go agent.MetricAgent(agentStorage, server.URL+"/update", 2, 10, stopSymb)

	time.Sleep(1 * time.Second)
	close(stopSymb)

	assert.Equal(t, agentStorage, serverStorage)
}
