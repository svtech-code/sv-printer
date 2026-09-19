package job

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"time"
)

type PrintJobStatus string

const (
	StatusQueued    PrintJobStatus = "queued"
	StatusPrinting  PrintJobStatus = "printing"
	StatusCompleted PrintJobStatus = "completed"
	StatusFailed    PrintJobStatus = "failed"
	StatusCancelled PrintJobStatus = "cancelled"
)

type PrintJob struct {
	ID        string         `json:"id"`
	PrinterID string         `json:"printer_id"`
	Payload   []byte         `json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	Status    PrintJobStatus `json:"status"`
}

func NewID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "job_" + strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	return "job_" + hex.EncodeToString(b)
}
