package main

import (
	"thirdopinion/internal/pkg/handlers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/ws", handlers.WS)

	// Start server
	e.Logger.Fatal(e.Start(":8081"))
}
