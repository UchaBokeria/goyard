package utils

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

func Exec(command string, envs ...string) (string, error) {
	cmd := exec.Command(command)
	var errBuf, outBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
	cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)
	for _, env := range envs {
		cmd.Env = append(os.Environ(), env)
	}
	output, err := cmd.CombinedOutput()
	return string(output), err
}
