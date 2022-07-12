package http

import (
	"context"

	"github.com/Abdulsametileri/package-tracking-app/domain"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type PackageHandler struct {
	upgrader websocket.Upgrader
	PUsecase domain.PackageUsecase
}

func NewPackageHandler(e *echo.Echo, pu domain.PackageUsecase) {
	handler := &PackageHandler{
		upgrader: websocket.Upgrader{},
		PUsecase: pu,
	}

	e.GET("/packages/track/:vehicleId", handler.TrackByVehicleID)
}

func (p *PackageHandler) TrackByVehicleID(c echo.Context) error {
	wsConn, err := p.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	go func() {
		_, _, err = wsConn.ReadMessage()
		if err != nil {
			cancelFunc()
		}
	}()

	for {
		select {
		case <-ctx.Done():
			wsConn.Close()
			return nil
		default:
			p, err := p.PUsecase.TrackByVehicleID(ctx, c.Param("vehicleId"))
			if err != nil {
				c.Logger().Error(err)
				continue
			}

			err = wsConn.WriteJSON(p)
			if err != nil {
				c.Logger().Error(err)
			}
		}
	}
}
