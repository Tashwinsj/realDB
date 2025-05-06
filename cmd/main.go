package main

import (
	"fmt"
	"net"
	"os" 
 	"os/signal"
	"realDB/internal/server"
	"syscall"
	"time"
	"github.com/sirupsen/logrus"
) 

var (
	shutdownChan = make(chan struct{})
)

func init() {
	file, err := os.OpenFile("realdb.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logrus.SetOutput(file)
	} else {
		logrus.Info("Failed to log to file, using default stderr")
	}
}

func main(){
	ln , err := net.Listen("tcp" , ":6369") 
	if err != nil {
		panic(err)
	} 
	fmt.Println("Server started on port 6369")   
	logrus.Info("Server Started on port 6369")

	go HandleShutDown(ln)

	for {
		conn , err := ln.Accept()
		if err != nil {
				select {
					case <-shutdownChan:
						// Expected error during shutdown — suppress or log cleanly
						logrus.Info("Stopped accepting new connections (listener closed).")
						return
					default:
						// Unexpected error — log it
						logrus.Errorf("Connection error : %v", err)
						continue
				}
		} 
		conn.Write([]byte( 
		`
.-----------------------------------------------------------.
|                                                           |
|     ██████╗ ███████╗ █████╗ ██╗      ██████╗ ██████╗      |
|     ██╔══██╗██╔════╝██╔══██╗██║      ██╔══██╗██╔══██╗     |
|     ██████╔╝█████╗  ███████║██║█████╗██║  ██║██████╔╝     |
|     ██╔══██╗██╔══╝  ██╔══██║██║╚════╝██║  ██║██╔══██╗     |
|     ██║  ██║███████╗██║  ██║███████╗ ██████╔╝██████╔╝     |
|     ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚═════╝      |
|                                                           |
|                  Connected to real-db                     |
'-----------------------------------------------------------'` + "\n")) 
		conn.Write([]byte("real-db> "))  

		remoteAddress := conn.RemoteAddr() // Logic to log the connected client detials [IP:PORT]
		if tcpAddr, ok := remoteAddress.(*net.TCPAddr); ok {
			ipAddress := tcpAddr.IP
			portNumber := tcpAddr.Port
			logrus.Infof("New client connected IP: %s  PORT: %d", ipAddress, portNumber)
		} else {
			logrus.Warn("New client connected, could not cast the remote address of the connected client")
		}

 
		go server.HandleConnection(conn)      // This infinite loop keeps looking for connections and wherever they come a seperate go routine is 
											  // spawned to handle the request ( seperate go routine for every client connected )
	}
  
}
 
func HandleShutDown(listener net.Listener) {
	// Capture system interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	logrus.Infof("Received signal: %s. Shutting down gracefully!", sig)

	// Stop new connections
	close(shutdownChan)
	listener.Close()

	// Optional: wait for ongoing goroutines to finish
	time.Sleep(1 * time.Second) // adjust as needed

	// Clean up resources (watchers, etc.)
	logrus.Info("Graceful shutdown complete.")
	time.Sleep(1000 * time.Millisecond)
	logrus.Exit(0)
}
