package bashutil

import (
	"os"
	"os/exec"
)

// SyncRun runs bash command synchronously and returns output and error
func SyncRun(cmd string) (string, error) {
	b, err := buildBashCmd(cmd).Output()
	if err != nil {
		return string(b), err
	}
	return string(b), err
}

func buildBashCmd(cmdline string) *exec.Cmd {
	cmd := exec.Command("bash", "-c", cmdline)
	cmd.Env = os.Environ()
	return cmd
}
