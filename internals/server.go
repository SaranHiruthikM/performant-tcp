package internals

import (
	"github.com/SaranHiruthikM/performant-tcp/internals/ratelimiter"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	listener     net.Listener
	pool         *WorkerPool
	rateLimiter  *ratelimiter.TokenBucket
	shuttingDown atomic.Bool // true when shutting down
}

func NewServer(listener net.Listener, pool *WorkerPool, rateLimiter *ratelimiter.TokenBucket) *Server {
	return &Server{
		listener:    listener,
		pool:        pool,
		rateLimiter: rateLimiter,
	}
}

func (s *Server) Start() {
	log.Println("server started")
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			// check if error is because we shut down intentionally
			if s.shuttingDown.Load() {
				log.Println("listener closed, stopping accept loop")
				return // clean exit
			}
			// otherwise it's a real error, log and continue
			log.Println("accept error:", err)
			continue
		}

		if !s.rateLimiter.Allow() {
			conn.Write([]byte("HTTP/1.1 429 Too Many Requests\r\n\r\nRate limit exceeded\n"))
			conn.Close()
			continue
		}

		s.pool.Submit(Job{Conn: conn})
	}
}

func (s *Server) Shutdown() {
	log.Println("shutting down...")
	s.shuttingDown.Store(true) // set flag first
	s.listener.Close()         // causes Accept() to return error
	s.pool.Close()             // wait for workers to finish
	log.Println("shutdown complete")
}
