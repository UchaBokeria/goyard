package middlewares

import (
	"github.com/UchaBokeria/goyard/controller"
	"github.com/labstack/echo/v4"
)

func Interceptor() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return controller.Use(func(ctx *controller.Context) error {
			return next(ctx)
		})
	}
}
