package internals

import (
	"log"
	"net"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type Job struct {
	Conn net.Conn
}

type WorkerPool struct {
	jobs           chan Job
	maxWorkers     int
	wg             sync.WaitGroup
	processedCount prometheus.Counter
}

func NewWorkerPool(maxWorkers, queueSize int, processedCount prometheus.Counter) *WorkerPool {
	jobs := make(chan Job, queueSize+maxWorkers)

	newWorker := &WorkerPool{
		maxWorkers:     maxWorkers,
		jobs:           jobs,
		processedCount: processedCount,
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
		w.processedCount.Inc()
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
