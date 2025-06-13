package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	template string
)

// newCmd represents the new command to create a new GoYard project
var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new GoYard project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		projectPath := filepath.Join(".", projectName)

		fmt.Printf("Creating new GoYard project: %s\n", projectName)

		// Create project directory
		if err := os.MkdirAll(projectPath, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create project directory: %v\n", err)
			os.Exit(1)
		}

		// Create project structure based on template
		createProjectStructure(projectPath, template)

		fmt.Printf("Project created successfully at %s\n", projectPath)
		fmt.Println("Run the following commands to get started:")
		fmt.Printf("  cd %s\n", projectName)
		fmt.Println("  go mod tidy")
		fmt.Println("  go run main.go")
	},
}

func init() {
	newCmd.Flags().StringVarP(&template, "template", "t", "basic", "Template to use (basic, api, htmx, full)")
}

func createProjectStructure(projectPath, template string) {
	// Create main.go
	mainContent := `package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/UchaBokeria/goyard"
)

func main() {
	fmt.Println("Starting GoYard application...")
	
	// Initialize GoYard
	app := goyard.New()
	
	// Start the server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
`

	// Create go.mod
	goModContent := fmt.Sprintf(`module %s

go 1.21

require (
	"github.com/UchaBokeria/goyard v0.1.0
)
`, filepath.Base(projectPath))

	// Write files
	writeFile(filepath.Join(projectPath, "main.go"), mainContent)
	writeFile(filepath.Join(projectPath, "go.mod"), goModContent)

	// Create other directories based on the template
	createDirs(projectPath, []string{
		"internal",
		"pkg",
		"web/templates",
		"web/static/css",
		"web/static/js",
	})
}

func createDirs(basePath string, dirs []string) {
	for _, dir := range dirs {
		path := filepath.Join(basePath, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create directory %s: %v\n", path, err)
		}
	}
}

func writeFile(path, content string) {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write file %s: %v\n", path, err)
	}
}
