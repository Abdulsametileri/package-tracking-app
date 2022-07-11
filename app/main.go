package main

import (
	"fmt"

	_packageHttpDelivery "github.com/Abdulsametileri/package-tracking-app/package/delivery/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	r := e.File("/", "website/index.html")
	fmt.Println(r)

	_packageHttpDelivery.NewPackageHandler(e)

	e.Logger.Fatal(e.Start(":1323"))
}
