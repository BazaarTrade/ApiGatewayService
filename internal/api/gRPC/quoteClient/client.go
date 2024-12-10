package qClient

import (
	"log/slog"
	"os"

	ws "github.com/BazaarTrade/ApiGatewayService/internal/api/websocket"
	"github.com/BazaarTrade/QuoteProtoGen/pbQ"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client pbQ.QuoteClient
	conn   *grpc.ClientConn
	hub    *ws.Hub
	logger *slog.Logger
}

func New(hub *ws.Hub, logger *slog.Logger) *Client {
	return &Client{
		hub:    hub,
		logger: logger,
	}
}

func (c *Client) Run() error {
	var err error
	c.conn, err = grpc.NewClient(os.Getenv("CONN_ADDR_QUOTE"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.logger.Error("failed connecting to localhost:50052", "error", err)
		return err
	}

	c.client = pbQ.NewQuoteClient(c.conn)

	return nil
}

func (c *Client) CloseConnection() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			c.logger.Error("failed to close Quote Service connection", "error", err)
		}
	}
}
