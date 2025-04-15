package converter

import (
	"time"

	"github.com/BazaarTrade/ApiGatewayService/internal/models"
	"github.com/BazaarTrade/MatchingEngineProtoGen/pbM"
	"github.com/BazaarTrade/QuoteProtoGen/pbQ"
)

func PbMOrderToModelsOrder(order *pbM.Order) models.Order {
	return models.Order{
		ID:         int(order.ID),
		UserID:     int(order.UserID),
		IsBid:      order.IsBid,
		Pair:       order.Pair,
		Price:      order.Price,
		Qty:        order.Qty,
		SizeFilled: order.SizeFilled,
		Status:     order.Status,
		Type:       order.Type,
		CreatedAt:  order.CreatedAt.AsTime().Format(time.RFC3339),
		ClosedAt:   order.ClosedAt.AsTime().Format(time.RFC3339),
	}
}

func PbQTickerToModelsTicker(ticker *pbQ.Ticker) models.Ticker {
	return models.Ticker{
		Pair:      ticker.Pair,
		LastPrice: ticker.LastPrice,
		Change:    ticker.Change,
		HighPrice: ticker.HighPrice,
		LowPrice:  ticker.LowPrice,
		Volume:    ticker.Volume,
		Turnover:  ticker.Turnover,
	}
}

func PbQTradeToModelsTrade(pbQTrades *pbQ.Trades) []models.Trade {
	var trades = make([]models.Trade, 0, len(pbQTrades.Trades))
	for _, pbQTrade := range pbQTrades.Trades {
		trades = append(trades, models.Trade{
			Pair:  pbQTrade.Pair,
			IsBid: pbQTrade.IsBid,
			Price: pbQTrade.Price,
			Qty:   pbQTrade.Qty,
			Time:  pbQTrade.Time.AsTime().Format(time.RFC3339),
		})
	}
	return trades
}

func PbQCandleStickToModelsCandleStick(candleStick *pbQ.CandleStick) models.CandleStick {
	return models.CandleStick{
		ID:         int(candleStick.ID),
		Pair:       candleStick.Pair,
		Timeframe:  candleStick.Timeframe,
		OpenTime:   candleStick.OpenTime.AsTime().Format(time.RFC3339),
		CloseTime:  candleStick.CloseTime.AsTime().Format(time.RFC3339),
		OpenPrice:  candleStick.OpenPrice,
		ClosePrice: candleStick.ClosePrice,
		HighPrice:  candleStick.HighPrice,
		LowPrice:   candleStick.LowPrice,
		Volume:     candleStick.Volume,
		Turnover:   candleStick.Turnover,
		IsClosed:   candleStick.IsClosed,
	}
}

func PbQOBSToModelsOBS(OBS *pbQ.OrderBookSnapshot) models.OrderBookSnapshot {
	var bids = make([]models.Limit, 0, len(OBS.Bids))
	for _, bid := range OBS.Bids {
		bids = append(bids, models.Limit{
			Price: bid.Price,
			Qty:   bid.Qty,
		})
	}

	var asks = make([]models.Limit, 0, len(OBS.Asks))
	for _, ask := range OBS.Asks {
		asks = append(asks, models.Limit{
			Price: ask.Price,
			Qty:   ask.Qty,
		})
	}

	return models.OrderBookSnapshot{
		Pair:    OBS.Pair,
		Bids:    bids,
		Asks:    asks,
		BidsQty: OBS.BidsQty,
		AsksQty: OBS.AsksQty,
	}
}
