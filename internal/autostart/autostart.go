package autostart

import "os"

const (
	appName     = "sv-printer"
	displayName = "SV Printer"
)

func getExecPath() (string, error) {
	return os.Executable()
}
