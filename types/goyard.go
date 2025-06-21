package types

import "github.com/labstack/echo/v4"

type Route struct {
	echo.Group
}

type Goyard struct {
	echo.Echo
	Run func(address string) error
}
