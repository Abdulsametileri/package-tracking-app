package http

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type PackageHandler struct {
	upgrader websocket.Upgrader
}

func NewPackageHandler(e *echo.Echo) {
	handler := &PackageHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	e.GET("/packages/track/:vehicleId", handler.TrackByVehicleID)
}

func (p *PackageHandler) TrackByVehicleID(c echo.Context) error {
	vehicleID := c.Param("vehicleId")
	_ = vehicleID

	wsConn, err := p.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer wsConn.Close()

	for {
		err = wsConn.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			c.Logger().Error(err)
		}

		_, msg, err := wsConn.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)
	}
}
