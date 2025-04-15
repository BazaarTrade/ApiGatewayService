package postgresPgx

import (
	"context"

	"github.com/BazaarTrade/ApiGatewayService/internal/models"
)

func (p *Postgres) GetClosedOrdersByUser(userID int) ([]models.Order, error) {
	rows, err := p.db.Query(context.Background(), `
	SELECT id, userID, isBid, pair, price, qty, sizeFilled, status, type, createdAt, closedAt
	FROM matchingEngine.orders
	WHERE userID = $1
	`, userID)
	if err != nil {
		p.logger.Error("failed to select orders", "error", err)
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.IsBid,
			&order.Pair,
			&order.Price,
			&order.Qty,
			&order.SizeFilled,
			&order.Status,
			&order.Type,
			&order.CreatedAt,
			&order.ClosedAt,
		)
		if err != nil {
			p.logger.Error("failed to scan order", "error", err)
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (p *Postgres) GetNotFilledOrdersByUser(userID int) ([]models.Order, error) {
	rows, err := p.db.Query(context.Background(), `
	SELECT id, userID, isBid, pair, price, qty, sizeFilled, status, type, createdAt, closedAt
	FROM matchingEngine.orders
	WHERE userID = $1 AND status IN ('filling', 'filled', 'canceled')
	`, userID)
	if err != nil {
		p.logger.Error("failed to select orders", "error", err)
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.IsBid,
			&order.Pair,
			&order.Price,
			&order.Qty,
			&order.SizeFilled,
			&order.Status,
			&order.Type,
			&order.CreatedAt,
			&order.ClosedAt,
		)
		if err != nil {
			p.logger.Error("failed to scan order", "error", err)
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
