package http

import (
	"encoding/json"
	"errors"
	"net/http"

	domainErrors "sv-print/internal/domain/errors"
)

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error errorBody `json:"error"`
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorResponse{Error: errorBody{Code: code, Message: message}})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func decodeBody(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			s, c, m := mapDomainError(domainErrors.ErrPayloadTooLarge)
			writeError(w, s, c, m)
			return false
		}
		s, c, m := mapDomainError(domainErrors.ErrInvalidPayload)
		writeError(w, s, c, m)
		return false
	}
	return true
}

func mapDomainError(err error) (int, string, string) {
	switch {
	case errors.Is(err, domainErrors.ErrPrinterNotFound):
		return http.StatusNotFound, "PRINTER_NOT_FOUND", "Printer not found."
	case errors.Is(err, domainErrors.ErrPrinterOffline):
		return http.StatusServiceUnavailable, "PRINTER_OFFLINE", "Printer is currently offline."
	case errors.Is(err, domainErrors.ErrPrinterBusy):
		return http.StatusConflict, "PRINTER_BUSY", "Printer is busy."
	case errors.Is(err, domainErrors.ErrPrinterPaperOut):
		return http.StatusConflict, "PRINTER_PAPER_OUT", "Printer is out of paper."
	case errors.Is(err, domainErrors.ErrPrinterConnectionFailed):
		return http.StatusBadGateway, "PRINTER_CONNECTION_FAILED", "Could not connect to the printer."
	case errors.Is(err, domainErrors.ErrPrintFailed):
		return http.StatusInternalServerError, "PRINT_FAILED", "Printing failed."
	case errors.Is(err, domainErrors.ErrInvalidPayload):
		return http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload."
	case errors.Is(err, domainErrors.ErrUnauthorized):
		return http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized."
	case errors.Is(err, domainErrors.ErrOriginNotAllowed):
		return http.StatusForbidden, "ORIGIN_NOT_ALLOWED", "Origin not allowed."
	case errors.Is(err, domainErrors.ErrPayloadTooLarge):
		return http.StatusRequestEntityTooLarge, "PAYLOAD_TOO_LARGE", "Payload too large."
	case errors.Is(err, domainErrors.ErrUnsupportedProtocol):
		return http.StatusBadRequest, "UNSUPPORTED_PROTOCOL", "Unsupported protocol."
	case errors.Is(err, domainErrors.ErrJobNotFound):
		return http.StatusNotFound, "JOB_NOT_FOUND", "Job not found."
	case errors.Is(err, domainErrors.ErrJobAlreadyExists):
		return http.StatusConflict, "JOB_ALREADY_EXISTS", "Job already exists."
	case errors.Is(err, domainErrors.ErrLicenseRequired):
		return http.StatusForbidden, "LICENSE_REQUIRED", "A valid license is required for this feature."
	case errors.Is(err, domainErrors.ErrLicenseQuotaExceeded):
		return http.StatusTooManyRequests, "LICENSE_QUOTA_EXCEEDED", "Daily trial print limit reached."
	default:
		return http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error."
	}
}
