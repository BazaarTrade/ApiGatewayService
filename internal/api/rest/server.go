package rest

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	aClient "github.com/BazaarTrade/ApiGatewayService/internal/api/gRPC/authClient"
	mClient "github.com/BazaarTrade/ApiGatewayService/internal/api/gRPC/matchingEngineClient"
	qClient "github.com/BazaarTrade/ApiGatewayService/internal/api/gRPC/quoteClient"
	ws "github.com/BazaarTrade/ApiGatewayService/internal/api/websocket"
	"github.com/BazaarTrade/ApiGatewayService/internal/middleware"
	"github.com/BazaarTrade/ApiGatewayService/internal/repository"
	"github.com/labstack/echo/v4"
	eMiddleware "github.com/labstack/echo/v4/middleware"
)

type Server struct {
	mClient *mClient.Client
	qClient *qClient.Client
	aClient *aClient.Client
	echo    *echo.Echo

	hub                     *ws.Hub
	db                      repository.Repository
	accessTokenSecretPhrase string
	logger                  *slog.Logger
}

func New(mClient *mClient.Client, qClient *qClient.Client, aClient *aClient.Client, hub *ws.Hub, db repository.Repository, ACCESS_TOKEN_SECRET_PHRASE string, logger *slog.Logger) *Server {
	return &Server{
		mClient:                 mClient,
		qClient:                 qClient,
		aClient:                 aClient,
		hub:                     hub,
		db:                      db,
		accessTokenSecretPhrase: ACCESS_TOKEN_SECRET_PHRASE,
		logger:                  logger,
	}
}

func (s *Server) Run(ADDR string) error {
	e := echo.New()
	s.echo = e

	CORS(e)
	s.init(e)

	s.logger.Info("server is listeninig on port" + ADDR)

	if err := e.Start(ADDR); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("failed to serve", "err", err)
		return err
	}
	return nil
}

func (s *Server) init(e *echo.Echo) {
	e.Use(eMiddleware.Recover())

	//public routes
	e.POST("/register", s.register)
	e.POST("/login", s.login)
	e.POST("/refresh/:userID", s.refreshAccessToken)

	e.POST("/orderbook", s.createOrderBook)
	e.DELETE("/orderbook/:pair", s.deleteOrderBook)
	e.GET("/ws/:userID", s.hub.HandleWebsocket)
	e.GET("/orderBookPricePrecisions/:pair", s.getOrderBookPricePrecisions)
	e.GET("/candleSticks", s.getCandleStickHistory)

	//private routes
	g := e.Group("")
	g.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return middleware.AccessTokenMiddleware(next, s.accessTokenSecretPhrase)
	})

	g.POST("/logout", s.logout)
	g.POST("/changePassword", s.changePassword)

	g.POST("/order", s.placeOrder)
	g.DELETE("/order/:orderID", s.cancelOrder)
	g.GET("/orders/current/:userID", s.getCurrentOrders)
	g.GET("/orders/:userID", s.getOrders)
}

func CORS(e *echo.Echo) {
	e.Use(eMiddleware.CORSWithConfig(eMiddleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
}

func (s *Server) Stop() error {
	if err := s.echo.Shutdown(context.Background()); err != nil {
		s.logger.Error("failed to shutdown server", "error", err)
		return err
	}

	return nil
}
