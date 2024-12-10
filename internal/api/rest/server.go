package rest

import (
	"log/slog"

	mClient "github.com/BazaarTrade/ApiGatewayService/internal/api/gRPC/matchingEngineClient"
	qClient "github.com/BazaarTrade/ApiGatewayService/internal/api/gRPC/quoteClient"
	ws "github.com/BazaarTrade/ApiGatewayService/internal/api/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	mClient *mClient.Client
	qClient *qClient.Client

	hub    *ws.Hub
	logger *slog.Logger
}

func New(mClient *mClient.Client, qClient *qClient.Client, hub *ws.Hub, logger *slog.Logger) *Server {
	return &Server{
		mClient: mClient,
		qClient: qClient,
		hub:     hub,
		logger:  logger,
	}
}

func (s *Server) Run() error {
	e := echo.New()
	CORS(e)
	s.init(e)

	s.logger.Info("server is listeninig on port 8080...")

	if err := e.Start(":8080"); err != nil {
		s.logger.Error("failed to serve", "err", err)
		return err
	}
	return nil
}

func (s *Server) init(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/order", s.placeOrder)
	e.DELETE("/order/:orderID", s.cancelOrder)
	e.GET("/orders/current/:userID", s.getCurrentOrders)
	e.GET("/orders/:userID", s.getOrders)
	e.POST("/orderbook", s.createOrderBook)
	e.DELETE("/orderbook/:pair", s.deleteOrderBook)
	e.GET("/ws/:userID", s.hub.HandleWebsocket)
	e.GET("/pricePrecisions/:pair", s.getPairPricePrecisions)
}

func CORS(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
}
