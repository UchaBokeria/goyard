package utils

import (
	"bufio"
	"fmt"
	"os/exec"
)

func Exec(command string, envs ...string) (string, error) {
	cmd := exec.Command(command, envs...)

	// Get stdout pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return "", err
	}

	// Read output line by line and print with command prefix
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Printf("[%s] %s\n", command, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	// Wait for command to finish
	if err := cmd.Wait(); err != nil {
		return "", err
	}

	return "", nil
}
