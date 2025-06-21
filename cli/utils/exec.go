package utils

import (
	"os"
	"os/exec"
)

func Exec(command string, envs ...string) (string, error) {
	cmd := exec.Command(command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	for _, env := range envs {
		cmd.Env = append(os.Environ(), env)
	}
	output, err := cmd.CombinedOutput()
	return string(output), err
}
