package serial

import (
	"context"
	"fmt"
	"time"

	goserial "go.bug.st/serial"

	domainErrors "sv-printer/internal/domain/errors"
	"sv-printer/internal/infrastructure/transport"
)

type writer interface {
	Write(p []byte) (n int, err error)
	Close() error
}

type SerialTransport struct {
	portName string
	baud     int
	timeout  time.Duration
	port     writer
}

func NewSerialTransport(portName string, baud int) transport.PrinterTransport {
	if baud == 0 {
		baud = 9600
	}
	return &SerialTransport{
		portName: portName,
		baud:     baud,
		timeout:  5 * time.Second,
	}
}

func (t *SerialTransport) Open(ctx context.Context) error {
	mode := &goserial.Mode{
		BaudRate: t.baud,
		DataBits: 8,
		Parity:   goserial.NoParity,
		StopBits: goserial.OneStopBit,
	}

	port, err := goserial.Open(t.portName, mode)
	if err != nil {
		return fmt.Errorf("%w: %v", domainErrors.ErrPrinterConnectionFailed, err)
	}
	t.port = port
	return nil
}

func (t *SerialTransport) Write(ctx context.Context, data []byte) error {
	if t.port == nil {
		return fmt.Errorf("%w: not connected", domainErrors.ErrPrinterConnectionFailed)
	}

	done := make(chan error, 1)
	go func() {
		_, err := t.port.Write(data)
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("%w: %v", domainErrors.ErrPrintFailed, err)
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *SerialTransport) Close() error {
	if t.port != nil {
		err := t.port.Close()
		t.port = nil
		return err
	}
	return nil
}
