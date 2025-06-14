package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	port       int
	production bool
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the GoYard application",
	Long:  `Run the GoYard application with hot reloading in development mode`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Starting GoYard application...")

		if production {
			runProductionServer()
		} else {
			runDevelopmentServer()
		}
	},
}

func init() {
	runCmd.Flags().IntVarP(&port, "port", "p", 3000, "Port to run the server on")
	runCmd.Flags().BoolVarP(&production, "production", "P", false, "Run in production mode")
}

func runDevelopmentServer() {
	fmt.Printf("Running in development mode on port %d\n", port)
	fmt.Println("Hot reloading is enabled")

	// Check if air is installed for hot reloading
	airPath, err := exec.LookPath("air")
	if err != nil {
		fmt.Println("Air not found. Installing air for hot reloading...")
		installAir()
		airPath, _ = exec.LookPath("air")
	}

	// Create .air.toml if it doesn't exist
	if _, err := os.Stat(".air.toml"); os.IsNotExist(err) {
		createAirConfig()
	}

	// Run with air for hot reloading
	cmd := exec.Command(airPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), fmt.Sprintf("PORT=%d", port))

	fmt.Println("Starting development server...")
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run development server: %v\n", err)
		os.Exit(1)
	}
}

func runProductionServer() {
	fmt.Printf("Running in production mode on port %d\n", port)

	// Build and run the application
	buildCmd := exec.Command("go", "build", "-o", "app", ".")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	fmt.Println("Building application...")
	if err := buildCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to build application: %v\n", err)
		os.Exit(1)
	}

	// Run the compiled application
	appPath, _ := filepath.Abs("./app")
	runCmd := exec.Command(appPath)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	runCmd.Env = append(os.Environ(), fmt.Sprintf("PORT=%d", port), "GO_ENV=production")

	fmt.Println("Starting production server...")
	if err := runCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run production server: %v\n", err)
		os.Exit(1)
	}
}

func installAir() {
	installCmd := exec.Command("go", "install", "github.com/cosmtrek/air@latest")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr

	if err := installCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to install air: %v\n", err)
		fmt.Println("You can install it manually with: go install github.com/cosmtrek/air@latest")
		os.Exit(1)
	}
}

func createAirConfig() {
	airConfig := `root = "."
tmp_dir = "tmp"
[build]
  cmd = "go build -o ./tmp/main ."
  bin = "./tmp/main"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor"]
  include_ext = ["go", "tpl", "tmpl", "templ", "html"]
  exclude_regex = ["_test\\.go"]
[screen]
  clear_on_rebuild = true
`

	if err := os.WriteFile(".air.toml", []byte(airConfig), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create .air.toml: %v\n", err)
	}
}
