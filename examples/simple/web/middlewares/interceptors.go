package middlewares

import (
	"github.com/UchaBokeria/goyard/controller"
	"github.com/labstack/echo/v4"
)

func Interceptor() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return controller.Set[any](func(ctx *controller.Context[any]) error {
			return next(ctx)
		})
	}
}
