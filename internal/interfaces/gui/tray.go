package gui

import (

	"github.com/atotto/clipboard"
	"github.com/getlantern/systray"
)

// RunTray starts the system tray on the main OS thread.
// It calls onReady when the tray is drawn, which should start the background agent.
// It calls onExit when the user clicks "Salir" to shut down the agent.
func RunTray(token string, onReady func(), onExit func()) {
	systray.Run(
		func() {
			systray.SetTitle("SV Printer")
			systray.SetTooltip("SV Printer Agent")
			
			// We don't have an icon byte slice right now, so it will just show the Title on Mac,
			// or a default blank square on Windows. In production you would call systray.SetIcon(iconBytes).

			mCopyToken := systray.AddMenuItem("Copiar Token", "Copia el token de seguridad al portapapeles")
			systray.AddSeparator()
			mQuit := systray.AddMenuItem("Salir", "Cerrar el agente SV Printer")

			// Start the background process
			go onReady()

			// Listen for menu clicks
			for {
				select {
				case <-mCopyToken.ClickedCh:
					clipboard.WriteAll(token)
				case <-mQuit.ClickedCh:
					systray.Quit()
					return
				}
			}
		},
		func() {
			// Trigger the shutdown callback when the tray exits
			onExit()
		},
	)
}
