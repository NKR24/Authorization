package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	fmt.Println("http://localhost:5000")
	e := echo.New()
	e.Use(middleware.CORS())
	e.Logger.Fatal(e.Start(":5000"))
}
