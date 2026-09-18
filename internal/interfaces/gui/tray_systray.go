//go:build (darwin && cgo) || windows

package gui

import (
	"log/slog"

	"github.com/atotto/clipboard"
	"github.com/getlantern/systray"
	"sv-printer/internal/autostart"
)

// RunTray starts the system tray on the main OS thread.
// It calls onReady when the tray is drawn, which should start the background agent.
// It calls onExit when the user clicks "Salir" to shut down the agent.
func RunTray(token string, onReady func(), onExit func()) {
	systray.Run(
		func() {
			systray.SetIcon(iconData)
			// systray.SetTitle("SV Printer") // Removed so only the logo appears
			systray.SetTooltip("SV Printer Agent")

			mCopyToken := systray.AddMenuItem("Copiar Token", "Copia el token de seguridad al portapapeles")
			systray.AddSeparator()

			// Initial autostart status
			enabled, _ := autostart.Status()
			mAutostart := systray.AddMenuItemCheckbox("Iniciar con el sistema", "Ejecuta SV Printer al encender la computadora", enabled)

			systray.AddSeparator()
			mQuit := systray.AddMenuItem("Salir", "Cerrar el agente SV Printer")

			// Start the background process
			go onReady()

			// Listen for menu clicks
			for {
				select {
				case <-mCopyToken.ClickedCh:
					clipboard.WriteAll(token)
				case <-mAutostart.ClickedCh:
					if mAutostart.Checked() {
						if err := autostart.Uninstall(); err == nil {
							mAutostart.Uncheck()
						} else {
							slog.Error("failed to disable autostart", "error", err)
						}
					} else {
						if err := autostart.Install(); err == nil {
							mAutostart.Check()
						} else {
							slog.Error("failed to enable autostart", "error", err)
						}
					}
				case <-mQuit.ClickedCh:
					systray.Quit()
					return
				}
			}
		},
		func() {
			onExit()
		},
	)
}
