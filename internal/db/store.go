package db 

import ( 
	"fmt"  
	"net" 
	"sync"
) 

var (
	store 	 = make(map[string]string) // key value store 
	watchers = make(map[string][]net.Conn) 
	mu		 	sync.RWMutex // synchornous mutex allows only one goroutine to lock and unlock at once
)

func HandleSet(conn net.Conn , key string , val string ) {
	mu.Lock() 
	store[key] = val 
	conns := watchers[key]
	mu.Unlock()  

	for _ , watcher := range conns {
		if watcher != conn { // the setting function is sending out exept to itself !
			watcher.Write([]byte(fmt.Sprintf("WATCH : '%s' ---> '%s' \n" , key ,val ))) 
			watcher.Write([]byte("real-db> "))
		}
	}  
	conn.Write([]byte("real-db> "))
	 
} 

func HandleGet( conn net.Conn , key string ) {
	mu.RLock()
	val , ok := store[key] 
	mu.RUnlock()

	if ok {
		conn.Write([]byte(fmt.Sprintf("%s \n", val))) 
	} else  {
		conn.Write([]byte("nil \n"))
	}
	conn.Write([]byte("real-db> "))
}  

func HandleDelete(conn net.Conn , key string ){
	mu.Lock()  
	conns := watchers[key]
	delete(store, key) 
	mu.Unlock()  

	for _ , watcher :=  range conns { 
		if watcher != conn { // the setting function is sending out exept to itself !
			watcher.Write([]byte(fmt.Sprintf("WATCH : '%s' ---> '%s' \n" , key , "nil" ))) 
			watcher.Write([]byte("real-db> "))
		}
	}
	conn.Write([]byte("real-db> "))
}