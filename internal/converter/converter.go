package converter

import (
	"time"

	"github.com/BazaarTrade/ApiGatewayService/internal/models"
	"github.com/BazaarTrade/MatchingEngineProtoGen/pbM"
	"github.com/BazaarTrade/QuoteProtoGen/pbQ"
)

func ProtoOrderToModelsOrder(order *pbM.Order) models.Order {
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

func ProtoPairParamsToModelsPairParams(pairParams *pbM.PairParams) models.PairParams {
	return models.PairParams{
		Pair:            pairParams.Pair,
		PricePrecisions: pairParams.PricePrecisions,
		QtyPrecision:    pairParams.QtyPrecision,
	}
}

func ProtoTickerToModelsTicker(ticker *pbQ.Ticker) models.Ticker {
	return models.Ticker{
		Price:     ticker.Price,
		Change:    ticker.Change,
		HighPrice: ticker.HighPrice,
		LowPrice:  ticker.LowPrice,
		Turnover:  ticker.Turnover,
	}
}
