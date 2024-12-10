package qClient

import (
	"context"
	"errors"
	"io"

	"github.com/BazaarTrade/ApiGatewayService/internal/converter"
	"github.com/BazaarTrade/ApiGatewayService/internal/models"
	"github.com/BazaarTrade/QuoteProtoGen/pbQ"
)

func (c *Client) StartStreamReaders(pair string) {
	go c.readPrecisedOrderBookSnapshot(pair)
	go c.readPrecisedTrades(pair)
	go c.readTicker(pair)
}

func (c *Client) readPrecisedOrderBookSnapshot(pair string) {
	stream, err := c.client.StreamPrecisedOrderBookSnapshot(context.Background(), &pbQ.Pair{Pair: pair})
	if err != nil {
		c.logger.Error("failed to connect to pbQ OBS stream", "error", err)
		return
	}

	defer func() {
		if err := stream.CloseSend(); err != nil {
			c.logger.Error("failed to close pOBSs stream", "pair", pair, "error", err)
		}
	}()

	c.logger.Info("successfully connected to pbQ POBSs stream")

	for {
		POBSs, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				c.logger.Info("precised OBSs stream closed by server", "pair", pair)
				return
			}
			c.logger.Error("failed to receive POBSs", "error", err)
			return
		}

		for precision, pbPOBS := range POBSs.PrecisedOrderBookSnapshot {
			var POBS = models.OrderBookSnapshot{
				Pair:    pbPOBS.Pair,
				BidsQty: pbPOBS.BidsQty,
				AsksQty: pbPOBS.AsksQty,
				Bids:    make([]models.Limit, len(pbPOBS.Bids)),
				Asks:    make([]models.Limit, len(pbPOBS.Asks)),
			}

			for i, pbBidLimit := range pbPOBS.Bids {
				POBS.Bids[i] = models.Limit{
					Price: pbBidLimit.Price,
					Qty:   pbBidLimit.Qty,
				}
			}

			for i, pbAskLimit := range pbPOBS.Asks {
				POBS.Asks[i] = models.Limit{
					Price: pbAskLimit.Price,
					Qty:   pbAskLimit.Qty,
				}
			}

			go func(POBS models.OrderBookSnapshot, precision int32) {
				c.hub.BroadcastPrecisedOrderBookSnapshot(POBS, precision)
			}(POBS, precision)
		}
	}
}

func (c *Client) readPrecisedTrades(pair string) {
	stream, err := c.client.StreamPrecisedTrades(context.Background(), &pbQ.Pair{Pair: pair})
	if err != nil {
		c.logger.Error("failed to connect to pbQ precised trades stream", "error", err)
		return
	}

	defer func() {
		if err := stream.CloseSend(); err != nil {
			c.logger.Error("failed to close precised trades stream", "pair", pair, "error", err)
		}
	}()

	c.logger.Info("successfully connected to pbQ precised trades stream")

	for {
		pbQPrecisedTrades, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				c.logger.Info("precised trades stream closed by server", "pair", pair)
				return
			}
			c.logger.Error("failed to receive precised trades", "error", err)
			return
		}

		var precisedTrades = models.Trades{
			Pair:   pbQPrecisedTrades.Pair,
			Trades: make([]models.Trade, len(pbQPrecisedTrades.Trades)),
		}

		for i, pbQTrade := range pbQPrecisedTrades.Trades {
			precisedTrades.Trades[i] = models.Trade{
				IsBid: pbQTrade.IsBid,
				Price: pbQTrade.Price,
				Qty:   pbQTrade.Qty,
				Time:  pbQTrade.Time.AsTime(),
			}
		}

		go func(trades models.Trades) {
			c.hub.BroadcastPrecisedTrades(trades)
		}(precisedTrades)
	}
}

func (c *Client) readTicker(pair string) {
	stream, err := c.client.StreamTicker(context.Background(), &pbQ.Pair{Pair: pair})
	if err != nil {
		c.logger.Error("failed to connect to pbQ ticker stream", "error", err)
		return
	}

	defer func() {
		if err := stream.CloseSend(); err != nil {
			c.logger.Error("failed to close precised trades stream", "pair", pair, "error", err)
		}
	}()

	for {
		pbQTicker, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				c.logger.Info("ticker stream closed by server", "pair", pair)
				return
			}
			c.logger.Error("failed to receive ticker", "error", err)
			return
		}

		go func(ticker models.Ticker) {
			c.hub.BroadcastTicker(ticker, pair)
		}(converter.ProtoTickerToModelsTicker(pbQTicker))
	}
}
