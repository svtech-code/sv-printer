package autostart

import (
	"fmt"
	"os"

	"github.com/emersion/go-autostart"
)

// App defines the application properties for autostart
var App = &autostart.App{
	Name:        "sv-printer",
	DisplayName: "SV Printer Agent",
}

// Setup initializes the Exec path for the app
func Setup() error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}
	App.Exec = []string{execPath}
	return nil
}

// Install enables autostart
func Install() error {
	if err := Setup(); err != nil {
		return err
	}
	if App.IsEnabled() {
		return fmt.Errorf("autostart is already enabled")
	}
	return App.Enable()
}

// Uninstall disables autostart
func Uninstall() error {
	if err := Setup(); err != nil {
		return err
	}
	if !App.IsEnabled() {
		return fmt.Errorf("autostart is not enabled")
	}
	return App.Disable()
}

// Status checks if autostart is enabled
func Status() (bool, error) {
	if err := Setup(); err != nil {
		return false, err
	}
	return App.IsEnabled(), nil
}
