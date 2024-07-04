package utils

import (
	"encoding/json"
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

func PrettyPrint(data interface{}) (string, error) {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", fmt.Errorf("error marshalling json: %v", err.Error())
	}
	return fmt.Sprintln(string(jsonData)), nil
}
