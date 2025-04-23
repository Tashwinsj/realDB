# Concurrency and Reliability 
- Proper connection management (Close watchers on disconnect) 
- Concurrent-safe access with sync.Map or channels 
- Client timeout/idle disconnection 
- Logging and structured error handling 

# Core Redis-like Features
- Key Expiration (SET key value EX 10)
- DELETE (DEL key) command
- Increment/Decrement (INCR, DECR)
- Multi-key operations (MGET, MSET)
- Data types:
    - Strings ( already done )  
    - Lists ( LPUSH , LRANGE )
    - Hashes ( HSET , HGET )
    - Sets (SADD , SMEMBERS ) 
- Pattern matching keys (KEYS pattern*) 

# Persistence 
- AOF (Append Only File) logging
- Snapshot-based persistence (RDB-like)
- Loading from disk on startup  

# Performance Features
- LRU cache for memory optimization
- Sharding support (manual for now)
- Compression for large values 

# Pub/Sub
- Publish/Subscribe system (PUBLISH, SUBSCRIBE)
- Keyspace notifications for watchers


