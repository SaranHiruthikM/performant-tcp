package main

import (
	"log"
	"net"
	"sync"
	"time"
)

type Job struct {
	Conn net.Conn
}

type WorkerPool struct {
	jobs       chan Job
	maxWorkers int
	wg         sync.WaitGroup
}

type TokenBucket struct {
	MaxTokens  int64
	Tokens     int64
	Rate       int64
	LastRefill time.Time
	Mutex      sync.Mutex
}

func NewTokenBucket(rate, maxTokens int64) *TokenBucket {
	newBucket := &TokenBucket{
		MaxTokens:  maxTokens,
		Tokens:     maxTokens,
		Rate:       rate,
		LastRefill: time.Now(),
	}

	return newBucket
}

func (tb *TokenBucket) Allow() bool {
	tb.Mutex.Lock()
	defer tb.Mutex.Unlock()

	timeElapsed := time.Since(tb.LastRefill).Seconds()
	tokensToAdd := timeElapsed * float64(tb.Rate)
	tb.Tokens = min(tb.Tokens+int64(tokensToAdd), tb.MaxTokens)
	tb.LastRefill = time.Now()

	if tb.Tokens > 0 {
		tb.Tokens--
		return true
	}

	return false
}

func NewWorkerPool(maxWorkers, queueSize int) *WorkerPool {
	jobs := make(chan Job, queueSize+maxWorkers)

	newWorker := &WorkerPool{
		maxWorkers: maxWorkers,
		jobs:       jobs,
	}
	for i := range maxWorkers {
		newWorker.wg.Add(1)
		go newWorker.worker(i)
	}

	// newWorker.Close()
	return newWorker
}

func (w *WorkerPool) worker(id int) {
	defer w.wg.Done()
	for job := range w.jobs {
		conn := job.Conn
		log.Printf("Worker: %d processing job", id)
		reader := make([]byte, 1024)
		_, err := conn.Read(reader)
		if err != nil {
			log.Println("Error in reading request")
			conn.Close()
			continue
		}
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nHello\n"))
		conn.Close()
	}

}

func (w *WorkerPool) Submit(job Job) {
	w.jobs <- job
}

func (w *WorkerPool) Close() {
	close(w.jobs)
	w.wg.Wait()
}

func main() {
	workers := NewWorkerPool(1, 10)
	rateLimiter := NewTokenBucket(2, 5)
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Listener failed..")
	}
	log.Println("Server started at :8080")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Client not connected")
			continue
		}

		if !rateLimiter.Allow() {
			conn.Write([]byte("HTTP/1.1 429 Too Many Requests\r\n\r\nRate limit exceeded\n"))
			conn.Close()
			continue
		}

		workers.Submit(Job{Conn: conn})
	}
}
