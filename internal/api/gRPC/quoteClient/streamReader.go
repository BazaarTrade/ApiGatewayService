package qClient

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/BazaarTrade/ApiGatewayService/internal/converter"
	"github.com/BazaarTrade/QuoteProtoGen/pbQ"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var retry = time.Second * 5

func (c *Client) StartStreamReaders(pair string) {
	ctx, cancel := context.WithCancel(c.ctx)
	c.cancelByPair[pair] = cancel
	go c.readPrecisedOrderBookSnapshots(ctx, pair)
	go c.readPrecisedTrade(ctx, pair)
	go c.readCandleStick(ctx, pair)
	go c.readTicker(ctx, pair)
}

func (c *Client) readPrecisedOrderBookSnapshots(ctx context.Context, pair string) {
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("stopped pOBSs stream reader", "pair", pair)
			return
		default:
		}

		stream, err := c.client.StreamPrecisedOrderBookSnapshots(ctx, &pbQ.Pair{Pair: pair})
		if err != nil {
			c.logger.Warn("failed to connect to pbQ pOBSs stream", "pair", pair)
			time.Sleep(retry)
			continue
		}

		c.logger.Info("successfully connected to pOBSs stream", "pair", pair)

		for {
			select {
			case <-ctx.Done():
				c.logger.Info("stopped pOBSs stream reader", "pair", pair)
				return
			default:
			}

			pOBSs, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					c.logger.Info("pOBSs stream closed by server", "pair", pair)
					time.Sleep(retry)
					break
				}

				if status.Code(err) == codes.Canceled {
					c.logger.Info("stopped pOBSs stream reader", "pair", pair)
					return
				}

				c.logger.Error("failed to receive pOBSs", "pair", pair, "error", err, "status", status.Code(err))
				time.Sleep(retry)
				break
			}

			c.logger.Debug("recieved pOBSs", "pair", pair)

			for orderBookprecision, pbPOBS := range pOBSs.PrecisedOrderBookSnapshot {
				go c.hub.BroadcastPrecisedOrderBookSnapshot(converter.PbQOBSToModelsOBS(pbPOBS), orderBookprecision)
			}
		}
	}
}

func (c *Client) readPrecisedTrade(ctx context.Context, pair string) {
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("stopped precised trades stream reader", "pair", pair)
			return
		default:
		}

		stream, err := c.client.StreamPrecisedTrades(ctx, &pbQ.Pair{Pair: pair})
		if err != nil {
			c.logger.Warn("failed to connect to pbQ precised trade stream", "pair", pair)
			time.Sleep(retry)
			continue
		}

		c.logger.Info("successfully connected to pbQ precised trade stream", "pair", pair)

		for {
			select {
			case <-ctx.Done():
				c.logger.Info("stopped precised trades stream reader", "pair", pair)
				return
			default:
			}

			pbQPrecisedTrades, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					c.logger.Info("precised trade stream closed by server", "pair", pair)
					time.Sleep(retry)
					break
				}

				if status.Code(err) == codes.Canceled {
					c.logger.Info("stopped precised trades stream reader", "pair", pair)
					return
				}

				c.logger.Error("failed to receive precised trade", "pair", pair, "error", err, "status", status.Code(err))
				time.Sleep(retry)
				break
			}

			c.logger.Debug("recieved precised trade", "pair", pair)

			go c.hub.BroadcastPrecisedTrades(converter.PbQTradeToModelsTrade(pbQPrecisedTrades))
		}
	}

}

func (c *Client) readTicker(ctx context.Context, pair string) {
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("stopped ticker stream reader", "pair", pair)
			return
		default:
		}

		stream, err := c.client.StreamTicker(ctx, &pbQ.Pair{Pair: pair})
		if err != nil {
			c.logger.Warn("failed to connect to pbQ ticker stream", "pair", pair)
			time.Sleep(retry)
			continue
		}

		c.logger.Info("successfully connected to ticker stream", "pair", pair)

		for {
			select {
			case <-ctx.Done():
				c.logger.Info("stopped ticker stream reader", "pair", pair)
				return
			default:
			}

			pbQTicker, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					c.logger.Info("ticker stream closed by server", "pair", pair)
					time.Sleep(retry)
					break
				}

				if status.Code(err) == codes.Canceled {
					c.logger.Info("stopped ticker stream reader", "pair", pair)
					return
				}

				c.logger.Error("failed to receive ticker", "error", err)
				time.Sleep(retry)
				break
			}

			c.logger.Debug("recieved ticker", "pair", pair)

			go c.hub.BroadcastTicker(converter.PbQTickerToModelsTicker(pbQTicker))
		}
	}

}

func (c *Client) readCandleStick(ctx context.Context, pair string) {
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("stopped candleStick stream reader", "pair", pair)
			return
		default:
		}

		stream, err := c.client.StreamCandleStick(ctx, &pbQ.Pair{Pair: pair})
		if err != nil {
			c.logger.Warn("failed to connect to pbQ candle stick stream", "pair", pair)
			time.Sleep(retry)
			continue
		}

		c.logger.Info("successfully connected to candleStick stream", "pair", pair)

		for {
			select {
			case <-ctx.Done():
				c.logger.Info("stopped candleStick stream reader", "pair", pair)
				return
			default:
			}

			pbQCandleStick, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					c.logger.Info("candle stick stream closed by server", "pair", pair)
					time.Sleep(retry)
					break
				}

				if status.Code(err) == codes.Canceled {
					c.logger.Info("stopped candleStick stream reader", "pair", pair)
					return
				}

				c.logger.Error("failed to receive candle stick", "error", err)
				time.Sleep(retry)
				break
			}

			c.logger.Debug("recieved candleStick", "pair", pair)

			go c.hub.BroadcastCandleStick(converter.PbQCandleStickToModelsCandleStick(pbQCandleStick))
		}
	}
}
