package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"time"

	"go-edge-proxy/middleware"
)

func main() {

	// Create backend pool
	backendURLs := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}
	pool := NewBackendPool(backendURLs)

	// Create the reverse proxy
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			backend := pool.NextBackend()
			req.URL.Scheme = backend.Scheme
			req.URL.Host = backend.Host
		},
	}

	// Create a new ServeMux (router)
	mux := http.NewServeMux()
	// Add metrics endpoint
	mux.Handle("/metrics", middleware.MetricsHandler())
	// Add reverse proxy wrapped with rate limiter and logging
	mux.Handle("/", middleware.RateLimit(middleware.Logging(proxy)))

	// Use mux as the server handler
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Channel to listen for interrupt (CTRL+C)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Println("Proxy running on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	<-stop // wait for CTRL+C
	log.Println("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	log.Println("Server stopped")
}
