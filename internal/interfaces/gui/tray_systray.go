//go:build (darwin && cgo) || windows

package gui

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/atotto/clipboard"
	"github.com/getlantern/systray"
	"github.com/sqweek/dialog"
	"sv-printer/internal/autostart"
	"sv-printer/internal/deviceid"
)

// RunTray starts the system tray on the main OS thread.
// It calls onReady when the tray is drawn, which should start the background agent.
// It calls onExit when the user clicks "Salir" to shut down the agent.
func RunTray(token string, configPath string, isPro bool, onReady func(), onExit func()) {
	systray.Run(
		func() {
			systray.SetIcon(iconData)
			systray.SetTooltip("SV Printer Agent")

			// License Status Header
			if isPro {
				mStatus := systray.AddMenuItem("Licencia: PRO", "Aplicación activada")
				mStatus.Disable()
			} else {
				mStatus := systray.AddMenuItem("Licencia: FREE", "Modo de prueba limitado")
				mStatus.Disable()
			}
			systray.AddSeparator()

			mCopyToken := systray.AddMenuItem("Copiar Token", "Copia el token de seguridad al portapapeles")
			mCopyID := systray.AddMenuItem("Copiar ID del Equipo", "Copia el ID único para generar una licencia")
			mInstallLic := systray.AddMenuItem("Instalar Licencia", "Seleccionar archivo license.key")

			if isPro {
				mCopyID.Hide()
				mInstallLic.Hide()
			}
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
					dialog.Message("Token copiado al portapapeles:\n%s", token).Title("Token de Seguridad").Info()
				case <-mCopyID.ClickedCh:
					id, err := deviceid.ID()
					if err == nil {
						clipboard.WriteAll(id)
						dialog.Message("ID copiado al portapapeles:\n%s", id).Title("ID del Equipo").Info()
					} else {
						dialog.Message("No se pudo obtener el ID: %v", err).Title("Error").Error()
					}
				case <-mInstallLic.ClickedCh:
					file, err := dialog.File().Filter("License Key", "key").Title("Seleccionar Licencia").Load()
					if err == nil && file != "" {
						dest := filepath.Join(filepath.Dir(configPath), "license.key")
						if err := copyFile(file, dest); err != nil {
							dialog.Message("Error al instalar licencia: %v", err).Title("Error").Error()
						} else {
							dialog.Message("Licencia instalada correctamente.\n\nPor favor, sal de la aplicación y vuelve a abrirla para aplicar los cambios.").Title("Éxito").Info()
						}
					}
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

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
