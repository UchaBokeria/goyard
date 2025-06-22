package run

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/UchaBokeria/goyard/cli/utils"
)

var (
	Host string
	Port int
)

func Dev() {
	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Define your development processes
	processes := []struct {
		Name string
		Cmd  string
	}{
		{
			"AIR",
			`go run github.com/air-verse/air@latest --build.cmd "go build -o ./bin/app ./cmd/app" --build.bin "./bin/app" --build.delay "100" --build.exclude_dir "node_modules" --build.include_ext "go, templ" --build.stop_on_error "false" --log.time "true" --log.main_only "true"`,
		},
		{
			"TEMPL",
			fmt.Sprintf(`go run github.com/a-h/templ/cmd/templ@latest generate --open-browser=false --watch --proxy="http://%s:%d" --proxyport="%d" --proxybind="%s"`, Host, Port-1, Port, Host),
		},
		{
			"TAILWIND",
			`bunx --yes tailwindcss -i ./public/assets/styles/tailwind.css -o ./public/assets/styles/style.css --watch`,
		},
	}

	var wg sync.WaitGroup
	processErrors := make(chan error, len(processes))

	// Start all processes using ExecLive
	for _, p := range processes {
		wg.Add(1)
		go func(name, cmd string) {
			defer wg.Done()

			// Use ExecLiveWithContext for cancellable execution
			outChan, errChan := utils.ExecLiveWithContext(ctx, name, cmd)

			// Handle output and errors in real-time
			for {
				select {
				case line, ok := <-outChan:
					if !ok {
						return // Channel closed, process finished
					}
					fmt.Println(line) // Print live output
				case err, ok := <-errChan:
					if !ok {
						return // Channel closed
					}
					if err != nil {
						fmt.Fprintf(os.Stderr, "[%s] Error: %v\n", name, err)
						// Send error to main goroutine
						select {
						case processErrors <- err:
						default:
						}
						return
					}
				case <-ctx.Done():
					return // Context cancelled
				}
			}
		}(p.Name, p.Cmd)
	}

	// Wait for either all processes to finish, a signal, or an error
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-sigChan:
		fmt.Println("\n🛑 Received shutdown signal, stopping all processes...")
		cancel() // This will stop all ExecLive processes
		<-done   // Wait for all goroutines to finish
		fmt.Println("✅ All processes stopped gracefully")
	case err := <-processErrors:
		fmt.Fprintf(os.Stderr, "❌ A process failed: %v\nStopping all processes...\n", err)
		cancel() // Stop all processes
		<-done   // Wait for cleanup
		fmt.Println("✅ All processes stopped")
		os.Exit(1)
	case <-done:
		fmt.Println("✅ All processes completed")
	}
}
