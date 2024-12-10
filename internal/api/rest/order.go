package rest

import (
	"net/http"
	"strconv"

	"github.com/BazaarTrade/MatchingEngineProtoGen/pbM"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/status"
)

func (s *Server) placeOrder(c echo.Context) error {
	var req pbM.PlaceOrderReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request format",
		})
	}

	switch {
	case req.Type != "market" && req.Type != "limit":
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid order type",
		})

	case req.Type == "limit" && req.Price == "":
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid price",
		})

	case req.Qty == "":
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid qty",
		})
	}

	order, matchOrders, err := s.mClient.PlaceOrder(&req)
	if err != nil {
		st, _ := status.FromError(err)
		if st.Message() != "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": st.Message(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal error in matching engine",
		})
	}

	for _, matchOrder := range matchOrders {
		go s.hub.BroadcastOrder(matchOrder)
	}
	return c.JSON(http.StatusOK, order)
}

func (s *Server) cancelOrder(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("orderID"))
	if err != nil {
		s.logger.Error("failed to convert orderID to int", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid orderId",
		})
	}

	if id < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "orderId must be greater than 0",
		})
	}

	order, err := s.mClient.CancelOrder(&pbM.OrderID{OrderID: int64(id)})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal error in matching engine service",
		})
	}
	return c.JSON(http.StatusOK, order)
}

func (s *Server) getCurrentOrders(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		s.logger.Error("failed to convert userID to int", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid user_id",
		})
	}

	if userID < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid userId",
		})
	}

	orders, err := s.mClient.GetCurrentOrders(&pbM.UserID{UserID: int64(userID)})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal error in matching engine service",
		})
	}
	return c.JSON(http.StatusOK, orders)
}

func (s *Server) getOrders(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		s.logger.Error("failed to convert userID to int", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid user_id",
		})
	}

	if userID < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid userId",
		})
	}

	orders, err := s.mClient.GetOrders(&pbM.UserID{UserID: int64(userID)})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal error in matching engine service",
		})
	}
	return c.JSON(http.StatusOK, orders)
}
