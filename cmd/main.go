package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/SaranHiruthikM/performant-tcp/internals"
	"github.com/SaranHiruthikM/performant-tcp/internals/config"
	"github.com/SaranHiruthikM/performant-tcp/internals/metrics"
	ratelimiter "github.com/SaranHiruthikM/performant-tcp/internals/ratelimiter"
)

func main() {
	// load config
	cfg := config.Load()

	// start metrics
	processed, rateLimited := metrics.RunMetrics(cfg.Metrics.Path, cfg.Metrics.Port)
	_ = rateLimited // used inside server

	// create components
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool := internals.NewWorkerPool(cfg.Server.Workers, int(cfg.Server.QueueSize), processed)
	limiter := ratelimiter.NewTokenBucket(cfg.Server.TokenRate, cfg.Server.TokenLimit)
	server := internals.NewServer(listener, pool, limiter)

	// start server in background goroutine
	go server.Start()
	log.Printf("server running on :%d", cfg.Server.Port)

	// block here until Ctrl+C or SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	// signal received — shut down cleanly
	server.Shutdown()
}
