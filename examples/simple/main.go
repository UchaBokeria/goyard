package main

import (
	"log"

	"github.com/UchaBokeria/goyard"
	"github.com/UchaBokeria/goyard/examples/simple/api"
	"github.com/UchaBokeria/goyard/examples/simple/app"
)

func main() {
	web := goyard.New()
	app.New(web.Group(""))
	api.New(web.Group("/api"))
	log.Fatal(web.Start(":8080"))
}
