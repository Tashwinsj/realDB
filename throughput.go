package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {
	serverAddr := "localhost:6369"
	numOps := 10000 // total operations
	numClients := 10 // number of concurrent clients

	var wg sync.WaitGroup
	start := time.Now()

	for c := 0; c < numClients; c++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()

			conn, err := net.Dial("tcp", serverAddr)
			if err != nil {
				panic(err)
			}
			defer conn.Close()

			reader := bufio.NewReader(conn)
			// Read the welcome message
			reader.ReadString('>')

			for i := 0; i < numOps/numClients; i++ {
				// SET command
				fmt.Fprintf(conn, "SET key%d value%d\n", i, i)
				reader.ReadString('>') // wait for prompt

				// GET command
				fmt.Fprintf(conn, "GET key%d\n", i)
				reader.ReadString('>') // wait for prompt
			}
		}(c)
	}

	wg.Wait()
	elapsed := time.Since(start)
	totalOps := numOps * 2 // because we do SET + GET each loop
	throughput := float64(totalOps) / elapsed.Seconds()

	fmt.Printf("Completed %d operations in %v\n", totalOps, elapsed)
	fmt.Printf("Throughput: %.2f ops/sec\n", throughput)
}
