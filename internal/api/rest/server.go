package rest

import (
	"log/slog"

	clientGRPC "github.com/BazaarTrade/ApiGatewayService/internal/api/gRPC"
	ws "github.com/BazaarTrade/ApiGatewayService/internal/api/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	hub        *ws.Hub
	clientGRPC *clientGRPC.GRPCClients
	logger     *slog.Logger
}

func New(hub *ws.Hub, clientGRPC *clientGRPC.GRPCClients, logger *slog.Logger) *Server {
	return &Server{
		hub:        hub,
		clientGRPC: clientGRPC,
		logger:     logger,
	}
}

func (s *Server) Run() {
	e := echo.New()
	CORS(e)
	s.init(e)
	e.Start(":8080")
}

func (s *Server) init(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/order", s.placeOrder)
	e.DELETE("/order/:orderID", s.cancelOrder)
	e.GET("/orders/current/:userID", s.getCurrentOrders)
	e.GET("/orders/:userID", s.getOrders)
	e.POST("/orderbook/:symbol", s.createOrderBook)
	e.DELETE("/orderbook/:symbol", s.deleteOrderBook)
	e.GET("/ws/:userID", s.hub.HandleWebsocket)
}

func CORS(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
}
