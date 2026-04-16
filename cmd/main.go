package main

import (
	_ "embed"
	"fmt"
	"log"
	"net"

	"github.com/SaranHiruthikM/performant-tcp/internals"
	"github.com/SaranHiruthikM/performant-tcp/internals/config"
	"github.com/SaranHiruthikM/performant-tcp/internals/metrics"
	ratelimiter "github.com/SaranHiruthikM/performant-tcp/internals/rateLimiter"
)

func main() {
	cfg := config.Load()
	requestsProcessed, requestsRateLimited := metrics.RunMetrics(cfg.Metrics.Path, cfg.Metrics.Port)
	workers := internals.NewWorkerPool(cfg.Server.Workers, int(cfg.Server.QueueSize), requestsProcessed)
	rateLimiter := ratelimiter.NewTokenBucket(cfg.Server.TokenRate, cfg.Server.TokenLimit)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		log.Fatal("Listener failed..")
	}
	log.Printf("Server started at :%d", cfg.Server.Port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Client not connected")
			continue
		}

		if !rateLimiter.Allow() {
			requestsRateLimited.Inc()
			conn.Write([]byte("HTTP/1.1 429 Too Many Requests\r\n\r\nRate limit exceeded\n"))
			conn.Close()
			continue
		}

		workers.Submit(internals.Job{Conn: conn})
	}
}
