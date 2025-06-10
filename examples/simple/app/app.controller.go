package app

import (
	"github.com/yourorg/goyard/context"
)

func index(ctx *context.Context) error {
	return ctx.Html()
}
