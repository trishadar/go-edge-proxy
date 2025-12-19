package middleware

import (
	"go-edge-proxy/ratelimit"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	clients = make(map[string]*ratelimit.Limiter)
	mu      sync.Mutex
)

func getLimiter(ip string) *ratelimit.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if limiter, ok := clients[ip]; ok {
		return limiter
	}

	// Example: VIP IPs have higher limits
	var limiter *ratelimit.Limiter
	if ip == "127.0.0.2" { // pretend VIP IP
		limiter = ratelimit.NewLimiter(5, 10) // 5/sec, burst 10
	} else {
		limiter = ratelimit.NewLimiter(1, 5) // normal user
	}

	clients[ip] = limiter
	return limiter
}

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Skip rate limiting for /metrics
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		limiter := getLimiter(ip)

		start := time.Now()
		if !limiter.Allow() {
			RecordRequest(time.Since(start), true)
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
		RecordRequest(time.Since(start), false)
	})
}
