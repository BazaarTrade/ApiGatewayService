package mClient

import (
	"context"

	"github.com/BazaarTrade/ApiGatewayService/internal/converter"
	"github.com/BazaarTrade/ApiGatewayService/internal/models"
	"github.com/BazaarTrade/MatchingEngineProtoGen/pbM"
)

func (c *Client) PlaceOrder(req *pbM.PlaceOrderReq) (models.Order, []models.Order, error) {
	placeOrderRes, err := c.client.PlaceOrder(context.Background(), req)
	if err != nil {
		c.logger.Error("failed calling matching engine method", "error", err)
		return models.Order{}, nil, err
	}

	var matchOrders = make([]models.Order, len(placeOrderRes.MatchOrders))
	for i, matchOrder := range placeOrderRes.MatchOrders {
		matchOrders[i] = converter.PbMOrderToModelsOrder(matchOrder)
	}
	return converter.PbMOrderToModelsOrder(placeOrderRes.Order), matchOrders, nil
}

func (c *Client) CancelOrder(req *pbM.OrderID) (models.Order, error) {
	pbOrder, err := c.client.CancelOrder(context.Background(), req)
	if err != nil {
		c.logger.Error("failed calling matching engine method", "error", err)
		return models.Order{}, err
	}
	return converter.PbMOrderToModelsOrder(pbOrder), nil
}

func (c *Client) CreateOrderBook(pair string) error {
	_, err := c.client.CreateOrderBook(context.Background(), &pbM.Pair{Pair: pair})
	if err != nil {
		c.logger.Error("failed calling matching engine method", "error", err)
		return err
	}
	return nil
}

func (c *Client) DeleteOrderBook(pair string) error {
	_, err := c.client.DeleteOrderBook(context.Background(), &pbM.Pair{Pair: pair})
	if err != nil {
		c.logger.Error("failed calling matching engine method", "error", err)
		return err
	}
	return nil
}
