//go:build (darwin && !cgo) || linux

package gui

import (
	"log/slog"
)

// RunTray starts the background agent without a system tray
// because CGO is disabled on this platform build.
func RunTray(token string, onReady func(), onExit func()) {
	slog.Info("System tray GUI is disabled in this build (CGO disabled). Running in headless mode.")

	// Start the background process
	go onReady()

	// Since there is no GUI loop to block, we just block forever.
	// The user must kill the process (Ctrl+C) to trigger graceful shutdown via signal context.
	select {}
}
