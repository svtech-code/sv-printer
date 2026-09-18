package autostart

import (
	"fmt"
	"os"
	"path/filepath"
)

func getDesktopPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "autostart", fmt.Sprintf("%s.desktop", appName)), nil
}

func Install() error {
	path, err := getExecPath()
	if err != nil {
		return err
	}
	desktopPath, err := getDesktopPath()
	if err != nil {
		return err
	}

	content := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=%s
Exec=%s
Hidden=false
NoDisplay=false
X-GNOME-Autostart-enabled=true
`, displayName, path)

	if err := os.MkdirAll(filepath.Dir(desktopPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(desktopPath, []byte(content), 0644)
}

func Uninstall() error {
	desktopPath, err := getDesktopPath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(desktopPath); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(desktopPath)
}

func Status() (bool, error) {
	desktopPath, err := getDesktopPath()
	if err != nil {
		return false, err
	}
	_, err = os.Stat(desktopPath)
	return err == nil, nil
}
