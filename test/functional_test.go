package test

import (
	"bufio"
	"net"
	"strings"
	"testing"
)

func sendCommand(t *testing.T, cmd string) string {
	conn, err := net.Dial("tcp", "localhost:6366")
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Flush the welcome banner and prompt
	reader := bufio.NewReader(conn)
	// Just skip initial banner (read until we get prompt or timeout after 5 lines)
	for i := 0; i < 5; i++ {
	line, err := reader.ReadString('\n')
	if err != nil || strings.Contains(line, "real-db>") {
		break
	}
}


	_, err = conn.Write([]byte(cmd + "\n"))
	if err != nil {
		t.Fatalf("Failed to write: %v", err)
	}

	resp, _ := reader.ReadString('\n')
	return strings.TrimSpace(resp)
}

func TestSetGet(t *testing.T) {
	resp := sendCommand(t, "set foo bar")
	if resp != "" {
		t.Errorf("Expected empty response for SET, got: %s", resp)
	}

	resp = sendCommand(t, "get foo")
	if resp != "bar" {
		t.Errorf("Expected 'bar', got: %s", resp)
	}
}

func TestIncDec(t *testing.T) {
	sendCommand(t, "set num 10")

	resp := sendCommand(t, "inc num")
	if resp != "11" {
		t.Errorf("Expected 11, got: %s", resp)
	}

	resp = sendCommand(t, "dec num")
	if resp != "10" {
		t.Errorf("Expected 10, got: %s", resp)
	}
}

func TestInvalidInc(t *testing.T) {
	sendCommand(t, "set str val")

	resp := sendCommand(t, "inc str")
	if resp != "ERR: value is not an integer" {
		t.Errorf("Expected error for non-int INC, got: %s", resp)
	}
}
