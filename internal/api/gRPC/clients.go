package clientGRPC

import (
	"log/slog"
	"os"

	ws "github.com/BazaarTrade/ApiGatewayService/internal/api/websocket"
	"github.com/BazaarTrade/MatchingEngineProtoGen/pbM"
	"github.com/BazaarTrade/QuoteProtoGen/pbQ"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClients struct {
	pbMClient pbM.MatchingEngineClient
	pbQClient pbQ.QuoteClient
	pbMConn   *grpc.ClientConn
	pbQConn   *grpc.ClientConn
	logger    *slog.Logger
	hub       *ws.Hub
}

func New(hub *ws.Hub, logger *slog.Logger) *GRPCClients {
	return &GRPCClients{
		logger: logger,
		hub:    hub,
	}
}

func (c *GRPCClients) RunPBMClient() error {
	var err error
	c.pbMConn, err = grpc.NewClient(os.Getenv("CONN_ADDR_MATCHING_ENGINE"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.logger.Error("failed connecting to localhost:50051", "error", err)
		return err
	}

	c.pbMClient = pbM.NewMatchingEngineClient(c.pbMConn)

	return nil
}

func (c *GRPCClients) RunPBQClient() error {
	var err error
	c.pbQConn, err = grpc.NewClient(os.Getenv("CONN_ADDR_QUOTE"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.logger.Error("failed connecting to localhost:50052", "error", err)
		return err
	}

	c.pbQClient = pbQ.NewQuoteClient(c.pbQConn)

	return nil
}

func (c *GRPCClients) CloseConnections() {
	if c.pbMConn != nil {
		if err := c.pbMConn.Close(); err != nil {
			c.logger.Error("failed to close Matching Engine connection", "error", err)
		}
	}
	if c.pbQConn != nil {
		if err := c.pbQConn.Close(); err != nil {
			c.logger.Error("failed to close Quote Service connection", "error", err)
		}
	}
}
