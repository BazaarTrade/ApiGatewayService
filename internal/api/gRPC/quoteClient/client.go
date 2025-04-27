package qClient

import (
	"context"
	"log/slog"

	ws "github.com/BazaarTrade/ApiGatewayService/internal/api/websocket"
	"github.com/BazaarTrade/QuoteProtoGen/pbQ"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client       pbQ.QuoteClient
	conn         *grpc.ClientConn
	hub          *ws.Hub
	ctx          context.Context
	cancel       context.CancelFunc
	cancelByPair map[string]context.CancelFunc
	logger       *slog.Logger
}

func New(hub *ws.Hub, logger *slog.Logger) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		hub:          hub,
		ctx:          ctx,
		cancel:       cancel,
		cancelByPair: make(map[string]context.CancelFunc),
		logger:       logger,
	}
}

func (c *Client) Run(CONN_ADDR_QUOTE string) error {
	var err error
	c.conn, err = grpc.NewClient(CONN_ADDR_QUOTE, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.logger.Error("failed connecting to localhost:50052", "error", err)
		return err
	}

	c.client = pbQ.NewQuoteClient(c.conn)

	return nil
}

func (c *Client) CloseConnection() {
	if c.cancel != nil {
		c.cancel()
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			c.logger.Error("failed to close qClient connection", "error", err)
		}
	}
}

func (c *Client) StopStreamReadersByPair(pair string) {
	if cancel, ok := c.cancelByPair[pair]; ok {
		cancel()
		delete(c.cancelByPair, pair)
	}
}
