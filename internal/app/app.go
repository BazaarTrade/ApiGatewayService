package app

import (
	"log/slog"
	"os"

	mClient "github.com/BazaarTrade/ApiGatewayService/internal/api/gRPC/matchingEngineClient"
	qClient "github.com/BazaarTrade/ApiGatewayService/internal/api/gRPC/quoteClient"
	"github.com/BazaarTrade/ApiGatewayService/internal/api/rest"
	ws "github.com/BazaarTrade/ApiGatewayService/internal/api/websocket"
	"github.com/BazaarTrade/ApiGatewayService/internal/repository"
	"github.com/BazaarTrade/ApiGatewayService/internal/repository/postgresPgx"
	"github.com/joho/godotenv"
)

func Run() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	if _, err := os.Stat("../.env"); err == nil {
		if err := godotenv.Load("../.env"); err != nil {
			logger.Error("failed to load .env file", "error", err)
			return
		}
	}

	logger.Info("starting aplication...")

	repository, err := postgresPgx.NewPostgres(logger)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		return
	}

	hub := ws.NewHub(logger)

	mClient := mClient.New(logger)
	if err := mClient.Run(); err != nil {
		return
	}
	defer mClient.CloseConnection()

	qClient := qClient.New(hub, logger)
	if err := qClient.Run(); err != nil {
		return
	}
	defer qClient.CloseConnection()

	if err := initPairs(qClient, hub, repository); err != nil {
		return
	}

	rest := rest.New(mClient, qClient, hub, repository, logger)
	if err := rest.Run(); err != nil {
		return
	}
}

func initPairs(qClient *qClient.Client, hub *ws.Hub, repository repository.Repository) error {
	pairsParams, err := repository.GetPairsParams()
	if err != nil {
		return nil
	}

	for _, pairParams := range pairsParams {
		hub.AddOrderBookSnapshotTopic(pairParams.Pair, pairParams.OrderBookPricePrecisions)
		hub.AddTradesTopic(pairParams.Pair)
		hub.AddTickerTopic(pairParams.Pair)
		hub.AddCandleStickTopic(pairParams.Pair, pairParams.CandleStickTimeframes)
		qClient.StartStreamReaders(pairParams.Pair)
	}
	return nil
}
