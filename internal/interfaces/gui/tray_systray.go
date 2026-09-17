//go:build (darwin && cgo) || (linux && cgo) || windows


package gui

import (
	"github.com/atotto/clipboard"
	"github.com/getlantern/systray"
)

// RunTray starts the system tray on the main OS thread.
// It calls onReady when the tray is drawn, which should start the background agent.
// It calls onExit when the user clicks "Salir" to shut down the agent.
func RunTray(token string, onReady func(), onExit func()) {
	var isHiding bool

	systray.Run(
		func() {
			systray.SetTitle("SV Printer")
			systray.SetTooltip("SV Printer Agent")

			mCopyToken := systray.AddMenuItem("Copiar Token", "Copia el token de seguridad al portapapeles")
			mHide := systray.AddMenuItem("Ocultar Icono", "Oculta el icono de la barra (el agente seguirá corriendo)")
			systray.AddSeparator()
			mQuit := systray.AddMenuItem("Salir", "Cerrar el agente SV Printer")

			// Start the background process
			go onReady()

			// Listen for menu clicks
			for {
				select {
				case <-mCopyToken.ClickedCh:
					clipboard.WriteAll(token)
				case <-mHide.ClickedCh:
					isHiding = true
					systray.Quit()
					return
				case <-mQuit.ClickedCh:
					isHiding = false
					systray.Quit()
					return
				}
			}
		},
		func() {
			// Trigger the shutdown callback only if we are actually quitting
			if !isHiding {
				onExit()
			}
		},
	)
}
