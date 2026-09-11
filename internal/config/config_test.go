package config

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	cfg, err := Load(nil)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Host != DefaultHost {
		t.Errorf("Host = %q, want %q", cfg.Host, DefaultHost)
	}
	if cfg.Port != DefaultPort {
		t.Errorf("Port = %d, want %d", cfg.Port, DefaultPort)
	}
	if cfg.MaxPayloadSize != DefaultMaxPayloadSize {
		t.Errorf("MaxPayloadSize = %d, want %d", cfg.MaxPayloadSize, DefaultMaxPayloadSize)
	}
	if cfg.Token == "" {
		t.Error("Token should be generated when not provided")
	}
	if !cfg.TokenGenerated {
		t.Error("TokenGenerated should be true when no token is provided")
	}
}

func TestLoadFlags(t *testing.T) {
	cfg, err := Load([]string{
		"-host", "0.0.0.0",
		"-port", "9000",
		"-token", "abc123",
		"-origin", "https://app.svtech.cl",
		"-origin", "https://localhost:3000",
		"-max-payload", "1024",
		"-printer", "Caja 1@192.168.1.100:9100",
	})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Host != "0.0.0.0" {
		t.Errorf("Host = %q, want %q", cfg.Host, "0.0.0.0")
	}
	if cfg.Port != 9000 {
		t.Errorf("Port = %d, want %d", cfg.Port, 9000)
	}
	if cfg.Token != "abc123" {
		t.Errorf("Token = %q, want %q", cfg.Token, "abc123")
	}
	if cfg.TokenGenerated {
		t.Error("TokenGenerated should be false when token is provided")
	}
	if len(cfg.AllowedOrigins) != 2 {
		t.Errorf("AllowedOrigins len = %d, want 2", len(cfg.AllowedOrigins))
	}
	if cfg.MaxPayloadSize != 1024 {
		t.Errorf("MaxPayloadSize = %d, want 1024", cfg.MaxPayloadSize)
	}
	if len(cfg.Printers) != 1 {
		t.Fatalf("Printers len = %d, want 1", len(cfg.Printers))
	}

	p := cfg.Printers[0]
	if p.Name != "Caja 1" {
		t.Errorf("Printer.Name = %q, want %q", p.Name, "Caja 1")
	}
	if p.Address != "192.168.1.100:9100" {
		t.Errorf("Printer.Address = %q, want %q", p.Address, "192.168.1.100:9100")
	}
	if p.ID != "net-192.168.1.100:9100" {
		t.Errorf("Printer.ID = %q, want %q", p.ID, "net-192.168.1.100:9100")
	}
	if p.Protocol != "escpos" {
		t.Errorf("Printer.Protocol = %q, want %q", p.Protocol, "escpos")
	}
}

func TestLoadInvalidPrinter(t *testing.T) {
	_, err := Load([]string{"-printer", "@"})
	if err == nil {
		t.Fatal("Load() expected error for printer with empty address")
	}
}

func TestLoadSerialPrinter(t *testing.T) {
	cfg, err := Load([]string{"-serial", "Receipt@/dev/ttyUSB0@9600"})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(cfg.Printers) != 1 {
		t.Fatalf("Printers len = %d, want 1", len(cfg.Printers))
	}
	p := cfg.Printers[0]
	if p.Connection != "serial" {
		t.Errorf("Connection = %q, want serial", p.Connection)
	}
	if p.Address != "/dev/ttyUSB0" {
		t.Errorf("Address = %q, want /dev/ttyUSB0", p.Address)
	}
	if p.ID != "serial-/dev/ttyUSB0" {
		t.Errorf("ID = %q, want serial-/dev/ttyUSB0", p.ID)
	}
	if p.BaudRate != 9600 {
		t.Errorf("BaudRate = %d, want 9600", p.BaudRate)
	}
}

func TestLoadSerialPrinterDefaultBaud(t *testing.T) {
	cfg, err := Load([]string{"-serial", "/dev/ttyUSB0"})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Printers[0].BaudRate != 9600 {
		t.Errorf("BaudRate = %d, want default 9600", cfg.Printers[0].BaudRate)
	}
}

func TestLoadSerialPrinterInvalidBaud(t *testing.T) {
	_, err := Load([]string{"-serial", "Receipt@/dev/ttyUSB0@fast"})
	if err == nil {
		t.Fatal("Load() expected error for invalid baud rate")
	}
}

func TestSaveAndLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")

	cfg := Config{
		Host:           "127.0.0.1",
		Port:           9000,
		Token:          "persisted-token",
		AllowedOrigins: []string{"https://app.svtech.cl"},
		MaxPayloadSize: 2048,
		Printers: []ManualPrinter{
			{ID: "net-1.2.3.4:9100", Name: "Caja", Address: "1.2.3.4:9100", Protocol: "escpos"},
		},
		Path: path,
	}
	if err := Save(cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := Load([]string{"-config", path})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.Port != 9000 {
		t.Errorf("Port = %d, want 9000", loaded.Port)
	}
	if loaded.Token != "persisted-token" {
		t.Errorf("Token = %q, want persisted-token", loaded.Token)
	}
	if loaded.TokenGenerated {
		t.Error("TokenGenerated should be false when token comes from file")
	}
	if loaded.MaxPayloadSize != 2048 {
		t.Errorf("MaxPayloadSize = %d, want 2048", loaded.MaxPayloadSize)
	}
	if len(loaded.Printers) != 1 {
		t.Fatalf("Printers len = %d, want 1", len(loaded.Printers))
	}
	if loaded.Printers[0].Address != "1.2.3.4:9100" {
		t.Errorf("Printer.Address = %q", loaded.Printers[0].Address)
	}
}

func TestFlagOverridesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")

	if err := Save(Config{Path: path, Port: 9876, Token: "file-token"}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := Load([]string{"-config", path, "-port", "9000"})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.Port != 9000 {
		t.Errorf("Port = %d, want 9000 (flag overrides file)", loaded.Port)
	}
	if loaded.Token != "file-token" {
		t.Errorf("Token = %q, want file-token (from file)", loaded.Token)
	}
}

func TestLoadMissingFile(t *testing.T) {
	cfg, err := Load([]string{"-config", filepath.Join(t.TempDir(), "nope.json")})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Port != DefaultPort {
		t.Errorf("Port = %d, want default %d", cfg.Port, DefaultPort)
	}
}

func TestDefaultPath(t *testing.T) {
	p := DefaultPath()
	if p == "" {
		t.Fatal("DefaultPath() is empty")
	}
	if !strings.HasSuffix(p, "config.json") {
		t.Errorf("DefaultPath() = %q, want suffix config.json", p)
	}
}
