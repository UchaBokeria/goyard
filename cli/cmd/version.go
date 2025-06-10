package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourorg/goyard"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version of GoYard",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("GoYard v%s\n", goyard.Version)
	},
}
