package network

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestTCPTransport(t *testing.T) {
	// Start a dummy TCP server
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}
	defer ln.Close()

	address := ln.Addr().String()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			// Just read and discard
			buf := make([]byte, 1024)
			conn.Read(buf)
			conn.Close()
		}
	}()

	transport := NewTCPTransport(address, 2*time.Second)
	ctx := context.Background()

	err = transport.Open(ctx)
	if err != nil {
		t.Fatalf("Expected successful connection, got: %v", err)
	}

	err = transport.Write(ctx, []byte("test payload"))
	if err != nil {
		t.Fatalf("Expected successful write, got: %v", err)
	}

	err = transport.Close()
	if err != nil {
		t.Fatalf("Expected successful close, got: %v", err)
	}
}
