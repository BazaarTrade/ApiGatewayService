package rest

import (
	"context"
	"net/http"
	"strconv"

	"github.com/BazaarTrade/ApiGatewayService/internal/converter"
	"github.com/BazaarTrade/ApiGatewayService/internal/models"
	"github.com/BazaarTrade/GeneratedProto/pb"
	"github.com/labstack/echo/v4"
)

func (h *Handler) placeOrder(c echo.Context) error {
	var req pb.PlaceOrderReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request format",
		})
	}

	if req.Price == "" || req.Qty == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "price and quantity must be specified",
		})
	}

	pbUpdatedOrders, err := h.pbClient.PlaceOrder(context.Background(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal error in matching engine",
		})
	}

	if len(pbUpdatedOrders.Orders) > 1 {
		for i := 1; i < len(pbUpdatedOrders.Orders); i++ {
			go h.hub.BroadcastUpdatedOrder(converter.PbOrderToOrder(pbUpdatedOrders.Orders[i]))
		}
	}

	return c.JSON(http.StatusOK, converter.PbOrderToOrder(pbUpdatedOrders.Orders[0]))
}

func (h *Handler) cancelOrder(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("orderID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid order_id",
		})
	}

	if id < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid order_id",
		})
	}

	order, err := h.pbClient.CancelOrder(context.Background(), &pb.OrderID{OrderID: int64(id)})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal error in matching engine service",
		})
	}

	return c.JSON(http.StatusOK, converter.PbOrderToOrder(order))
}

func (h *Handler) getCurrentOrders(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid user_id",
		})
	}

	if userID < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid user_id",
		})
	}

	pbOrders, err := h.pbClient.GetCurrentOrders(context.Background(), &pb.UserID{UserID: int64(userID)})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal error in matching engine service",
		})
	}

	var orders []models.Order
	for _, pbOrder := range pbOrders.Orders {
		orders = append(orders, converter.PbOrderToOrder(pbOrder))
	}

	return c.JSON(http.StatusOK, orders)
}

func (h *Handler) getOrders(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid user_id",
		})
	}

	if userID < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid user_id",
		})
	}

	pbOrders, err := h.pbClient.GetOrders(context.Background(), &pb.UserID{UserID: int64(userID)})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal error in matching engine service",
		})
	}

	var orders []models.Order
	for _, pbOrder := range pbOrders.Orders {
		orders = append(orders, converter.PbOrderToOrder(pbOrder))
	}

	return c.JSON(http.StatusOK, orders)
}

func (h *Handler) createOrderBook(c echo.Context) error {
	symbol := c.Param("symbol")
	if symbol == "" || len(symbol) < 4 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid symbol",
		})
	}

	_, err := h.pbClient.CreateOrderBook(context.Background(), &pb.OrderBookSymbol{Symbol: symbol})
	if err != nil {
		return c.JSON(http.StatusOK, map[string]string{
			"error": "internal error in matching engine service",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "orderbook created successfully",
	})
}

func (h *Handler) deleteOrderBook(c echo.Context) error {
	symbol := c.Param("symbol")
	if symbol == "" || len(symbol) < 4 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid symbol",
		})
	}

	_, err := h.pbClient.DeleteOrderBook(context.Background(), &pb.OrderBookSymbol{Symbol: symbol})
	if err != nil {
		return c.JSON(http.StatusOK, map[string]string{
			"error": "internal error in matching engine service",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "orderbook created successfully",
	})
}
