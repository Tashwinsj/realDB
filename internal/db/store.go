package db

import (
	"fmt"
	"net"
	"realDB/internal/cache"
	"strconv"
	"sync"
	"github.com/sirupsen/logrus"
) 


var (
	lru = cache.NewLRUCache(928202913952)
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

func HandleINC( conn net.Conn , key string) {
	mu.RLock()  
	val , ok := lru.Get(key)  
	mu.RUnlock()

	// if the value is not present we treat it as zero
	if !ok {
		val = "0"
	} 

	intVal , err := strconv.Atoi(val) 
	if err != nil {
		conn.Write([]byte("ERR: value is not an integer \n real-db> ")) 
		return 
	}
	
	intVal++ 
 
	newVal :=strconv.Itoa(intVal)  
	mu.Lock()
	lru.Set(key , newVal) 
	mu.Unlock()

	conns := watchers[key] 
	for _, watcher := range conns {
		if watcher != conn {
			watcher.Write([]byte(fmt.Sprintf("WATCH : '%s' ---> '%s'\n", key, newVal)))
			watcher.Write([]byte("real-db> "))
		}
	}  

	logrus.Infof("INCR %s = %s", key, newVal)
	conn.Write([]byte(fmt.Sprintf("%s\nreal-db> ", newVal)))

	
}


func HandleDEC( conn net.Conn , key string) {
	mu.RLock()  
	val , ok := lru.Get(key)  
	mu.RUnlock()

	// if the value is not present return an error
	if !ok {
		conn.Write([]byte("ERR: Value not found in the DB \n real-db> "))  
		return 
	} 

	intVal , err := strconv.Atoi(val) 
	if err != nil {
		conn.Write([]byte("ERR: value is not an integer \n real-db> ")) 
		return 
	}
	
	intVal--
 
	newVal :=strconv.Itoa(intVal)  
	mu.Lock()
	lru.Set(key , newVal) 
	mu.Unlock()

	conns := watchers[key] 
	for _, watcher := range conns {
		if watcher != conn {
			watcher.Write([]byte(fmt.Sprintf("WATCH : '%s' ---> '%s'\n", key, newVal)))
			watcher.Write([]byte("real-db> "))
		}
	}  

	logrus.Infof("INCR %s = %s", key, newVal)
	conn.Write([]byte(fmt.Sprintf("%s\nreal-db> ", newVal)))
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