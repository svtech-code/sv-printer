package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunCLI(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errContains string
		wantOut     string
	}{
		{
			name:    "version command",
			args:    []string{"version"},
			wantErr: false,
			wantOut: "SV Print Agent v0.1.0\n",
		},
		{
			name:    "status command",
			args:    []string{"status"},
			wantErr: false,
			wantOut: "SV Print Agent status: OK\n",
		},
		{
			name:    "printers command with no printers",
			args:    []string{"printers"},
			wantErr: false,
			wantOut: "No printers configured yet.\n",
		},
		{
			name:        "unknown subcommand",
			args:        []string{"foo"},
			wantErr:     true,
			errContains: "unknown subcommand: foo",
		},
		{
			name:        "unimplemented subcommand",
			args:        []string{"service"},
			wantErr:     true,
			errContains: "service is not implemented yet",
		},
		{
			name:        "test requires printer id",
			args:        []string{"test"},
			wantErr:     true,
			errContains: "usage: sv-print test",
		},
		{
			name:        "print requires file",
			args:        []string{"print"},
			wantErr:     true,
			errContains: "usage: sv-print print",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var outBuf bytes.Buffer
			err := runCLI(context.Background(), tt.args, &outBuf)

			if (err != nil) != tt.wantErr {
				t.Errorf("runCLI() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("runCLI() error = %v, want err to contain %v", err, tt.errContains)
				}
			}
			if tt.wantOut != "" && outBuf.String() != tt.wantOut {
				t.Errorf("runCLI() output = %q, want %q", outBuf.String(), tt.wantOut)
			}
		})
	}
}

func TestRunCLIPrintersList(t *testing.T) {
	var outBuf bytes.Buffer
	err := runCLI(context.Background(), []string{"printers", "-printer", "Caja 1@192.168.1.100:9100"}, &outBuf)
	if err != nil {
		t.Fatalf("runCLI() error = %v", err)
	}

	out := outBuf.String()
	if !strings.Contains(out, "net-192.168.1.100:9100") {
		t.Errorf("runCLI() output = %q, want to contain printer ID", out)
	}
	if !strings.Contains(out, "Caja 1") {
		t.Errorf("runCLI() output = %q, want to contain printer name", out)
	}
}

func TestRunCLIDiscover(t *testing.T) {
	var outBuf bytes.Buffer
	err := runCLI(context.Background(), []string{"discover", "-printer", "Caja 1@192.168.1.100:9100"}, &outBuf)
	if err != nil {
		t.Fatalf("runCLI() error = %v", err)
	}
	if !strings.Contains(outBuf.String(), "net-192.168.1.100:9100") {
		t.Errorf("runCLI() output = %q, want to contain discovered printer ID", outBuf.String())
	}
}

func TestRunCLIDoctor(t *testing.T) {
	var outBuf bytes.Buffer
	err := runCLI(context.Background(), []string{"doctor"}, &outBuf)
	if err != nil {
		t.Fatalf("runCLI() error = %v", err)
	}
	if !strings.Contains(outBuf.String(), "SV Print Doctor") {
		t.Errorf("runCLI() output = %q, want to contain doctor banner", outBuf.String())
	}
}

func TestRunCLILogsNoFile(t *testing.T) {
	var outBuf bytes.Buffer
	err := runCLI(context.Background(), []string{"logs"}, &outBuf)
	if err != nil {
		t.Fatalf("runCLI() error = %v", err)
	}
	if !strings.Contains(outBuf.String(), "No logs found") {
		t.Errorf("runCLI() output = %q, want no-logs message", outBuf.String())
	}
}
