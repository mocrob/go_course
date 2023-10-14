package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/mocrob/go_course.git/internal/handler"
	"github.com/mocrob/go_course.git/internal/repository"
	"github.com/mocrob/go_course.git/internal/storage"
	"log"
	"net/http"
)

func main() {
	memoryStorage := storage.NewMemoryStorage()
	repo := repository.MetricRepo(memoryStorage)

	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", handler.MetricUpdateHandler(repo))
	r.Get("/value/{type}/{name}", handler.MetricGetHandler(repo))
	r.Get("/", handler.MetricGetAllHandler(repo))

	log.Fatal(http.ListenAndServe(":8080", r))
}
