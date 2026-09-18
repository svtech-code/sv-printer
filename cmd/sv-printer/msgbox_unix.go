//go:build !windows

package main

import (
	"fmt"
	"os"
)

func showFatalError(title, msg string) {
	fmt.Fprintf(os.Stderr, "FATAL: %s\n%s\n", title, msg)
}
