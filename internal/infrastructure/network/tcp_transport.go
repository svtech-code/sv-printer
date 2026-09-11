package network

import (
	"context"
	"fmt"
	"net"
	"time"

	"sv-print/internal/domain/errors"
	"sv-print/internal/infrastructure/transport"
)

type TCPTransport struct {
	address string
	timeout time.Duration
	conn    net.Conn
}

func NewTCPTransport(address string, timeout time.Duration) transport.PrinterTransport {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return &TCPTransport{
		address: address,
		timeout: timeout,
	}
}

func (t *TCPTransport) Open(ctx context.Context) error {
	var d net.Dialer
	d.Timeout = t.timeout

	conn, err := d.DialContext(ctx, "tcp", t.address)
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrPrinterConnectionFailed, err)
	}
	t.conn = conn
	return nil
}

func (t *TCPTransport) Write(ctx context.Context, data []byte) error {
	if t.conn == nil {
		return fmt.Errorf("%w: not connected", errors.ErrPrinterConnectionFailed)
	}

	// Set write deadline based on context
	if deadline, ok := ctx.Deadline(); ok {
		t.conn.SetWriteDeadline(deadline)
	} else {
		t.conn.SetWriteDeadline(time.Now().Add(t.timeout))
	}

	_, err := t.conn.Write(data)
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrPrintFailed, err)
	}
	return nil
}

func (t *TCPTransport) Close() error {
	if t.conn != nil {
		err := t.conn.Close()
		t.conn = nil
		return err
	}
	return nil
}
