package main

import (
	"net/url"
	"sync"
)

type BackendPool struct {
	backends []*url.URL
	index    int
	mu       sync.Mutex
}

// NewBackendPool creates a pool of backend URLs
func NewBackendPool(urls []string) *BackendPool {
	backends := []*url.URL{}
	for _, u := range urls {
		parsed, _ := url.Parse(u)
		backends = append(backends, parsed)
	}
	return &BackendPool{
		backends: backends,
		index:    0,
	}
}

// NextBackend returns the next backend in round-robin order
func (p *BackendPool) NextBackend() *url.URL {
	p.mu.Lock()
	defer p.mu.Unlock()
	backend := p.backends[p.index]
	p.index = (p.index + 1) % len(p.backends)
	return backend
}
