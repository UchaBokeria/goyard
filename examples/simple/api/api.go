package api

import (
	"github.com/UchaBokeria/goyard/controller"
	"github.com/UchaBokeria/goyard/examples/simple/view"
	"github.com/labstack/echo/v4"
)

func index(ctx *controller.Context[any]) error {
	return ctx.Html(view.Index())
}

func New(router *echo.Group) {
	router.GET("/", controller.Set[any](index))
}
