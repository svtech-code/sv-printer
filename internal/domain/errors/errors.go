package errors

import "errors"

var (
	ErrPrinterNotFound         = errors.New("PRINTER_NOT_FOUND")
	ErrPrinterOffline          = errors.New("PRINTER_OFFLINE")
	ErrPrinterBusy             = errors.New("PRINTER_BUSY")
	ErrPrinterPaperOut         = errors.New("PRINTER_PAPER_OUT")
	ErrPrinterConnectionFailed = errors.New("PRINTER_CONNECTION_FAILED")
	ErrPrintFailed             = errors.New("PRINT_FAILED")
	ErrInvalidPayload          = errors.New("INVALID_PAYLOAD")
	ErrUnauthorized            = errors.New("UNAUTHORIZED")
	ErrOriginNotAllowed        = errors.New("ORIGIN_NOT_ALLOWED")
	ErrPayloadTooLarge         = errors.New("PAYLOAD_TOO_LARGE")
	ErrUnsupportedProtocol     = errors.New("UNSUPPORTED_PROTOCOL")
	ErrJobNotFound             = errors.New("JOB_NOT_FOUND")
	ErrJobAlreadyExists        = errors.New("JOB_ALREADY_EXISTS")
	ErrLicenseRequired         = errors.New("LICENSE_REQUIRED")
	ErrLicenseQuotaExceeded    = errors.New("LICENSE_QUOTA_EXCEEDED")
)
