package postgresPgx

import (
	"context"

	"github.com/BazaarTrade/ApiGatewayService/internal/models"
)

func (p *Postgres) GetPairsOrderBookPricePrecisions() (map[string][]int32, error) {
	rows, err := p.db.Query(context.Background(), `
	SELECT pair, orderBookPricePrecisions
	FROM quote.pairs
	`)
	if err != nil {
		p.logger.Error("failed to select pairs", "error", err)
		return nil, err
	}
	defer rows.Close()

	var orderBookPricePrecisions = make(map[string][]int32)
	for rows.Next() {
		var (
			pair            string
			pricePrecisions []int32
		)
		err := rows.Scan(&pair, &pricePrecisions)
		if err != nil {
			p.logger.Error("failed to scan pair", "error", err)
			return nil, err
		}
		orderBookPricePrecisions[pair] = pricePrecisions
	}
	return orderBookPricePrecisions, nil
}

func (p *Postgres) GetOrderBookPricePrecisions(pair string) ([]int32, error) {
	var orderBookpricePrecisions []int32
	err := p.db.QueryRow(context.Background(), `
	SELECT orderBookPricePrecisions
	FROM quote.pairs
	WHERE pair = $1
	`, pair).Scan(&orderBookpricePrecisions)
	if err != nil {
		p.logger.Error("failed to select pairs", "error", err)
		return nil, err
	}

	return orderBookpricePrecisions, nil
}

func (p *Postgres) CreatePair(pairParams models.PairParams) error {
	_, err := p.db.Exec(context.Background(), `
	INSERT INTO quote.pairs
	(pair, orderBookPricePrecisions, qtyPrecision, candleStickTimeframes)
	VALUES($1, $2, $3, $4)
	`, pairParams.Pair, pairParams.OrderBookPricePrecisions, pairParams.QtyPrecision, pairParams.CandleStickTimeframes)
	if err != nil {
		p.logger.Error("failed to insert pair", "error", err)
		return err
	}
	return nil
}

func (p *Postgres) GetCandleStickTimeframes(pair string) ([]string, error) {
	var candleStickTimeframes []string
	err := p.db.QueryRow(context.Background(), `
	SELECT candleStickTimeframes
	FROM quote.pairs
	WHERE pair = $1
	`, pair).Scan(&candleStickTimeframes)
	if err != nil {
		p.logger.Error("failed to select pairs", "error", err)
		return nil, err
	}

	return candleStickTimeframes, nil
}

func (p *Postgres) GetPairsCandleStickTimeframes() (map[string][]string, error) {
	rows, err := p.db.Query(context.Background(), `
	SELECT pair, candleStickTimeframes
	FROM quote.pairs
	`)
	if err != nil {
		p.logger.Error("failed to select pairs", "error", err)
		return nil, err
	}
	defer rows.Close()

	var candleStickTimeframes = make(map[string][]string)
	for rows.Next() {
		var (
			pair                 string
			candleStickTimeframe []string
		)
		err := rows.Scan(&pair, &candleStickTimeframe)
		if err != nil {
			p.logger.Error("failed to scan pair", "error", err)
			return nil, err
		}
		candleStickTimeframes[pair] = candleStickTimeframe
	}
	return candleStickTimeframes, nil
}

func (p *Postgres) GetPairsParams() ([]models.PairParams, error) {
	rows, err := p.db.Query(context.Background(), `
	SELECT pair, orderBookPricePrecisions, qtyPrecision, candleStickTimeframes
	FROM quote.pairs
	`)
	if err != nil {
		p.logger.Error("failed to select pairs", "error", err)
		return nil, err
	}
	defer rows.Close()

	var pairParams []models.PairParams
	for rows.Next() {
		var pairParam models.PairParams
		err := rows.Scan(&pairParam.Pair, &pairParam.OrderBookPricePrecisions, &pairParam.QtyPrecision, &pairParam.CandleStickTimeframes)
		if err != nil {
			p.logger.Error("failed to scan pair", "error", err)
			return nil, err
		}
		pairParams = append(pairParams, pairParam)
	}
	return pairParams, nil
}
