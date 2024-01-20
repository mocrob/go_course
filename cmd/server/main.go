package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/mocrob/go_course.git/internal/agent"
	"github.com/mocrob/go_course.git/internal/handler"
	"github.com/mocrob/go_course.git/internal/storage"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var sugar zap.SugaredLogger

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}(logger)

	sugar = *logger.Sugar()
	parseFlags()
	memoryStorage := storage.NewMemoryStorage()
	if flagRestore {
		memoryStorage = storage.NewMemoryStorageFromFile(filepath.Join(rootDir, flagFileStoragePath))
	} else {
		memoryStorage = storage.NewMemoryStorage()
	}

	r := chi.NewRouter()
	r.Post("/update/", handler.WithGzip(handler.WithLogging(handler.MetricUpdateHandler(memoryStorage), sugar)))
	r.Post("/update/{type}/{name}/{value}", handler.WithGzip(handler.WithLogging(handler.MetricUpdateHandler(memoryStorage), sugar)))
	r.Get("/value/{type}/{name}", handler.WithGzip(handler.WithLogging(handler.MetricGetHandler(memoryStorage), sugar)))
	r.Post("/value/", handler.WithGzip(handler.WithLogging(handler.MetricPostHandler(memoryStorage), sugar)))
	r.Get("/", handler.WithGzip(handler.WithLogging(handler.MetricGetAllHandler(memoryStorage), sugar)))
	sugar.Infow(
		"Starting server",
		"addr", flagRunAddr,
	)
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
	g.Go(func() error {
		// Запускаем агент с использованием контекста
		return agent.SaveMetricsInFileAgent(memoryStorage, filepath.Join(rootDir, flagFileStoragePath), time.Duration(flagStoreInterval), gCtx)
	})
	if err := g.Wait(); err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
		fmt.Printf("exit reason: %s \n", err)
	}
}
