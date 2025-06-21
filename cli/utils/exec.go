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
		fmt.Printf("Error getting stdout pipe for %s: %v\n", name, err)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("Error getting stderr pipe: %v\n", err)
		return
	}
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Printf("[%s][ERR] %s\n", name, scanner.Text())
		}
	}()
	// Start the command
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting %s: %v\n", name, err)
		return
	}

	// Read output line by line and print with command prefix
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Printf("[%s] %s\n", name, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading stdout for %s: %v\n", name, err)
	}

	// Wait for command to finish
	if err := cmd.Wait(); err != nil {
		fmt.Printf("Command %s finished with error: %v\n", name, err)
	}
}
