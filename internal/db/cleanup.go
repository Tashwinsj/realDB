package db

import "net"

func RemoveConnFromWatchers(conn net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	for key , conns := range watchers {
		newConns := make([]net.Conn , 0 , len(conns)) 
		for _ , c := range conns {
			if c != conn {
				newConns = append(newConns, c) 
			} 
		} 
		if len(newConns) == 0 {
			delete(watchers, key) 
		} else { 
			watchers[key] = newConns
		}
	}
}