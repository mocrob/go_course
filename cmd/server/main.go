package main

import (
	"github.com/mocrob/go_course.git/internal/handler"
	"github.com/mocrob/go_course.git/internal/repository"
	"github.com/mocrob/go_course.git/internal/storage"
	"net/http"
)

func main() {
	memStorage := storage.NewMemoryStorage()
	repo := repository.MetricRepo(memStorage)

	http.Handle("/update/", handler.MetricUpdateHandler(repo))
	http.ListenAndServe(":8080", nil)
}
