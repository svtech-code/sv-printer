package autostart

import (
	"fmt"
	"os"
	"path/filepath"
)

func getPlistPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Library", "LaunchAgents", fmt.Sprintf("com.svtech.%s.plist", appName)), nil
}

func Install() error {
	path, err := getExecPath()
	if err != nil {
		return err
	}
	plistPath, err := getPlistPath()
	if err != nil {
		return err
	}

	content := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>com.svtech.%s</string>
	<key>ProgramArguments</key>
	<array>
		<string>%s</string>
	</array>
	<key>RunAtLoad</key>
	<true/>
</dict>
</plist>`, appName, path)

	if err := os.MkdirAll(filepath.Dir(plistPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(plistPath, []byte(content), 0644)
}

func Uninstall() error {
	plistPath, err := getPlistPath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(plistPath); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(plistPath)
}

func Status() (bool, error) {
	plistPath, err := getPlistPath()
	if err != nil {
		return false, err
	}
	_, err = os.Stat(plistPath)
	return err == nil, nil
}
