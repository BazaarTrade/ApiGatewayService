package converter

import (
	"time"

	"github.com/BazaarTrade/ApiGatewayService/internal/models"
	"github.com/BazaarTrade/MatchingEngineProtoGen/pbM"
)

func ProtoOrderToOrder(order *pbM.Order) models.Order {
	var createdAt, closedAt string

	if order.CreatedAt != nil {
		createdAt = order.CreatedAt.AsTime().Format(time.RFC3339)
	}
	if order.ClosedAt != nil {
		closedAt = order.ClosedAt.AsTime().Format(time.RFC3339)
	}

	return models.Order{
		ID:         order.ID,
		UserID:     order.UserID,
		IsBid:      order.IsBid,
		Symbol:     order.Symbol,
		Price:      order.Price,
		Qty:        order.Qty,
		SizeFilled: order.SizeFilled,
		Status:     order.Status,
		Type:       order.Type,
		CreatedAt:  createdAt,
		ClosedAt:   closedAt,
	}
}
