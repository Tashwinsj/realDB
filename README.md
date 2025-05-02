# Real-DB 🧱 
### About Real-DB
Real-DB is an in-memory real-time database with a focus on high performance, concurrency, and reliability. It aims to provide a feature set similar to Redis, while also incorporating enhancements for improved data management and scalability. Real-DB is suitable for applications requiring fast data access, efficient memory utilization, and reliable data persistence. 
 
## Get Started  
### Setting up Real-DB with Docker 
Quick way to setup Real-DB is through running this commmand 
```
docker run -p 6369:6369 tashwinsj/real-db:latest 
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

### Concurrency and Reliability 
- Proper connection management (Close watchers on disconnect) - ✔️ 
- Concurrent-safe access with sync.Map or channels 
- Client timeout/idle disconnection 
- Logging and structured error handling 

### Core Redis-like Features
- Key Expiration (SET key value EX 10)
- DELETE (DEL key) command - ✔️ 
- Increment/Decrement (INCR, DECR)
- Multi-key operations (MGET, MSET)
- Data types:
    - Strings ( already done )  
    - Lists ( LPUSH , LRANGE )
    - Hashes ( HSET , HGET )
    - Sets (SADD , SMEMBERS ) 
- Pattern matching keys (KEYS pattern*) 

### Persistence 
- Basic logging - ✔️ 
- AOF (Append Only File) logging
- Snapshot-based persistence (RDB-like)
- Loading from disk on startup  

### Performance Features
- LRU cache for memory optimization - ✔️ 
- Sharding support (manual for now)
- Compression for large values 

### Pub/Sub
- Publish/Subscribe system (PUBLISH, SUBSCRIBE)
- Keyspace notifications for watchers


