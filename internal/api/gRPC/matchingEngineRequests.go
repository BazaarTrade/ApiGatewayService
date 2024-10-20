package clientGRPC

import (
	"context"

	"github.com/BazaarTrade/ApiGatewayService/internal/converter"
	"github.com/BazaarTrade/ApiGatewayService/internal/models"
	"github.com/BazaarTrade/MatchingEngineProtoGen/pbM"
)

func (c *GRPCClients) PlaceOrder(req *pbM.PlaceOrderReq) ([]models.Order, error) {
	pbUpdatedOrders, err := c.pbMClient.PlaceOrder(context.Background(), req)
	if err != nil {
		c.logger.Error("failed calling matching engine method", "error", err)
		return nil, err
	}

	var orders []models.Order
	for _, order := range pbUpdatedOrders.Orders {
		orders = append(orders, converter.ProtoOrderToOrder(order))
	}

	return orders, nil
}

func (c *GRPCClients) CancelOrder(req *pbM.OrderID) (models.Order, error) {
	pbOrder, err := c.pbMClient.CancelOrder(context.Background(), req)
	if err != nil {
		c.logger.Error("failed calling matching engine method", "error", err)
		return models.Order{}, err
	}

	return converter.ProtoOrderToOrder(pbOrder), nil
}

func (c *GRPCClients) GetCurrentOrders(req *pbM.UserID) ([]models.Order, error) {
	pbOrders, err := c.pbMClient.GetCurrentOrders(context.Background(), req)
	if err != nil {
		c.logger.Error("failed calling matching engine method", "error", err)
		return nil, err
	}

	var orders []models.Order
	for _, order := range pbOrders.Orders {
		orders = append(orders, converter.ProtoOrderToOrder(order))
	}

	return orders, nil
}

func (c *GRPCClients) GetOrders(req *pbM.UserID) ([]models.Order, error) {
	pbOrders, err := c.pbMClient.GetOrders(context.Background(), req)
	if err != nil {
		c.logger.Error("failed calling matching engine method", "error", err)
		return nil, err
	}

	var orders []models.Order
	for _, order := range pbOrders.Orders {
		orders = append(orders, converter.ProtoOrderToOrder(order))
	}

	return orders, nil
}

func (c *GRPCClients) CreateOrderBook(req *pbM.OrderBookSymbol) error {
	_, err := c.pbMClient.CreateOrderBook(context.Background(), req)
	if err != nil {
		c.logger.Error("failed calling matching engine method", "error", err)
		return err
	}

	return nil
}

func (c *GRPCClients) DeleteOrderBook(req *pbM.OrderBookSymbol) error {
	_, err := c.pbMClient.DeleteOrderBook(context.Background(), req)
	if err != nil {
		c.logger.Error("failed calling matching engine method", "error", err)
		return err
	}

	return nil
}
