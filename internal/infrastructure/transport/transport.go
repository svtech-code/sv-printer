package transport

import "context"

type PrinterTransport interface {
	Open(ctx context.Context) error
	Write(ctx context.Context, data []byte) error
	Close() error
}
