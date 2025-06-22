package middlewares

import "github.com/labstack/echo/v4"

func Htmx() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Header.Get("HX-Request") == "true" {
				return next(c)
			}
			return next(c)
		}
	}
}
