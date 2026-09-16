package serial

import (
	"context"
	"errors"
	"testing"
	"time"

	domainErrors "sv-printer/internal/domain/errors"
)

type fakePort struct {
	written []byte
	closed  bool
	err     error
}

func (f *fakePort) Write(p []byte) (int, error) {
	if f.err != nil {
		return 0, f.err
	}
	f.written = append(f.written, p...)
	return len(p), nil
}

func (f *fakePort) Close() error {
	f.closed = true
	return nil
}

type blockingPort struct {
	started chan struct{}
	release chan struct{}
}

func (b *blockingPort) Write(p []byte) (int, error) {
	close(b.started)
	<-b.release
	return len(p), nil
}

func (b *blockingPort) Close() error {
	close(b.release)
	return nil
}

func TestNewSerialTransportDefaultBaud(t *testing.T) {
	tr := NewSerialTransport("/dev/ttyUSB0", 0).(*SerialTransport)
	if tr.baud != 9600 {
		t.Errorf("baud = %d, want 9600", tr.baud)
	}
}

func TestWriteSuccess(t *testing.T) {
	fp := &fakePort{}
	tr := &SerialTransport{port: fp}

	if err := tr.Write(context.Background(), []byte("hello")); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if string(fp.written) != "hello" {
		t.Errorf("written = %q, want hello", fp.written)
	}
}

func TestWriteWrapsError(t *testing.T) {
	fp := &fakePort{err: errors.New("boom")}
	tr := &SerialTransport{port: fp}

	err := tr.Write(context.Background(), []byte("x"))
	if err == nil {
		t.Fatal("Write() expected error")
	}
	if !errors.Is(err, domainErrors.ErrPrintFailed) {
		t.Errorf("Write() error = %v, want wrapped ErrPrintFailed", err)
	}
}

func TestWriteNotConnected(t *testing.T) {
	tr := &SerialTransport{}

	err := tr.Write(context.Background(), []byte("x"))
	if !errors.Is(err, domainErrors.ErrPrinterConnectionFailed) {
		t.Errorf("Write() error = %v, want ErrPrinterConnectionFailed", err)
	}
}

func TestWriteContextCancellation(t *testing.T) {
	bp := &blockingPort{started: make(chan struct{}), release: make(chan struct{})}
	tr := &SerialTransport{port: bp}

	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() {
		errCh <- tr.Write(ctx, []byte("x"))
	}()

	<-bp.started
	cancel()

	select {
	case err := <-errCh:
		if err == nil || !errors.Is(err, context.Canceled) {
			t.Errorf("Write() error = %v, want context.Canceled", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Write() did not return on cancellation")
	}

	bp.Close()
}

func TestClose(t *testing.T) {
	fp := &fakePort{}
	tr := &SerialTransport{port: fp}

	if err := tr.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if !fp.closed {
		t.Error("Close() should close the port")
	}
	if tr.port != nil {
		t.Error("Close() should clear the port")
	}
}

func TestCloseNilPort(t *testing.T) {
	tr := &SerialTransport{}
	if err := tr.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}
