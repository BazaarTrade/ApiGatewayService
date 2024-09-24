package rest

import (
	ws "github.com/BazaarTrade/ApiGatewayService/internal/api/websocket"
	"github.com/BazaarTrade/GeneratedProto/pb"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handler struct {
	pbClient pb.MatchingEngineClient
	hub      *ws.Hub
}

func NewHandler(hub *ws.Hub, pbClient pb.MatchingEngineClient) *Handler {
	return &Handler{
		hub:      hub,
		pbClient: pbClient,
	}
}

func (h *Handler) Init(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/order", h.placeOrder)
	e.DELETE("/order/:orderID", h.cancelOrder)
	e.GET("/orders/current/:userID", h.getCurrentOrders)
	e.GET("/orders/:userID", h.getOrders)
	e.POST("/orderbook/:symbol", h.createOrderBook)
	e.DELETE("/orderbook/:symbol", h.deleteOrderBook)
	e.GET("/ws/:userID", h.hub.HandleWebsocket)
}
