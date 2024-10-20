package ws

import (
	"context"
	"encoding/json"

	"github.com/BazaarTrade/ApiGatewayService/internal/models"
	"github.com/coder/websocket"
)

func (h *Hub) BroadcastOrder(order models.Order) error {
	h.mu.RLock()
	user, exists := h.Users[order.UserID]
	h.mu.RUnlock()

	if !exists {
		return nil
	}

	message := struct {
		Topic string       `json:"topic"`
		Order models.Order `json:"order"`
	}{
		Topic: "orderUpdate",
		Order: order,
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		h.logger.Error("failed to marshal order message", "error", err)
		return err
	}

	for client := range user.Clients {
		go func(c *Client) {
			c.mu.Lock()
			defer c.mu.Unlock()

			err := c.Conn.Write(context.Background(), websocket.MessageText, messageJSON)
			if err != nil {
				h.logger.Error("failed to write order message to websocket", "error", err)
				return
			}
		}(client)
	}
	return nil
}

func (h *Hub) BroadcastPrecisedOrderBookSnapshot(pOBS models.OrderBookSnapshot, precision int32) {
	h.mu.RLock()
	subscribers, exists := h.Subscribers["orderBook"][OrderBookParams{
		Symbol:    pOBS.Symbol,
		Precision: precision,
	}]
	h.mu.RUnlock()

	if !exists || subscribers == nil {
		h.logger.Info("no subscribers for orderBook", "symbol", pOBS.Symbol, "precision", precision)
		return
	}

	message := struct {
		Topic  string                   `json:"topic"`
		Params models.OrderBookSnapshot `json:"params"`
	}{
		Topic:  "orderBook",
		Params: pOBS,
	}

	OBSJSON, err := json.Marshal(message)
	if err != nil {
		h.logger.Error("failed to marshal pOBS message", "error", err)
		return
	}

	for client := range subscribers.Clients {
		go func(c *Client) {
			c.mu.Lock()
			defer c.mu.Unlock()

			err := c.Conn.Write(context.Background(), websocket.MessageText, OBSJSON)
			if err != nil {
				h.logger.Error("failed to write pOBS message to websocket", "error", err)
				return
			}
		}(client)
	}
}

func (h *Hub) BroadcastPrecisedTrades(precisedTrades models.Trades) {
	h.mu.RLock()
	subscribers, exists := h.Subscribers["trades"][TradesParams{
		Symbol: precisedTrades.Symbol,
	}]
	h.mu.RUnlock()

	if !exists || subscribers == nil {
		h.logger.Info("no subscribers for trades", "symbol", precisedTrades.Symbol)
		return
	}

	message := struct {
		Topic  string        `json:"topic"`
		Params models.Trades `json:"params"`
	}{
		Topic:  "trades",
		Params: precisedTrades,
	}

	OBSJSON, err := json.Marshal(message)
	if err != nil {
		h.logger.Error("failed to marshal precised trades message", "error", err)
		return
	}

	for client := range subscribers.Clients {
		go func(c *Client) {
			c.mu.Lock()
			defer c.mu.Unlock()

			err := c.Conn.Write(context.Background(), websocket.MessageText, OBSJSON)
			if err != nil {
				h.logger.Error("failed to write precised trades message to websocket", "error", err)
				return
			}
		}(client)
	}
}
