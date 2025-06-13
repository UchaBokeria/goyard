package app

import (
	"github.com/UchaBokeria/goyard/controller"
	"github.com/UchaBokeria/goyard/examples/simple/app/dto"
	"github.com/UchaBokeria/goyard/examples/simple/view/components"
	view "github.com/UchaBokeria/goyard/examples/simple/view/pages"
)

func index(ctx *controller.Context[any]) error {
	return ctx.Html(view.Page())
}

func createUser(ctx *controller.Context[any], user *dto.UserCreateDto) error {
	return ctx.Html(components.User(user))
}
