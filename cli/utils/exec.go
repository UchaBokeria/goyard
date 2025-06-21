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
	cmd.Stderr = io.MultiWriter(&errBuf, os.Stderr)
	cmd.Stdout = io.MultiWriter(&outBuf)
	for _, env := range envs {
		cmd.Env = append(os.Environ(), env)
	}
	err := cmd.Run()
	return outBuf.String(), err
}
