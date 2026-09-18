package autostart

import (
	"golang.org/x/sys/windows/registry"
)

func Install() error {
	path, err := getExecPath()
	if err != nil {
		return err
	}
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()
	return k.SetStringValue(appName, path)
}

func Uninstall() error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()
	return k.DeleteValue(appName)
}

func Status() (bool, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE)
	if err != nil {
		return false, nil
	}
	defer k.Close()
	_, _, err = k.GetStringValue(appName)
	if err != nil {
		return false, nil
	}
	return true, nil
}
