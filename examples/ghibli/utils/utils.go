package utils

import (
	"os/exec"
	"regexp"
	"runtime"

	uuid "github.com/satori/go.uuid"
)

// NewUUID создание нового строкового UUID
func NewUUID() string {
	return uuid.Must(uuid.NewV4(), nil).String()
}

// сравнение строки с шаблоном uuid
func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

// open opens the specified URL in the default browser of the user.
func OpenInBrowser(url string) error {
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
