//go:build !windows

package system

import (
	"context"
	"os/exec"
	"strings"
	"testing"
)

func TestSystemTransport_Open(t *testing.T) {
	tr := NewSystemTransport("TestPrinter")

	err := tr.Open(context.Background())
	// On Mac/Linux, lp should exist, so err should be nil
	// However, if the environment doesn't have lp, it might fail.
	_, lookErr := exec.LookPath("lp")
	if lookErr == nil && err != nil {
		t.Errorf("Expected Open to succeed, got %v", err)
	} else if lookErr != nil && err == nil {
		t.Errorf("Expected Open to fail since lp is missing, but it succeeded")
	}
}

func TestSystemTransport_Write(t *testing.T) {
	// We can't really test Write easily without a real lp command and printer,
	// but we can check if it returns an error when given an invalid printer name
	// on a system that has lp.
	_, lookErr := exec.LookPath("lp")
	if lookErr != nil {
		t.Skip("Skipping Write test because lp command is not available")
	}

	tr := NewSystemTransport("InvalidPrinterThatDoesNotExist12345")
	err := tr.Write(context.Background(), []byte("test"))
	if err == nil {
		t.Errorf("Expected Write to fail for invalid printer, but it succeeded")
	} else if !strings.Contains(err.Error(), "lp command failed") {
		t.Errorf("Expected lp command failed error, got %v", err)
	}
}

func TestSystemTransport_Close(t *testing.T) {
	tr := NewSystemTransport("TestPrinter")
	if err := tr.Close(); err != nil {
		t.Errorf("Expected Close to succeed, got %v", err)
	}
}
