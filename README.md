# Threadis

Threadis is a high-performance, in-memory key-value database built from scratch in Go.

Inspired by Redis, it implements the **RESP (Redis Serialization Protocol)**, making it compatible with `redis-cli` and other Redis clients.

This project was created for educational purposes to explore the design and implementation of modern database systems, including event-driven network servers, custom data structures, and probabilistic algorithms.

## Key Features
#### Redis Protocol (RESP) Compliant
- Fully compatible with the RESP.
- Works seamlessly with `redis-cli` and Redis-compatible clients.

#### High-Performance I/O
- Event-driven, single-threaded architecture.
- Uses I/O multiplexing: epoll (Linux)
- Capable of handling thousands of concurrent client connections with minimal overhead.

#### Custom Data Structures
Built from scratch without relying on external database implementations:
- **Skiplist**:
    - Powers sorted set operations such as `ZADD`, `ZRANK`, `ZSCORE`, and `ZREM`.
    - Provides efficient logarithmic-time insertion, deletion, and lookup.

#### Probabilistic Data Structures
- **Scalable Bloom Filter**: For fast, memory-efficient set membership testing

- **Count-Min Sketch**: For estimating item frequencies in a data stream

## Getting Started
```
  # run server 
  go run cmd/main.go
```
## Supported Commands

The server supports a wide range of commands grouped by data type:

| Category | Commands |
| :--- | :--- |
| **General** | `PING` |
| **String** | `SET`, `GET`, `DEL`, `EXISTS`, `TTL`, `EXPIRE` |
| **Sorted Set**| `ZADD`, `ZRANK`, `ZREM`, `ZSCORE`, `ZRANGE` |
| **Set** | `SADD`, `SREM`, `SMEMBERS`, `SISMEMBER` |
| **Bloom Filter**| `BF.RESERVE`, `BF.MADD`, `BF.EXISTS` |
| **Count-Min** | `CMS.INITBYDIM`, `CMS.INITBYPROB`, `CMS.INCRBY`, `CMS.QUERY` |

## Test with Redis CLI
```bash
redis-cli -p 3000

SET session abc123
EXPIRE session 10
TTL session

GET session
# wait 10s

GET session
TTL session
```
