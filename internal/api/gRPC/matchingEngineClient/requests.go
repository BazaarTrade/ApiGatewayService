package mClient

import (
	"context"

	"github.com/BazaarTrade/ApiGatewayService/internal/converter"
	"github.com/BazaarTrade/ApiGatewayService/internal/models"
	"github.com/BazaarTrade/MatchingEngineProtoGen/pbM"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Client) PlaceOrder(req *pbM.PlaceOrderReq) (models.Order, []models.Order, error) {
	placeOrderRes, err := c.client.PlaceOrder(context.Background(), req)
	if err != nil {
		c.logger.Error("failed calling matching engine method", "error", err)
		return models.Order{}, nil, err
	}

	var matchOrders = make([]models.Order, len(placeOrderRes.MatchOrders))
	for i, matchOrder := range placeOrderRes.MatchOrders {
		matchOrders[i] = converter.ProtoOrderToModelsOrder(matchOrder)
	}
	return converter.ProtoOrderToModelsOrder(placeOrderRes.Order), matchOrders, nil
}

func (c *Client) CancelOrder(req *pbM.OrderID) (models.Order, error) {
	pbOrder, err := c.client.CancelOrder(context.Background(), req)
	if err != nil {
		c.logger.Error("failed calling matching engine method", "error", err)
		return models.Order{}, err
	}
	return converter.ProtoOrderToModelsOrder(pbOrder), nil
}

func (c *Client) GetCurrentOrders(req *pbM.UserID) ([]models.Order, error) {
	pbOrders, err := c.client.GetCurrentOrders(context.Background(), req)
	if err != nil {
		c.logger.Error("failed calling matching engine method", "error", err)
		return nil, err
	}

	var orders []models.Order
	for _, order := range pbOrders.Orders {
		orders = append(orders, converter.ProtoOrderToModelsOrder(order))
	}

	return orders, nil
}

func (c *Client) GetOrders(req *pbM.UserID) ([]models.Order, error) {
	pbOrders, err := c.client.GetOrders(context.Background(), req)
	if err != nil {
		c.logger.Error("failed calling matching engine method", "error", err)
		return nil, err
	}

	var orders []models.Order
	for _, order := range pbOrders.Orders {
		orders = append(orders, converter.ProtoOrderToModelsOrder(order))
	}
	return orders, nil
}

func (c *Client) CreateOrderBook(req models.PairParams) error {
	_, err := c.client.CreateOrderBook(context.Background(), &pbM.PairParams{Pair: req.Pair, PricePrecisions: req.PricePrecisions, QtyPrecision: req.QtyPrecision})
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

func (c *Client) GetPairsParams() ([]models.PairParams, error) {
	pbPairsParams, err := c.client.GetPairsParams(context.Background(), &emptypb.Empty{})
	if err != nil {
		c.logger.Error("failed calling matching engine method", "error", err)
		return nil, err
	}

	var pairsParams []models.PairParams = make([]models.PairParams, len(pbPairsParams.PairParams))
	for i, pbPairParams := range pbPairsParams.PairParams {
		pairsParams[i] = converter.ProtoPairParamsToModelsPairParams(pbPairParams)
	}
	return pairsParams, nil
}

func (c *Client) GetPairPricePrecisions(pair string) ([]int32, error) {
	pbPricePreacisions, err := c.client.GetPairPricePrecisions(context.Background(), &pbM.Pair{Pair: pair})
	if err != nil {
		c.logger.Error("failed calling matching engine method", "error", err)
		return nil, err
	}
	return pbPricePreacisions.Precisions, nil
}
