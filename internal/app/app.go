package app

import (
	"log/slog"
	"os"

	mClient "github.com/BazaarTrade/ApiGatewayService/internal/api/gRPC/matchingEngineClient"
	qClient "github.com/BazaarTrade/ApiGatewayService/internal/api/gRPC/quoteClient"
	"github.com/BazaarTrade/ApiGatewayService/internal/api/rest"
	ws "github.com/BazaarTrade/ApiGatewayService/internal/api/websocket"
	"github.com/joho/godotenv"
)

func Run() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			logger.Error("failed to load .env file")
		}
	}

	logger.Info("starting aplication...")

	hub := ws.NewHub(logger)

	mClient := mClient.New(hub, logger)
	if err := mClient.Run(); err != nil {
		return
	}
	defer mClient.CloseConnection()

	qClient := qClient.New(hub, logger)
	if err := qClient.Run(); err != nil {
		return
	}
	defer qClient.CloseConnection()

	if err := initOrderBooks(mClient, qClient, hub); err != nil {
		return
	}

	rest := rest.New(mClient, qClient, hub, logger)
	if err := rest.Run(); err != nil {
		return
	}
}

func initOrderBooks(mClient *mClient.Client, qClient *qClient.Client, hub *ws.Hub) error {
	pairsParams, err := mClient.GetPairsParams()
	if err != nil {
		return err
	}

	for _, pairParams := range pairsParams {
		qClient.StartStreamReaders(pairParams.Pair)
		hub.AddOrderBookUpdateTopic(pairParams.Pair, pairParams.PricePrecisions)
		hub.AddTradesTopic(pairParams.Pair)
		hub.AddTickerTopic(pairParams.Pair)
	}
	return nil
}
