package main

import (
	"os/exec"
	"runtime"
)

func openURL(url string) error {
	var cmd string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	return exec.Command(cmd, url).Start()
}
