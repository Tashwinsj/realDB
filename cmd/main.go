package main

import (
	"fmt"
	"net"
	"os" 
	"log"
	"os/signal"
	"realDB/internal/server"
	"syscall"
	"time"
) 

var (
	shutdownChan = make(chan struct{})
)


func main(){
	ln , err := net.Listen("tcp" , ":6369") 
	if err != nil {
		panic(err)
	} 
	fmt.Println("Server started on port 6369")  

	go HandleShutDown(ln)

	for {
		conn , err := ln.Accept()
		if err != nil {
				select {
					case <-shutdownChan:
						// Expected error during shutdown — suppress or log cleanly
						log.Println("Stopped accepting new connections (listener closed).")
						return
					default:
						// Unexpected error — log it
						log.Printf("Connection error : %v", err)
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
|                  Connected to real-db>                    |
'-----------------------------------------------------------'` + "\n")) 
		conn.Write([]byte("real-db> "))

		go server.HandleConnection(conn)      // This infinite loop keeps looking for connections and wherever they come a seperate go routine is 
										// created to handle the request ( seperate go routine for every client connected )
	}
  
}
 
func HandleShutDown(listener net.Listener) {
	// Capture system interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("Received signal: %s. Shutting down gracefully... \n", sig)

	// Stop new connections
	close(shutdownChan)
	listener.Close()

	// Optional: wait for ongoing goroutines to finish
	time.Sleep(1 * time.Second) // adjust as needed

	// Clean up resources (watchers, etc.)
	log.Printf("Graceful shutdown complete.")
	os.Exit(0)
}
