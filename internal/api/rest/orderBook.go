package rest

import (
	"net/http"

	"github.com/BazaarTrade/MatchingEngineProtoGen/pbM"
	"github.com/labstack/echo/v4"
)

func (s *Server) createOrderBook(c echo.Context) error {
	symbol := c.Param("symbol")
	if symbol == "" || len(symbol) < 4 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid symbol",
		})
	}

	err := s.clientGRPC.CreateOrderBook(&pbM.OrderBookSymbol{Symbol: symbol})
	if err != nil {
		return c.JSON(http.StatusOK, map[string]string{
			"error": "internal error in matching engine service",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "orderbook created successfully",
	})
}

func (s *Server) deleteOrderBook(c echo.Context) error {
	symbol := c.Param("symbol")
	if symbol == "" || len(symbol) < 4 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid symbol",
		})
	}

	err := s.clientGRPC.DeleteOrderBook(&pbM.OrderBookSymbol{Symbol: symbol})
	if err != nil {
		return c.JSON(http.StatusOK, map[string]string{
			"error": "internal error in matching engine service",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "orderbook created successfully",
	})
}
