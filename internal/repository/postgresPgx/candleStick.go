package postgresPgx

import (
	"context"
	"time"

	"github.com/BazaarTrade/ApiGatewayService/internal/models"
	"github.com/jackc/pgx/v5"
)

func (p *Postgres) GetCandleStickHistory(req models.CandleStickHistoryRequest) ([]models.CandleStick, error) {
	var (
		rows pgx.Rows
		err  error
	)

	if req.CandleID == 0 {
		rows, err = p.db.Query(context.Background(), `
		SELECT id, pair, timeframe, openTime, closeTime, openPrice, closePrice, highPrice, lowPrice, volume, turnover 
		FROM quote.candlestick 
		WHERE pair = $1 AND timeframe = $2
		ORDER BY id DESC LIMIT $3
		`, req.Pair, req.Timeframe, req.Limit)
	} else {
		rows, err = p.db.Query(context.Background(), `
		SELECT id, pair, timeframe, openTime, closeTime, openPrice, closePrice, highPrice, lowPrice, volume, turnover 
		FROM quote.candlestick 
		WHERE pair = $1 AND timeframe = $2 AND id < $3
		ORDER BY id DESC LIMIT $4
		`, req.Pair, req.Timeframe, req.CandleID, req.Limit)
	}
	if err != nil {
		p.logger.Error("failed to get candleStick history", "error", err)
		return nil, err
	}

	var candleSticks []models.CandleStick

	defer rows.Close()

	for rows.Next() {
		var (
			openTime    time.Time
			closeTime   time.Time
			candleStick models.CandleStick
		)

		err := rows.Scan(&candleStick.ID, &candleStick.Pair, &candleStick.Timeframe, &openTime, &closeTime,
			&candleStick.OpenPrice, &candleStick.ClosePrice, &candleStick.HighPrice, &candleStick.LowPrice, &candleStick.Volume, &candleStick.Turnover)
		if err != nil {
			p.logger.Error("failed to scan candlestick", "error", err)
			return nil, err
		}

		candleStick.OpenTime = openTime.Format(time.RFC3339)
		candleStick.CloseTime = closeTime.Format(time.RFC3339)

		candleSticks = append(candleSticks, candleStick)
	}

	return candleSticks, nil
}
