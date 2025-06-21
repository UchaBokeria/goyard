package utils

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func ExecLive(name, command string) (<-chan string, <-chan error) {
	outChan := make(chan string)
	errChan := make(chan error, 1)

	cmdArgs := strings.Fields(command)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		errChan <- fmt.Errorf("%s stdout pipe: %w", name, err)
		close(outChan)
		close(errChan)
		return outChan, errChan
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		errChan <- fmt.Errorf("%s stderr pipe: %w", name, err)
		close(outChan)
		close(errChan)
		return outChan, errChan
	}

	go func() {
		defer close(outChan)
		defer close(errChan)

		if err := cmd.Start(); err != nil {
			errChan <- fmt.Errorf("%s start error: %w", name, err)
			return
		}

		scanner := bufio.NewScanner(io.MultiReader(stdoutPipe, stderrPipe))
		for scanner.Scan() {
			outChan <- fmt.Sprintf("[%s] %s", name, scanner.Text())
		}

		if err := cmd.Wait(); err != nil {
			errChan <- fmt.Errorf("%s wait error: %w", name, err)
		}
	}()

	return outChan, errChan
}
