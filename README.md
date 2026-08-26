# concurrent-kv-store# In-Memory Thread-Safe Key-Value Store with TTL

A high-performance, concurrent in-memory key-value store built with Go standard primitives. It features dual eviction (lazy on-read and periodic background sweeps), fine-grained concurrency control via reader-writer mutual exclusion locks, and typed domain errors.

---

## Key Features

* **Thread-Safe Concurrency:** Implements `sync.RWMutex` to allow concurrent readers (`RLock`) while ensuring exclusive access for write operations (`Lock`).
* **Dual Eviction Strategy:**
  * **Lazy Eviction:** Keys checked during read operations (`Get`) are evaluated against their TTL and deleted immediately if expired.
  * **Active Sweeping:** A background goroutine managed via `time.Ticker` sweeps and purges expired entries to prevent memory retention leaks.
* **Graceful Lifecycle Management:** Explicit channel-based teardown (`Close()`) stops the background worker and releases runtime timer resources.
* **Race-Condition Free:** Verified under high concurrent load with Go's race detector (`go test -race`).

---

## Architecture Overview

```
                      +-----------------------------+
                      |         Client Calls        |
                      | (Set / Get / Delete / Close)|
                      +--------------+--------------+
                                     |
             +-----------------------+-----------------------+
             |                       |                       |
             v                       v                       v
      +--------------+        +--------------+        +--------------+
      |  Set(k,v,d)  |        |    Get(k)    |        |  Delete(k)   |
      |   [W-Lock]   |        |   [R-Lock]   |        |   [W-Lock]   |
      +--------------+        +--------------+        +--------------+
             |                       |                       |
             +-----------------------+-----------------------+
                                     |
                                     v
                      +-----------------------------+
                      |         MemoryStore         |
                      |   map[string]item + RWMutex |
                      +--------------^--------------+
                                     | [W-Lock]
                      +--------------+--------------+
                      |  Background Cleanup Worker  |
                      |  (time.Ticker interval)     |
                      +-----------------------------+
```

---

## Directory Structure

```text
kvstore/
├── store/
│   ├── errors.go       # Sentinel domain errors
│   ├── item.go         # Value wrapper and TTL evaluation logic
│   ├── store.go        # Thread-safe KeyValueStore implementation
│   └── store_test.go   # Table-driven, expiration, and race tests
├── go.mod
├── main.go             # CLI usage example
└── README.md
```

---

## Getting Started

### Prerequisites
* Go 1.21 or higher installed.

### Installation & Run

1. Clone the repository:
   ```bash
   git clone [https://github.com/](https://github.com/)<your-username>/<repo-name>.git
   cd <repo-name>
   ```

2. Run the application:
   ```bash
   go run main.go
   ```

---

## Running Tests

Execute the comprehensive test suite with verbose output and the Go race detector enabled:

```bash
go test -v -race ./...
```
