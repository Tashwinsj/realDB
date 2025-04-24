package db

import ("net" 
		"github.com/sirupsen/logrus" 
)

func HandleWatch(conn net.Conn , key string) {
	mu.Lock()
	watchers[key] = append(watchers[key] , conn) 
	mu.Unlock()	

	conn.Write([]byte("WATCHING " + key + "\n" )) 
	remoteAddress := conn.RemoteAddr() 
		if tcpAddr, ok := remoteAddress.(*net.TCPAddr); ok {
			ipAddress := tcpAddr.IP
			portNumber := tcpAddr.Port
			logrus.Infof("[IP: %s  PORT: %d] WATCHING %s", ipAddress, portNumber , key)
		} else {
			logrus.Warn("New client connected, could not cast the remote address of the connected client")
		}
	conn.Write([]byte("real-db> "))
} 
