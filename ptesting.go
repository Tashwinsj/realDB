package main

import (
	"bufio"
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

func percentile(latencies []time.Duration, p float64) time.Duration {
	if len(latencies) == 0 {
		return 0
	}
	index := int(float64(len(latencies)-1) * p)
	return latencies[index]
}

func main() {
	serverAddr := "localhost:6369"
	numOps := 10000 // total SET+GET pairs
	numClients := 10

	var wg sync.WaitGroup
	var mu sync.Mutex
	var setLatencies, getLatencies []time.Duration

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
			// Read welcome banner
			reader.ReadString('>')

			for i := 0; i < numOps/numClients; i++ {
				// Measure SET latency
				setStart := time.Now()
				fmt.Fprintf(conn, "SET key%d value%d\n", i, i)
				reader.ReadString('>') // wait for prompt
				setLatency := time.Since(setStart)

				// Measure GET latency
				getStart := time.Now()
				fmt.Fprintf(conn, "GET key%d\n", i)
				reader.ReadString('>') // wait for prompt
				getLatency := time.Since(getStart)

				// Store results thread-safe
				mu.Lock()
				setLatencies = append(setLatencies, setLatency)
				getLatencies = append(getLatencies, getLatency)
				mu.Unlock()
			}
		}(c)
	}

	wg.Wait()
	elapsed := time.Since(start)
	totalOps := numOps * 2
	throughput := float64(totalOps) / elapsed.Seconds()

	// Sort to compute percentiles
	sort.Slice(setLatencies, func(i, j int) bool { return setLatencies[i] < setLatencies[j] })
	sort.Slice(getLatencies, func(i, j int) bool { return getLatencies[i] < getLatencies[j] })

	fmt.Printf("Completed %d operations in %v\n", totalOps, elapsed)
	fmt.Printf("Throughput: %.2f ops/sec\n\n", throughput)

	fmt.Println("Write Latency (SET):")
	fmt.Printf("  p50: %v\n", percentile(setLatencies, 0.50))
	fmt.Printf("  p95: %v\n", percentile(setLatencies, 0.95))
	fmt.Printf("  p99: %v\n", percentile(setLatencies, 0.99))

	fmt.Println("\nRead Latency (GET):")
	fmt.Printf("  p50: %v\n", percentile(getLatencies, 0.50))
	fmt.Printf("  p95: %v\n", percentile(getLatencies, 0.95))
	fmt.Printf("  p99: %v\n", percentile(getLatencies, 0.99))
}
