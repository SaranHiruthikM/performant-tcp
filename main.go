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

	return newWorker
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	log.Println("Client Connected")
	reader := make([]byte, 1024)
	_, err := conn.Read(reader)
	if err != nil {
		log.Println("Error in reading request")
		return
	}
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nHello\n"))
}

func main() {
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

		go handleConn(conn)
	}
}
