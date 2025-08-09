# Real-DB 🧱 
### About Real-DB
RealDB is a project I started to learn and explore the internals of distributed systems, in-memory databases, and high-performance storage. Built entirely in Go, RealDB follows a master–slave architecture and uses consistent hashing to evenly distribute load and ensure fault tolerance. It supports LRU eviction policies to gracefully handle memory overload and is concurrency-safe through mutex locking.

My aim with RealDB was to build something lightweight, fast, and scalable that could handle real-time, data-intensive workloads while being simple to extend for new features. This project was both a learning experience and a way to experiment with concepts like distributed data sharding, replication, and efficient eviction mechanisms.
## Performance 
```
|       | GET (Read Latency)  | SET (Write Latency) |
|-------|---------------------|-------------------- |
| p50   | 0.159459 ms         | 0.279375 ms         |
| p95   | 0.298292 ms         | 0.563084 ms         |
| p99   | 0.385709 ms         | 0.759375 ms         |

Throughput: ~44k ops/sec
```
## Get Started  
### Setting up Real-DB with Docker 
Quick way to setup Real-DB is through running this commmand 
```
docker run -p 6369:6369 tashwinsj/realdb:latest 
``` 
### Connecting to the DB 
From another terminal window run this command to connect as a database client 
```
$ nc localhost 6369
``` 
### Usage example 
To set a key value pair, run the command 
```
real-db> $ SET key-string val-string 
``` 
To get a key  
```
real-db> $ GET key-string 
``` 
To get real-time data of a particular key
real-db> $ WATCH key-string 

## Functionalities

### Major milestones 
- Distributed caching - ✔️ 
   - Consistent hashing - ✔️ 
   - Local LRU Eviction - ✔️ 
   - Backup of Evictied keys to Persistance secondary storage
   - If its cache miss -> persistant backup -> bring on to cache(in-memory) again.
- Persistance to secondary memory
- Pub/sub model
### Concurrency and Reliability 
- Proper connection management (Close watchers on disconnect) - ✔️ 
- Concurrent-safe access with sync.Map or channels - ✔️ 
- Client timeout/idle disconnection 
- Logging and structured error handling - ✔️ 

### Core Redis-like Features
- Key Expiration (SET key value EX 10)
- DELETE (DEL key) command - ✔️ 
- Increment/Decrement (INCR, DECR) - ✔️ 
- Multi-key operations (MGET, MSET)
- Data types:
    - Strings - ✔️ 
    - Lists ( LPUSH , LRANGE )
    - Hashes ( HSET , HGET )
    - Sets (SADD , SMEMBERS ) 
- Pattern matching keys (KEYS pattern*) 

### Persistence 
- Basic logging - ✔️ 
- AOF (Append Only File) logging - ✔️ 
- Snapshot-based persistence (RDB-like)
- Loading from disk on startup  

### Performance Features
- LRU cache for memory optimization - ✔️ 
- Sharding support (manual for now)
- Compression for large values 

### Pub/Sub
- Publish/Subscribe system (PUBLISH, SUBSCRIBE)
- Keyspace notifications for watchers


