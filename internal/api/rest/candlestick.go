package rest

import (
	"net/http"

	"github.com/BazaarTrade/ApiGatewayService/internal/models"
	"github.com/labstack/echo/v4"
)

func (s *Server) getCandleStickHistory(c echo.Context) error {
	var req models.CandleStickHistoryRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request format",
		})
	}

	candleSticks, err := s.db.GetCandleStickHistory(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request format",
		})
	}

	return c.JSON(http.StatusOK, candleSticks)
}
