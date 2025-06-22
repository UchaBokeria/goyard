package main

import (
	"log"

	"github.com/UchaBokeria/goyard/examples/simple/web/api"
	"github.com/UchaBokeria/goyard/examples/simple/web/app"
	"github.com/UchaBokeria/goyard/examples/simple/web/middlewares"
	"github.com/UchaBokeria/goyard/goyard"
)

func main() {
	web := goyard.New()
	web.Use(middlewares.Htmx())
	web.Use(middlewares.Interceptor())
	app.New(web.Group(""))
	api.New(web.Group("/api"))
	log.Fatal(web.Run(":3000"))
}
