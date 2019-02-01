package exec

import (
	"os/exec"
)

// ExecuteCommand is used to execute a linux command line command and return the output of the command with an error if it exists.
func ExecuteCommand(cmd string, args []string) (string, error) {
	out, err := exec.Command(cmd, args...).Output()
	return string(out), err
}
