package qClient

import (
	"context"
	"errors"
	"io"

	"github.com/BazaarTrade/ApiGatewayService/internal/converter"
	"github.com/BazaarTrade/QuoteProtoGen/pbQ"
)

func (c *Client) StartStreamReaders(pair string) {
	go c.readPrecisedOrderBookSnapshot(pair)
	go c.readPrecisedTrade(pair)
	go c.readCandleStick(pair)
	go c.readTicker(pair)
}

func (c *Client) readPrecisedOrderBookSnapshot(pair string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := c.client.StreamPrecisedOrderBookSnapshot(ctx, &pbQ.Pair{Pair: pair})
	if err != nil {
		c.logger.Error("failed to connect to pbQ OBS stream", "error", err)
		return
	}

	c.logger.Debug("successfully connected to pOBSs stream", "pair", pair)

	for {
		pOBSs, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				c.logger.Debug("precised OBSs stream closed by server", "pair", pair)
				return
			}
			c.logger.Error("failed to receive pOBSs", "error", err)
			return
		}

		for orderBookprecision, pbPOBS := range pOBSs.PrecisedOrderBookSnapshot {
			go c.hub.BroadcastPrecisedOrderBookSnapshot(converter.PbQOBSToModelsOBS(pbPOBS), orderBookprecision)
		}
	}
}

func (c *Client) readPrecisedTrade(pair string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := c.client.StreamPrecisedTrades(ctx, &pbQ.Pair{Pair: pair})
	if err != nil {
		c.logger.Error("failed to connect to pbQ precised trade stream", "error", err)
		return
	}

	c.logger.Debug("successfully connected to pbQ precised trade stream", "pair", pair)

	for {
		pbQPrecisedTrades, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				c.logger.Debug("precised trade stream closed by server", "pair", pair)
				return
			}
			c.logger.Error("failed to receive precised trade", "error", err)
			return
		}

		go c.hub.BroadcastPrecisedTrades(converter.PbQTradeToModelsTrade(pbQPrecisedTrades))
	}
}

func (c *Client) readTicker(pair string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := c.client.StreamTicker(ctx, &pbQ.Pair{Pair: pair})
	if err != nil {
		c.logger.Error("failed to connect to pbQ ticker stream", "error", err)
		return
	}

	c.logger.Debug("successfully connected to ticker stream", "pair", pair)

	for {
		pbQTicker, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				c.logger.Debug("ticker stream closed by server", "pair", pair)
				return
			}
			c.logger.Error("failed to receive ticker", "error", err)
			return
		}

		go c.hub.BroadcastTicker(converter.PbQTickerToModelsTicker(pbQTicker))
	}
}

func (c *Client) readCandleStick(pair string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := c.client.StreamCandleStick(ctx, &pbQ.Pair{Pair: pair})
	if err != nil {
		c.logger.Error("failed to connect to pbQ candle stick stream", "error", err)
		return
	}

	c.logger.Debug("successfully connected to candleStick stream", "pair", pair)

	for {
		pbQCandleStick, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				c.logger.Debug("candle stick stream closed by server", "pair", pair)
				return
			}
			c.logger.Error("failed to receive candle stick", "error", err)
			return
		}

		go c.hub.BroadcastCandleStick(converter.PbQCandleStickToModelsCandleStick(pbQCandleStick))
	}
}
