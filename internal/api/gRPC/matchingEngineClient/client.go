package mClient

import (
	"log/slog"
	"os"

	"github.com/BazaarTrade/MatchingEngineProtoGen/pbM"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client pbM.MatchingEngineClient
	conn   *grpc.ClientConn
	logger *slog.Logger
}

func New(logger *slog.Logger) *Client {
	return &Client{
		logger: logger,
	}
}

func (c *Client) Run() error {
	var err error
	c.conn, err = grpc.NewClient(os.Getenv("CONN_ADDR_MATCHING_ENGINE"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.logger.Error("failed connecting to localhost:50051", "error", err)
		return err
	}

	c.client = pbM.NewMatchingEngineClient(c.conn)

	return nil
}

func (c *Client) CloseConnection() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			c.logger.Error("failed to close Matching Engine connection", "error", err)
		}
	}
}
