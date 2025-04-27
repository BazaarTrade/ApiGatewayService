package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	aClient "github.com/BazaarTrade/ApiGatewayService/internal/api/gRPC/authClient"
	mClient "github.com/BazaarTrade/ApiGatewayService/internal/api/gRPC/matchingEngineClient"
	qClient "github.com/BazaarTrade/ApiGatewayService/internal/api/gRPC/quoteClient"
	"github.com/BazaarTrade/ApiGatewayService/internal/api/rest"
	ws "github.com/BazaarTrade/ApiGatewayService/internal/api/websocket"
	"github.com/BazaarTrade/ApiGatewayService/internal/repository"
	"github.com/BazaarTrade/ApiGatewayService/internal/repository/postgresPgx"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	if _, err := os.Stat("../.env"); err == nil {
		if err := godotenv.Load("../.env"); err != nil {
			logger.Error("failed to load .env file", "error", err)
			return
		}
	}

	logger.Info("starting aplication...")

	DB_CONNECTION := os.Getenv("DB_CONNECTION")
	if DB_CONNECTION == "" {
		logger.Error("DB_CONNECTION environment variable is not set")
		return
	}

	repository, err := postgresPgx.NewPostgres(DB_CONNECTION, logger)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		return
	}

	hub := ws.NewHub(logger)

	CONN_ADDR_MATCHING_ENGINE := os.Getenv("CONN_ADDR_MATCHING_ENGINE")
	if CONN_ADDR_MATCHING_ENGINE == "" {
		logger.Error("CONN_ADDR_MATCHING_ENGINE environment variable is not set")
		return
	}

	mClient := mClient.New(logger)
	if err := mClient.Run(CONN_ADDR_MATCHING_ENGINE); err != nil {
		return
	}

	CON_ADDR_QUOTE := os.Getenv("CONN_ADDR_QUOTE")
	if CON_ADDR_QUOTE == "" {
		logger.Error("CONN_ADDR_QUOTE environment variable is not set")
		return
	}

	qClient := qClient.New(hub, logger)
	if err := qClient.Run(CON_ADDR_QUOTE); err != nil {
		return
	}

	CONN_ADDR_AUTH := os.Getenv("CONN_ADDR_AUTH")
	if CONN_ADDR_AUTH == "" {
		logger.Error("CONN_ADDR_AUTH environment variable is not set")
		return
	}

	aClient := aClient.New(logger)
	if err := aClient.Run(CONN_ADDR_AUTH); err != nil {
		return
	}

	if err := initPairs(qClient, hub, repository); err != nil {
		return
	}

	ACCESS_TOKEN_SECRET_PHRASE := os.Getenv("ACCESS_TOKEN_SECRET_PHRASE")
	if ACCESS_TOKEN_SECRET_PHRASE == "" {
		logger.Error("ACCESS_TOKEN_SECRET_PHRASE environment variable is not set")
		return
	}

	ADDR := os.Getenv("ADDR")
	if ADDR == "" {
		logger.Error("ADDR environment variable is not set")
		return
	}

	rest := rest.New(mClient, qClient, aClient, hub, repository, ACCESS_TOKEN_SECRET_PHRASE, logger)
	go func() {
		if err := rest.Run(ADDR); err != nil {
			os.Exit(1)
		}
	}()

	//Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	logger.Info("shutting down...")

	qClient.CloseConnection()
	logger.Info("closed qClient connection")

	mClient.CloseConnection()
	logger.Info("closed mClient connection")

	aClient.CloseConnection()
	logger.Info("closed aClient connection")

	rest.Stop()
	logger.Info("stopped gRPC server")

	repository.Close()
	logger.Info("closed database connection")

	logger.Info("gracefully stopped")
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
