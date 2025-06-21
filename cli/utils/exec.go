package utils

import (
	"bufio"
	"io"
	"os/exec"
	"strings"
)

func ExecLive(name, command string) (<-chan string, <-chan error) {
	outChan := make(chan string)
	errChan := make(chan error, 1)

	cmdArgs := strings.Fields(command)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	merged := io.MultiReader(stdout, stderr)
	scanner := bufio.NewScanner(merged)

	go func() {
		defer close(outChan)
		defer close(errChan)

		if err := cmd.Start(); err != nil {
			errChan <- err
			return
		}

		for scanner.Scan() {
			outChan <- "[" + name + "] " + scanner.Text()
		}

		if err := cmd.Wait(); err != nil {
			errChan <- err
			return
		}
	}()

	return outChan, errChan
}
