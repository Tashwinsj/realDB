package db

import "net" 

func HandleWatch(conn net.Conn , key string) {
	mu.Lock()
	watchers[key] = append(watchers[key] , conn) 
	mu.Unlock()	

	conn.Write([]byte("WATCHING " + key + "\n" )) 
	conn.Write([]byte("real-db> "))
} 
