package api

import (
	"context"
	"net/http"
	"strconv"

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

	updatedOrders, err := h.pbClient.PlaceOrder(context.Background(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal error in matching engine",
		})
	}

	//TODO
	//create websocket and deliver updated orders to online users

	return c.JSON(http.StatusOK, updatedOrders.Orders[0])
}

func (h *Handler) cancelOrder(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("order_id"))
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

	return c.JSON(http.StatusOK, order)
}

func (h *Handler) getCurrentOrders(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
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

	orders, err := h.pbClient.GetCurrentOrders(context.Background(), &pb.UserID{UserID: int64(userID)})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal error in matching engine service",
		})
	}

	return c.JSON(http.StatusOK, orders)
}

func (h *Handler) getOrders(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
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

	orders, err := h.pbClient.GetOrders(context.Background(), &pb.UserID{UserID: int64(userID)})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal error in matching engine service",
		})
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
