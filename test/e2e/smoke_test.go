//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// startTCPServer starts a TCP listener on a random port and returns the
// captured bytes and the address. The listener accepts connections until
// stop is called.
func startTCPServer(t *testing.T) (addr string, captured *[][]byte, stop func()) {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	var mu struct {
		bufs [][]byte
	}
	captured = &mu.bufs

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			buf, _ := io.ReadAll(conn)
			conn.Close()
			mu.bufs = append(mu.bufs, buf)
		}
	}()

	return ln.Addr().String(), captured, func() { ln.Close() }
}

// waitForHTTP polls the health endpoint until it responds or timeout.
func waitForHTTP(t *testing.T, baseURL, token string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		req, _ := http.NewRequest("GET", baseURL+"/health", nil)
		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("agent not ready after %v", timeout)
}

func TestSmoke(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e in short mode")
	}

	repoRoot, _ := os.Getwd()
	repoRoot = filepath.Join(repoRoot, "..", "..")

	// Build the agent binary.
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "sv-print")

	goBuild := exec.Command("go", "build", "-o", binPath, "./cmd/sv-print")
	goBuild.Dir = repoRoot
	if out, err := goBuild.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}

	// Start a simulated printer.
	tcpAddr, captured, tcpStop := startTCPServer(t)
	defer tcpStop()

	// Find a free port for the agent HTTP API.
	apiLn, _ := net.Listen("tcp", "127.0.0.1:0")
	apiAddr := apiLn.Addr().String()
	apiLn.Close()
	_, apiPort, _ := net.SplitHostPort(apiAddr)

	token := "e2e-test-token"
	printerID := fmt.Sprintf("Sim@%s", tcpAddr)

	cmd := exec.Command(binPath,
		"-token", token,
		"-port", apiPort,
		"-printer", printerID,
	)
	cmd.Dir = repoRoot

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout

	if err := cmd.Start(); err != nil {
		t.Fatalf("start agent: %v", err)
	}
	defer func() {
		cmd.Process.Signal(os.Interrupt)
		done := make(chan error)
		go func() { done <- cmd.Wait() }()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			cmd.Process.Kill()
		}
	}()

	baseURL := fmt.Sprintf("http://127.0.0.1:%s", apiPort)
	waitForHTTP(t, baseURL, token, 10*time.Second)

	client := &http.Client{Timeout: 5 * time.Second}

	// --- /health (no auth) ---
	resp, err := client.Get(baseURL + "/health")
	if err != nil {
		t.Fatalf("GET /health: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Errorf("/health status = %d, want 200; body: %s", resp.StatusCode, body)
	}

	// --- /api/v1/info ---
	resp = doAuth(t, client, "GET", baseURL+"/api/v1/info", token, nil)
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Errorf("/info status = %d, want 200", resp.StatusCode)
	}
	var info map[string]any
	json.Unmarshal(body, &info)
	if info["name"] != "sv-print" {
		t.Errorf("info.name = %v, want sv-print", info["name"])
	}
	if info["tier"] != "trial" {
		t.Errorf("info.tier = %v, want trial", info["tier"])
	}

	// --- /api/v1/printers ---
	resp = doAuth(t, client, "GET", baseURL+"/api/v1/printers", token, nil)
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	var printersResp struct {
		Printers []struct {
			ID string `json:"id"`
		} `json:"printers"`
	}
	json.Unmarshal(body, &printersResp)
	if len(printersResp.Printers) == 0 {
		t.Error("expected at least one printer")
	}

	// --- POST /api/v1/printers/{id}/test (trial watermark) ---
	simID := "net-" + tcpAddr
	resp = doAuth(t, client, "POST", baseURL+"/api/v1/printers/"+simID+"/test", token, nil)
	if resp.StatusCode != 201 {
		t.Errorf("/test status = %d, want 201", resp.StatusCode)
	}
	resp.Body.Close()

	// Wait for job to process.
	time.Sleep(500 * time.Millisecond)

	// Verify captured ESC/POS bytes.
	allBytes := joinBytes(*captured)
	escInit := []byte{0x1B, 0x40}
	cutCmd := []byte{0x1D, 0x56}
	watermark := []byte("*** SV PRINT")

	if !bytes.Contains(allBytes, escInit) {
		t.Error("missing ESC @ (init) in captured bytes")
	}
	if !bytes.Contains(allBytes, cutCmd) {
		t.Error("missing GS V (cut) in captured bytes")
	}
	if !bytes.Contains(allBytes, watermark) {
		t.Errorf("missing trial watermark in captured bytes")
	}

	// --- POST /api/v1/print/receipt (structured receipt) ---
	receiptDoc := map[string]any{
		"printer_id": simID,
		"cut":        true,
		"lines": []map[string]any{
			{"text": "E2E Test", "style": map[string]any{"bold": true, "align": "center"}},
			{"text": "Line 2"},
		},
	}
	resp = doAuth(t, client, "POST", baseURL+"/api/v1/print/receipt", token, receiptDoc)
	if resp.StatusCode != 201 {
		t.Errorf("/print/receipt status = %d, want 201", resp.StatusCode)
	}
	resp.Body.Close()
	time.Sleep(500 * time.Millisecond)

	// Re-capture after second receipt.
	allBytes = joinBytes(*captured)

	// Verify additional bytes contain "E2E Test".
	if !bytes.Contains(allBytes, []byte("E2E Test")) {
		t.Error("missing 'E2E Test' in captured bytes")
	}

	// --- License gating: POST /api/v1/print (raw, pro) → 403 ---
	resp = doAuth(t, client, "POST", baseURL+"/api/v1/print", token, map[string]any{
		"printer_id": simID,
		"payload":    "AQIDBA==",
		"format":     "escpos",
	})
	if resp.StatusCode != 403 {
		t.Errorf("/print (raw, trial) status = %d, want 403", resp.StatusCode)
	}
	resp.Body.Close()

	// --- License gating: GET /api/v1/events (ws, pro) → 403 ---
	resp = doAuth(t, client, "GET", baseURL+"/api/v1/events?token="+token, token, nil)
	if resp.StatusCode != 403 {
		t.Errorf("/events (trial) status = %d, want 403", resp.StatusCode)
	}
	resp.Body.Close()
}

func doAuth(t *testing.T, client *http.Client, method, url, token string, body any) *http.Response {
	t.Helper()

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("do %s %s: %v", method, url, err)
	}
	return resp
}

func joinBytes(bufs [][]byte) []byte {
	var total int
	for _, b := range bufs {
		total += len(b)
	}
	out := make([]byte, 0, total)
	for _, b := range bufs {
		out = append(out, b...)
	}
	return out
}

func TestSmoke_CLISubcommands(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e in short mode")
	}

	repoRoot, _ := os.Getwd()
	repoRoot = filepath.Join(repoRoot, "..", "..")

	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "sv-print")

	goBuild := exec.Command("go", "build", "-o", binPath, "./cmd/sv-print")
	goBuild.Dir = repoRoot
	if out, err := goBuild.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}

	tests := []struct {
		name    string
		args    []string
		wantOut string
	}{
		{"version", []string{"version"}, "SV Print v0.1.0"},
		{"status", []string{"status"}, "SV Print status: OK"},
		{"device-id", []string{"device-id"}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := exec.Command(binPath, tt.args...).CombinedOutput()
			if err != nil {
				t.Fatalf("sv-print %s: %v\n%s", tt.args[0], err, out)
			}
			if tt.wantOut != "" && !strings.Contains(string(out), tt.wantOut) {
				t.Errorf("output = %q, want to contain %q", out, tt.wantOut)
			}
		})
	}
}
