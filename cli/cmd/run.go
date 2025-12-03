package cmd

import (
	"fmt"

	"github.com/UchaBokeria/goyard/cli/cmd/run"
	"github.com/spf13/cobra"
)

var (
	host       string
	port       int
	production bool
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the GoYard project with live reload",
	Long:  `Run the GoYard application with hot reloading in development mode`,
	Run: func(cmd *cobra.Command, args []string) {
		if production {
			runProductionServer()
		} else {
			runDevelopmentServer()
		}
	},
}

func init() {
	runCmd.Flags().StringVarP(&host, "host", "H", "localhost", "Host to run the server on")
	runCmd.Flags().IntVarP(&port, "port", "p", 3000, "Port to run the server on")
	runCmd.Flags().BoolVarP(&production, "production", "P", false, "Run in production mode")
}

func runDevelopmentServer() {
	fmt.Printf("🔥 Development is starting on %s:%d\n", host, port)
	run.Host = host
	run.Port = port
	run.Dev()
}

func runProductionServer() {
}
