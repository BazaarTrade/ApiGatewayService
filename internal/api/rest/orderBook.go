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

	if pairParams.Pair == "" || len(pairParams.Pair) < 4 || len(pairParams.OrderBookPricePrecisions) < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}

	if err := s.db.CreatePair(pairParams); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "ivalid request",
		})
	}

	err := s.mClient.CreateOrderBook(pairParams.Pair)
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

	s.hub.AddOrderBookSnapshotTopic(pairParams.Pair, pairParams.OrderBookPricePrecisions)
	s.hub.AddTradesTopic(pairParams.Pair)
	s.hub.AddTickerTopic(pairParams.Pair)
	s.hub.AddCandleStickTopic(pairParams.Pair, pairParams.CandleStickTimeframes)
	s.qClient.StartStreamReaders(pairParams.Pair)

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

	orderBookPricePrecisions, err := s.db.GetOrderBookPricePrecisions(pair)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "pair not found",
		})
	}

	candleStickTimeframes, err := s.db.GetCandleStickTimeframes(pair)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "pair not found",
		})
	}

	s.qClient.StopStreamReadersByPair(pair)
	s.hub.RemoveOrderBookSnapshotTopic(pair, orderBookPricePrecisions)
	s.hub.RemoveTradesTopic(pair)
	s.hub.RemoveTickerTopic(pair)
	s.hub.RemoveCandleStickTopic(pair, candleStickTimeframes)

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

func (s *Server) getOrderBookPricePrecisions(c echo.Context) error {
	pair := c.Param("pair")
	if pair == "" || len(pair) < 4 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid pair",
		})
	}

	orderBookPricePrecisions, err := s.db.GetOrderBookPricePrecisions(pair)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "pair not found",
		})
	}

	return c.JSON(http.StatusOK, struct {
		OrderBookPricePrecisions []int32 `json:"orderBookPricePrecisions"`
	}{orderBookPricePrecisions})
}
