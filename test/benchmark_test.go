package test

import (
	"fmt"
	"net"
	"testing"
	"bufio"
	"strings"
)

func BenchmarkSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conn, _ := net.Dial("tcp", "localhost:6366")
		reader := bufio.NewReader(conn)
		for {
			line, _ := reader.ReadString('\n')
			if strings.HasSuffix(line, "real-db> ") {
				break
			}
		}
		fmt.Fprintf(conn, "set key%d val%d\n", i, i)
		reader.ReadString('\n')
		conn.Close()
	}
}
