package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

var (
	totalRequests int
	rateLimitHits int
	totalLatency  time.Duration
	metricsMu     sync.Mutex
)

func RecordRequest(latency time.Duration, limited bool) {
	metricsMu.Lock()
	defer metricsMu.Unlock()
	totalRequests++
	if limited {
		rateLimitHits++
	}
	totalLatency += latency
}

func MetricsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricsMu.Lock()
		defer metricsMu.Unlock()
		avgLatency := 0.0
		if totalRequests > 0 {
			avgLatency = totalLatency.Seconds() * 1000 / float64(totalRequests) // ms
		}
		data := map[string]interface{}{
			"total_requests":  totalRequests,
			"rate_limit_hits": rateLimitHits,
			"avg_latency_ms":  avgLatency,
		}
		json.NewEncoder(w).Encode(data)
	})
}
