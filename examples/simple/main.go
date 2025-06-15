package main

import (
	"log"

	"github.com/UchaBokeria/goyard/examples/simple/api"
	"github.com/UchaBokeria/goyard/examples/simple/app"
	"github.com/UchaBokeria/goyard/goyard"
)

func main() {
	web := goyard.New()
	app.New(web.Group(""))
	api.New(web.Group("/api"))
	log.Fatal(web.Start(":8080"))
}
