package qClient

import (
	"context"

	"github.com/BazaarTrade/ApiGatewayService/internal/models"
	"github.com/BazaarTrade/QuoteProtoGen/pbQ"
)

func (c *Client) CreateOrderBook(pairParams models.PairParams) error {
	_, err := c.client.CreateOrderBook(context.Background(), &pbQ.PairParams{Pair: pairParams.Pair, PricePrecisions: pairParams.OrderBookPricePrecisions, QtyPrecision: pairParams.QtyPrecision, CandleStickTimeframes: pairParams.CandleStickTimeframes})
	if err != nil {
		c.logger.Error("failed calling quote engine method", "error", err)
		return err
	}
	return nil
}

func (c *Client) DeleteOrderBook(pair string) error {
	_, err := c.client.DeleteOrderBook(context.Background(), &pbQ.Pair{Pair: pair})
	if err != nil {
		c.logger.Error("failed calling quote engine method", "error", err)
		return err
	}
	return nil
}
