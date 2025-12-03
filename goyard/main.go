package goyard

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/UchaBokeria/goyard/controller"
	"github.com/UchaBokeria/goyard/types"
	"github.com/labstack/echo/v4"
)

// Version represents the current version of the goyard framework
const Version = "0.1.3"

func RunRegister(echo *echo.Echo) func(address string) error {
	return func(address string) error {
		parts := strings.Split(address, ":")
		originalPort, err := strconv.Atoi(parts[len(parts)-1])
		newPort := originalPort - 1
		if err != nil {
			return fmt.Errorf("invalid port: %w", err)
		}
		if len(parts) > 1 {
			address = strings.Join(parts[:len(parts)-1], ":") + ":" + strconv.Itoa(newPort)
		} else {
			address = "localhost:" + strconv.Itoa(newPort)
		}
		fmt.Println("🚀 Starting server on", strings.Replace(address, strconv.Itoa(newPort), strconv.Itoa(originalPort), 1))
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
//	app := goyard.New()
//	app.Use(middlewares.Htmx())
//	app.Use(middlewares.Interceptor())
//	app.Run(":3000")
func New() *types.Goyard {
	echo := echo.New()
	echo.HidePort = true
	echo.HideBanner = true
	echo.Use(controller.Initialize())
	return &types.Goyard{Echo: echo, Run: RunRegister(echo)}
}

