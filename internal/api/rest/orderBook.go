package rest

import (
	"net/http"

	"github.com/BazaarTrade/ApiGatewayService/internal/models"
	"github.com/labstack/echo/v4"
)

func (s *Server) createOrderBook(c echo.Context) error {
	var pairParams models.PairParams

	if err := c.Bind(&pairParams); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request format",
		})
	}

	if pairParams.Pair == "" || len(pairParams.Pair) < 4 || len(pairParams.PricePrecisions) < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}

	err := s.mClient.CreateOrderBook(pairParams)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]string{
			"error": "internal error in matching engine service",
		})
	}

	err = s.qClient.CreateOrderBook(pairParams)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]string{
			"error": "internal error in quote service",
		})
	}

	s.qClient.StartStreamReaders(pairParams.Pair)
	s.hub.AddOrderBookUpdateTopic(pairParams.Pair, pairParams.PricePrecisions)
	s.hub.AddTradesTopic(pairParams.Pair)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "orderbook created successfully",
	})
}

func (s *Server) deleteOrderBook(c echo.Context) error {
	pair := c.Param("pair")
	if pair == "" || len(pair) < 4 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid pair",
		})
	}

	pricePrecisions, err := s.mClient.GetPairPricePrecisions(pair)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]string{
			"error": "internal error in matching engine service",
		})
	}

	s.hub.RemoveOrderBookUpdateTopic(pair, pricePrecisions)
	s.hub.RemoveTradesTopic(pair)
	s.hub.RemoveTickerTopic(pair)

	if err := s.qClient.DeleteOrderBook(pair); err != nil {
		return c.JSON(http.StatusOK, map[string]string{
			"error": "internal error in matching engine service",
		})
	}

	if err := s.mClient.DeleteOrderBook(pair); err != nil {
		return c.JSON(http.StatusOK, map[string]string{
			"error": "internal error in matching engine service",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "orderbook deleted successfully",
	})
}

func (s *Server) getPairPricePrecisions(c echo.Context) error {
	pair := c.Param("pair")
	if pair == "" || len(pair) < 4 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid pair",
		})
	}

	pricePrecisions, err := s.mClient.GetPairPricePrecisions(pair)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]string{
			"error": "internal error in matching engine service",
		})
	}

	return c.JSON(http.StatusOK, struct {
		PricePrecisions []int32 `json:"pricePrecisions"`
	}{pricePrecisions})
}
