package clientGRPC

import (
	"context"

	"github.com/BazaarTrade/ApiGatewayService/internal/models"
	"github.com/BazaarTrade/QuoteProtoGen/pbQ"
)

func (c *GRPCClients) ReadPrecisedOrderBookSnapshot() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := c.pbQClient.StreamPrecisedOrderBookSnapshot(ctx, &pbQ.Ping{})
	if err != nil {
		c.logger.Error("failed to connect to pbQ OBS stream", "error", err)
		return
	}

	c.logger.Info("successfully connected to pbQ POBSs stream")

	for {
		POBSs, err := stream.Recv()
		if err != nil {
			c.logger.Error("failed to receive POBSs", "error", err)
			return
		}

		for precision, pbPOBS := range POBSs.PrecisedOrderBookSnapshot {
			var POBS = models.OrderBookSnapshot{
				Symbol:  pbPOBS.Symbol,
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

func (c *GRPCClients) ReadPrecisedTrades() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := c.pbQClient.StreamPrecisedTrades(ctx, &pbQ.Ping{})
	if err != nil {
		c.logger.Error("failed to connect to pbQ precised trades stream", "error", err)
		return
	}

	c.logger.Info("successfully connected to pbQ precised trades stream")

	for {
		pbQPrecisedTrades, err := stream.Recv()
		if err != nil {
			c.logger.Error("failed to receive precised trades", "error", err)
			return
		}

		var precisedTrades = models.Trades{
			Symbol: pbQPrecisedTrades.Symbol,
			Trades: make([]models.Trade, len(pbQPrecisedTrades.Trade)),
		}

		for i, pbQTrade := range pbQPrecisedTrades.Trade {
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
