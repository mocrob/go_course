package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/mocrob/go_course.git/internal/handler"
	"github.com/mocrob/go_course.git/internal/repository"
	"github.com/mocrob/go_course.git/internal/storage"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()
	parseFlags()
	memoryStorage := storage.NewMemoryStorage()
	repo := repository.MetricRepo(memoryStorage)
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", handler.MetricUpdateHandler(repo))
	r.Get("/value/{type}/{name}", handler.MetricGetHandler(repo))
	r.Get("/", handler.MetricGetAllHandler(repo))
	httpServer := &http.Server{
		Addr:    flagRunAddr,
		Handler: r,
	}
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return httpServer.Shutdown(context.Background())
	})
	if err := g.Wait(); err != nil {
		fmt.Printf("exit reason: %s \n", err)
	}
}
