// Package goyard provides a comprehensive web framework toolkit for Go.
package goyard

import (
	"github.com/labstack/echo/v4"
	"github.com/yourorg/goyard/controller"
)

// Version represents the current version of the goyard framework
const Version = "0.1.0"

// New returns the goyard middleware which wraps echo.Context with additional
// functionality (see controller.Initialize). It enables:
//   - extended context helpers (HTML rendering, HTMX checks, etc.)
//   - automatic validation & DTO binding helpers via controller.Use / controller.Set
//
// Usage:
//
//	e := echo.New()
//	e.Use(goyard.New())
func New() echo.MiddlewareFunc {
	return controller.Initialize()
}
