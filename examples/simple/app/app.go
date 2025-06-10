package app

import "github.com/labstack/echo/v4"

func app(app *echo.Group) {
	app.GET("", controller.Set[any](index))
}
