//go:build !windows

package serial

func patterns() []string {
	return []string{
		"/dev/ttyUSB*",
		"/dev/ttyACM*",
		"/dev/cu.usbserial*",
		"/dev/cu.usbmodem*",
		"/dev/tty.usbserial*",
		"/dev/tty.usbmodem*",
	}
}
