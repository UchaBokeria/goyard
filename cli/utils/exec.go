package utils

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"syscall"
)

func ExecLive(name, command string) (<-chan string, <-chan error) {
	return ExecLiveWithContext(context.Background(), name, command)
}

func ExecLiveWithContext(ctx context.Context, name, command string) (<-chan string, <-chan error) {
	outChan := make(chan string, 100) // Buffered to prevent blocking
	errChan := make(chan error, 1)

	cmdArgs := strings.Fields(command)
	cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)

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
		stdoutPipe.Close() // Clean up stdout pipe
		close(outChan)
		close(errChan)
		return outChan, errChan
	}

	go func() {
		defer func() {
			// Ensure pipes are closed
			stdoutPipe.Close()
			stderrPipe.Close()
			close(outChan)
			close(errChan)
		}()

		if err := cmd.Start(); err != nil {
			errChan <- fmt.Errorf("%s start error: %w", name, err)
			return
		}

		// Create a combined reader for both stdout and stderr
		combinedReader := io.MultiReader(stdoutPipe, stderrPipe)
		scanner := bufio.NewScanner(combinedReader)

		// Read output in a separate goroutine to handle context cancellation
		done := make(chan bool)
		go func() {
			defer close(done)
			for scanner.Scan() {
				select {
				case outChan <- fmt.Sprintf("[%s] %s", name, scanner.Text()):
				case <-ctx.Done():
					return
				}
			}

			// Check for scanner errors
			if err := scanner.Err(); err != nil {
				select {
				case errChan <- fmt.Errorf("%s scanner error: %w", name, err):
				case <-ctx.Done():
				}
			}
		}()

		// Wait for either the process to finish or context cancellation
		waitDone := make(chan error, 1)
		go func() {
			waitDone <- cmd.Wait()
		}()

		select {
		case <-ctx.Done():
			// Context cancelled, kill the process group
			if cmd.Process != nil {
				syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
			}
			<-waitDone // Wait for process to actually exit
			errChan <- fmt.Errorf("%s cancelled: %w", name, ctx.Err())
		case err := <-waitDone:
			// Process finished normally
			<-done // Wait for scanner to finish
			if err != nil {
				errChan <- fmt.Errorf("%s wait error: %w", name, err)
			}
		}
	}()

	return outChan, errChan
}
