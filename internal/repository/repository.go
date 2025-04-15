package repository

import "github.com/BazaarTrade/ApiGatewayService/internal/models"

type Repository interface {
	GetPairsOrderBookPricePrecisions() (map[string][]int32, error)
	GetOrderBookPricePrecisions(string) ([]int32, error)
	GetCandleStickTimeframes(string) ([]string, error)
	GetNotFilledOrdersByUser(int) ([]models.Order, error)
	GetCandleStickHistory(models.CandleStickHistoryRequest) ([]models.CandleStick, error)
	GetClosedOrdersByUser(int) ([]models.Order, error)
	GetPairsParams() ([]models.PairParams, error)
	CreatePair(models.PairParams) error
}
