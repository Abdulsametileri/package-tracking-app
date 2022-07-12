package main

import (
	"log"
	"time"

	"github.com/Abdulsametileri/package-tracking-app/domain"
	_packageClient "github.com/Abdulsametileri/package-tracking-app/package/client"
	_packageHttpDelivery "github.com/Abdulsametileri/package-tracking-app/package/delivery/http"
	_packageUcase "github.com/Abdulsametileri/package-tracking-app/package/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.File("/", "website/index.html")

	pc, err := _packageClient.NewRabbitMQClient("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Panicln(err)
	}
	defer pc.Close()

	go func() {
		time.Sleep(3 * time.Second)
		pc.Publish(domain.Package{From: "Bordeaux", To: "Toulouse", VehicleID: "123"})
		time.Sleep(3 * time.Second)
		pc.Publish(domain.Package{From: "Toulouse", To: "Monaco", VehicleID: "123"})
		time.Sleep(3 * time.Second)
		pc.Publish(domain.Package{From: "Monaco", To: "Lyon", VehicleID: "123"})
		time.Sleep(3 * time.Second)
		pc.Publish(domain.Package{From: "Lyon", To: "Paris", VehicleID: "123"})
		time.Sleep(3 * time.Second)
		pc.Publish(domain.Package{From: "Paris", To: "Brussels", VehicleID: "123"})
		time.Sleep(3 * time.Second)
		pc.Publish(domain.Package{From: "Brussels", To: "Rotterdam", VehicleID: "123"})
		time.Sleep(3 * time.Second)
		pc.Publish(domain.Package{From: "Rotterdam", To: "Amsterdam", VehicleID: "123"})
	}()

	pu := _packageUcase.NewPackageUsecase(pc)

	_packageHttpDelivery.NewPackageHandler(e, pu)

	e.Logger.Fatal(e.Start(":1323"))
}
