package utils

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
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

func ExecLive(name, command string) (<-chan string, <-chan error) {
	return ExecLiveWithContext(context.Background(), name, command)
}

func ExecLiveWithContext(ctx context.Context, name, command string) (<-chan string, <-chan error) {
	outChan := make(chan string, 100) // Buffered to prevent blocking
	errChan := make(chan error, 1)

	// Use shell to properly handle quoted arguments
	cmd := exec.CommandContext(ctx, "sh", "-c", command)

	// Set process group to allow killing child processes
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

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
		defer func() {
			close(outChan)
			close(errChan)
		}()

		if err := cmd.Start(); err != nil {
			errChan <- fmt.Errorf("%s start error: %w", name, err)
			return
		}

		// Handle stdout and stderr concurrently instead of using MultiReader
		var wg sync.WaitGroup

		// Handle stdout
		wg.Add(1)
		go func() {
			defer wg.Done()
			scanner := bufio.NewScanner(stdoutPipe)
			for scanner.Scan() {
				select {
				case outChan <- fmt.Sprintf("[%s] %s", name, scanner.Text()):
				case <-ctx.Done():
					return
				}
			}
			if err := scanner.Err(); err != nil {
				select {
				case errChan <- fmt.Errorf("%s stdout scanner error: %w", name, err):
				case <-ctx.Done():
				}
			}
		}()

		// Handle stderr
		wg.Add(1)
		go func() {
			defer wg.Done()
			scanner := bufio.NewScanner(stderrPipe)
			for scanner.Scan() {
				select {
				case outChan <- fmt.Sprintf("[%s] %s", name, scanner.Text()):
				case <-ctx.Done():
					return
				}
			}
			if err := scanner.Err(); err != nil {
				select {
				case errChan <- fmt.Errorf("%s stderr scanner error: %w", name, err):
				case <-ctx.Done():
				}
			}
		}()

		// Wait for either the process to finish or context cancellation
		waitDone := make(chan error, 1)
		go func() {
			waitDone <- cmd.Wait()
		}()

		// Wait for output readers to finish
		outputDone := make(chan struct{})
		go func() {
			wg.Wait()
			close(outputDone)
		}()

		// Wait for either the process to finish or context cancellation
		select {
		case err := <-waitDone:
			if err != nil {
				errChan <- fmt.Errorf("%s execution error: %w", name, err)
			}
			// Wait for output readers to finish processing remaining output
			<-outputDone
		case <-ctx.Done():
			errChan <- ctx.Err()
		}
	}()

	return outChan, errChan
}
