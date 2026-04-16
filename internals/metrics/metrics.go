package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RunMetrics(path string, port int) (prometheus.Counter, prometheus.Counter) {
	requestsProcessed := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "total_requests_processed",
		Help: "Number of requests successfully processed",
	})
	requestsRateLimited := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "total_requests_rate_limited",
		Help: "Number of requests successfully rate limited",
	})
	prometheus.MustRegister(requestsProcessed)
	prometheus.MustRegister(requestsRateLimited)
	go func() {
		http.Handle(path, promhttp.Handler())
		http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()
	return requestsProcessed, requestsRateLimited
}
