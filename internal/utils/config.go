// Package utils provides internal utilities for the goyard framework
package utils

// Config represents application configuration
type Config struct {
	AppName  string
	Port     int
	LogLevel string
	Debug    bool
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		AppName:  "Goyard App",
		Port:     8080,
		LogLevel: "info",
		Debug:    false,
	}
}
