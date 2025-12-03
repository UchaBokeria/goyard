package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goyard",
	Short: "GoYard - A comprehensive Go web framework toolkit",
	Long: `* GoYard is a web framework toolkit that provides essential components
for web development including HTMX integration, Alpine.js, Tailwind CSS,
and hot reloading for development.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommands are provided, print the help message.
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.

	// Add subcommands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(runCmd)
}
