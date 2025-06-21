package goyard

import (
	"strings"

	"github.com/UchaBokeria/goyard/controller"
	"github.com/UchaBokeria/goyard/types"
	"github.com/labstack/echo/v4"
)

// Version represents the current version of the goyard framework
const Version = "0.1.0"

func RunRegister(echo *echo.Echo) func(address string) error {
	return func(address string) error {
		parts := strings.Split(address, ":")
		if len(parts) > 1 {
			address = strings.Join(parts[:len(parts)-1], ":") + ":7331"
		}
		return echo.Start(address)
	}
}

// New returns the goyard middleware which wraps echo.Context with additional
// functionality (see controller.Initialize). It enables:
//   - extended context helpers (HTML rendering, HTMX checks, etc.)
//   - automatic validation & DTO binding helpers via controller.Use / controller.Set
//
// Usage:
//
//	e := echo.New()
//	e.Use(goyard.New())
func New() *types.Goyard {
	echo := echo.New()
	echo.Use(controller.Initialize())
	return &types.Goyard{Echo: *echo, Run: RunRegister(echo)}
}

func Use(echo *types.Goyard) {
	echo.Use(controller.Initialize())
}
