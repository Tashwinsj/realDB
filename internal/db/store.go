package db

import (
	"fmt"
	"net"
	"sync"
	"github.com/sirupsen/logrus" 
	"realDB/internal/cache"
) 


var (
	lru = cache.NewLRUCache(25)
	watchers = make(map[string][]net.Conn) 
	mu		 	sync.RWMutex // synchornous mutex allows only one goroutine to lock and unlock at once
)

func HandleSet(conn net.Conn , key string , val string ) {
	mu.Lock() 
	lru.Set(key , val)
	conns := watchers[key]
	mu.Unlock()  

	for _ , watcher := range conns {
		if watcher != conn { // the setting function is sending out exept to itself !
			watcher.Write([]byte(fmt.Sprintf("WATCH : '%s' ---> '%s' \n" , key ,val ))) 
			watcher.Write([]byte("real-db> "))
		}
	}   
	logrus.Infof("SET %s %s", key , val )
	conn.Write([]byte("real-db> "))
	 
} 

func HandleGet( conn net.Conn , key string ) {
	mu.RLock()
	val , ok :=  lru.Get(key)
	mu.RUnlock()

	if ok {
		conn.Write([]byte(fmt.Sprintf("%s \n", val))) 
	} else  {
		conn.Write([]byte("nil \n"))
	} 
	logrus.Infof("GET %s", key)
	conn.Write([]byte("real-db> "))
}  

func HandleDelete(conn net.Conn , key string ){
	mu.Lock()  
	conns := watchers[key]
	lru.Del(key)
	mu.Unlock()  

	for _ , watcher :=  range conns { 
		if watcher != conn { // the setting function is sending out exept to itself !
			watcher.Write([]byte(fmt.Sprintf("WATCH : '%s' ---> '%s' \n" , key , "nil" ))) 
			watcher.Write([]byte("real-db> "))
		}
	} 
	logrus.Infof("DEL %s", key)
	conn.Write([]byte("real-db> "))
}