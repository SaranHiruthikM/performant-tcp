package main

import (
	"log"
	"net"
	"sync"
)

type Job struct {
	Conn net.Conn
}

type WorkerPool struct {
	jobs       chan Job
	maxWorkers int
	wg         sync.WaitGroup
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

		workers.Submit(Job{Conn: conn})
	}
}
