package aClient

import (
	"log/slog"

	"github.com/BazaarTrade/AuthProtoGen/pbA"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client pbA.AuthClient
	conn   *grpc.ClientConn
	logger *slog.Logger
}

func New(logger *slog.Logger) *Client {
	return &Client{
		logger: logger,
	}
}

func (c *Client) Run(CONN_ADDR_AUTH string) error {
	var err error
	c.conn, err = grpc.NewClient(CONN_ADDR_AUTH, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.logger.Error("failed connecting to localhost:50052", "error", err)
		return err
	}

	c.client = pbA.NewAuthClient(c.conn)

	return nil
}

func (c *Client) CloseConnection() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			c.logger.Error("failed to close aClient connection", "error", err)
		}
	}
}
