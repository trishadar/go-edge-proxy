# Go Edge Proxy

A high-performance reverse proxy built in Go that supports multiple backend servers with round-robin load balancing, per-IP rate limiting, request logging, real-time metrics, and graceful shutdown.

---

## Features

- **Reverse Proxy:** Forwards requests to backend servers.  
- **Load Balancing:** Round-robin across multiple backends.  
- **Rate Limiting:** Per-IP token bucket limits; configurable for VIPs.  
- **Logging:** Tracks method, path, status, and latency.  
- **Metrics:** `/metrics` endpoint shows requests, 429 hits, and average latency; exempt from rate limits.  
- **Graceful Shutdown:** Completes in-flight requests on exit.

---

## Quick Start

1. Run test backends (optional, using Python):

```bash
python3 -m http.server 8081
python3 -m http.server 8082
python3 -m http.server 8083
```

2. Run the proxy:

```bash
go run .
```

3. Access proxy at http://localhost:8080/ and metrics at /metrics.
