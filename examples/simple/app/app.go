package app

import (
	"github.com/UchaBokeria/goyard/controller"
	"github.com/UchaBokeria/goyard/examples/simple/app/dto"
	"github.com/labstack/echo/v4"
)

func New(router *echo.Group) {
	router.GET("", controller.Set[any](index))
	router.POST("user", controller.Set[dto.UserCreateDto](createUser))
}
