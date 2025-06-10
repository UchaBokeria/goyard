package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yourorg/goyard"
)

func main() {
	fmt.Printf("GoYard Version: %s\n", goyard.Version)

	e := echo.New()

	// Attach GoYard middleware (adds extended context capabilities)
	e.Use(goyard.New())

	// Define a simple route for demonstration
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from GoYard example!")
	})

	log.Fatal(e.Start(":8080"))
}
