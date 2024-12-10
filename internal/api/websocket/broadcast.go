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
		Pair:      pOBS.Pair,
		Precision: precision,
	}]
	h.mu.RUnlock()

	if !exists || subscribers == nil {
		h.logger.Info("no subscribers for orderBook", "pair", pOBS.Pair, "precision", precision)
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
		Pair: precisedTrades.Pair,
	}]
	h.mu.RUnlock()

	if !exists || subscribers == nil {
		h.logger.Info("no subscribers for trades", "pair", precisedTrades.Pair)
		return
	}

	message := struct {
		Topic  string        `json:"topic"`
		Params models.Trades `json:"params"`
	}{
		Topic:  "trades",
		Params: precisedTrades,
	}

	tradesJSON, err := json.Marshal(message)
	if err != nil {
		h.logger.Error("failed to marshal precised trades message", "error", err)
		return
	}

	for client := range subscribers.Clients {
		go func(c *Client) {
			c.mu.Lock()
			defer c.mu.Unlock()

			err := c.Conn.Write(context.Background(), websocket.MessageText, tradesJSON)
			if err != nil {
				h.logger.Error("failed to write precised trades message to websocket", "error", err)
				return
			}
		}(client)
	}
}

func (h *Hub) BroadcastTicker(ticker models.Ticker, pair string) {
	h.mu.RLock()
	subscribers, exists := h.Subscribers["ticker"][TickerParams{
		Pair: pair,
	}]
	h.mu.RUnlock()

	if !exists || subscribers == nil {
		h.logger.Info("no subscribers for ticker", "pair", pair)
		return
	}

	message := struct {
		Topic  string `json:"topic"`
		Params struct {
			Pair   string        `json:"pair"`
			Ticker models.Ticker `json:"ticker"`
		} `json:"params"`
	}{
		Topic: "ticker",
		Params: struct {
			Pair   string        `json:"pair"`
			Ticker models.Ticker `json:"ticker"`
		}{Pair: pair, Ticker: ticker},
	}

	tickerJSON, err := json.Marshal(message)
	if err != nil {
		h.logger.Error("failed to marshal ticker message", "error", err)
		return
	}

	for client := range subscribers.Clients {
		go func(c *Client) {
			c.mu.Lock()
			defer c.mu.Unlock()

			err := c.Conn.Write(context.Background(), websocket.MessageText, tickerJSON)
			if err != nil {
				h.logger.Error("failed to write ticker message to websocket", "error", err)
				return
			}
		}(client)
	}
}
