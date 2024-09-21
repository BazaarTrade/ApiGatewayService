package api

import (
	"github.com/BazaarTrade/GeneratedProto/pb"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handler struct {
	pbClient pb.MatchingEngineClient
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Init(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/order", h.placeOrder)
	e.DELETE("/order/:order_id", h.cancelOrder)
	e.GET("/orders/current/:user_id", h.getCurrentOrders)
	e.GET("/orders/:user_id", h.getOrders)
	e.POST("/orderbook/:symbol", h.createOrderBook)
	e.DELETE("/orderbook/:symbol", h.deleteOrderBook)
}
