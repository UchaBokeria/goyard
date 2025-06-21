package utils

import (
	"bufio"
	"fmt"
	"os/exec"
)

func Exec(name string, args ...string) {
	fmt.Printf("Running command: %s %v\n", name, args)
	cmd := exec.Command(name, args...)

	// Get stdout pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return
	}
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
		}
	}()
	// Start the command
	if err := cmd.Start(); err != nil {
		return
	}

	// Read output line by line and print with command prefix
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
	}

	if err := scanner.Err(); err != nil {
	}

	// Wait for command to finish
	if err := cmd.Wait(); err != nil {
	}
}
