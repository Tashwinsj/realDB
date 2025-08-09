package test

import (
	"fmt"
	"net"
	"strings"
	"bufio"
	"sync"
	"testing"
)

func TestConcurrentLoad(t *testing.T) {
	var wg sync.WaitGroup
	clientCount := 100

	for i := 0; i < clientCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			conn, err := net.Dial("tcp", "localhost:6366")
			if err != nil {
				t.Errorf("Client %d failed to connect: %v", id, err)
				return
			}
			defer conn.Close()

			reader := bufio.NewReader(conn)
			for {
				line, _ := reader.ReadString('\n')
				if strings.HasSuffix(line, "real-db> ") {
					break
				}
			}

			for j := 0; j < 10; j++ {
				cmd := fmt.Sprintf("set key%d-%d val%d", id, j, j)
				fmt.Fprintln(conn, cmd)
				reader.ReadString('\n')
			}
		}(i)
	}
	wg.Wait()
}
