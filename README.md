# High-Throughput In-Memory Cache Service

A lightweight, concurrent Go in-memory cache service engineered for high-concurrency ingestion pipelines.

## Architecture

- `pkg/cache/cache.go`: Sharded lock partitions using `sync.RWMutex` to prevent CPU lock contention at scale.
- `cmd/server/main.go`: HTTP API exposing GET/POST key-value cache operations.
