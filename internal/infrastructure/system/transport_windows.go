//go:build windows

package system

import (
	"context"
	"fmt"
	"syscall"
	"unsafe"

	domainErrors "sv-printer/internal/domain/errors"
	"sv-printer/internal/infrastructure/transport"
)

var (
	winspool             = syscall.NewLazyDLL("winspool.drv")
	procOpenPrinter      = winspool.NewProc("OpenPrinterW")
	procStartDocPrinter  = winspool.NewProc("StartDocPrinterW")
	procStartPagePrinter = winspool.NewProc("StartPagePrinter")
	procWritePrinter     = winspool.NewProc("WritePrinter")
	procEndPagePrinter   = winspool.NewProc("EndPagePrinter")
	procEndDocPrinter    = winspool.NewProc("EndDocPrinter")
	procClosePrinter     = winspool.NewProc("ClosePrinter")
)

type SystemTransport struct {
	printerName string
	handle      syscall.Handle
}

func NewSystemTransport(printerName string) transport.PrinterTransport {
	return &SystemTransport{
		printerName: printerName,
	}
}

func (t *SystemTransport) Open(ctx context.Context) error {
	name, err := syscall.UTF16PtrFromString(t.printerName)
	if err != nil {
		return fmt.Errorf("%w: invalid printer name: %v", domainErrors.ErrPrinterConnectionFailed, err)
	}

	var handle syscall.Handle
	ret, _, err := procOpenPrinter.Call(uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(&handle)), 0)
	if ret == 0 {
		// Error returned by Call is a syscall.Errno, so we can log it.
		return fmt.Errorf("%w: OpenPrinter failed: %v", domainErrors.ErrPrinterConnectionFailed, err)
	}
	t.handle = handle
	return nil
}

type docInfo1 struct {
	pDocName    *uint16
	pOutputFile *uint16
	pDatatype   *uint16
}

func (t *SystemTransport) Write(ctx context.Context, data []byte) error {
	if t.handle == 0 {
		return domainErrors.ErrPrinterConnectionFailed
	}

	if len(data) == 0 {
		return nil
	}

	docName, _ := syscall.UTF16PtrFromString("SV Printer Document")
	dataType, _ := syscall.UTF16PtrFromString("RAW")

	docInfo := docInfo1{
		pDocName:  docName,
		pDatatype: dataType,
	}

	ret, _, err := procStartDocPrinter.Call(
		uintptr(t.handle),
		1,
		uintptr(unsafe.Pointer(&docInfo)),
	)
	if ret == 0 {
		return fmt.Errorf("%w: StartDocPrinter failed: %v", domainErrors.ErrPrintFailed, err)
	}

	defer procEndDocPrinter.Call(uintptr(t.handle))

	ret, _, err = procStartPagePrinter.Call(uintptr(t.handle))
	if ret == 0 {
		return fmt.Errorf("%w: StartPagePrinter failed: %v", domainErrors.ErrPrintFailed, err)
	}

	defer procEndPagePrinter.Call(uintptr(t.handle))

	var bytesWritten uint32
	ret, _, err = procWritePrinter.Call(
		uintptr(t.handle),
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
		uintptr(unsafe.Pointer(&bytesWritten)),
	)
	if ret == 0 {
		return fmt.Errorf("%w: WritePrinter failed: %v", domainErrors.ErrPrintFailed, err)
	}

	return nil
}

func (t *SystemTransport) Close() error {
	if t.handle != 0 {
		procClosePrinter.Call(uintptr(t.handle))
		t.handle = 0
	}
	return nil
}
