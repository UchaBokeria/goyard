package cmd

import (
	"fmt"

	"github.com/UchaBokeria/goyard/goyard"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "versionZ",
	Short: "Display the version of GoYard",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("GoYard v%s\n", goyard.Version)
	},
}
