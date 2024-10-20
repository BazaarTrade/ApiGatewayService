package app

import (
	"log/slog"
	"os"

	clientGRPC "github.com/BazaarTrade/ApiGatewayService/internal/api/gRPC"
	"github.com/BazaarTrade/ApiGatewayService/internal/api/rest"
	ws "github.com/BazaarTrade/ApiGatewayService/internal/api/websocket"
)

func Run() {
	handler := slog.NewTextHandler(os.Stdout, nil)
	logger := slog.New(handler)

	logger.Info("starting aplication...")

	hub := ws.NewHub(logger)

	clientGRPC := clientGRPC.New(hub, logger)
	clientGRPC.RunPBMClient()
	clientGRPC.RunPBQClient()

	go clientGRPC.ReadPrecisedOrderBookSnapshot()
	go clientGRPC.ReadPrecisedTrades()

	defer clientGRPC.CloseConnections()

	server := rest.New(hub, clientGRPC, logger)
	server.Run()
}
