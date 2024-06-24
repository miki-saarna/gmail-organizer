package utils

import (
	"fmt"
	"os/exec"
	"runtime"
)

func OpenURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)

	return exec.Command(cmd, args...).Start()
}

func main() {
	url := "https://www.example.com"
	err := OpenURL(url)
	if err != nil {
		fmt.Println("Failed to open URL:", err)
	}
}
